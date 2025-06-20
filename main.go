package main

import (
	"bytes"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	lua "github.com/yuin/gopher-lua"
)

type Game struct {
	Screen       ScreenSpecs
	LuaVM        *lua.LState
	Navbar       Navbar
	Input        Input
	LinearBuffer []LinearBuffer
}

type ScreenSpecs = struct {
	Width           int
	Height          int
	UpscalingFactor int
	Font            string
	FontSize        int
	FontWidth       int
	Buffer          [128][128]color.RGBA
}

type Navbar = struct {
	Tabs         []string
	CurrentTab   int
	CliEnabled   bool
	NavbarColor  color.RGBA
	NavbarHeight int
}

type LinearBuffer = struct {
	Content []string
	IsInput bool
}

type Input = struct {
	Keys               []ebiten.Key
	CurrentInputString string
}

var (
	TextFaceSource *text.GoTextFaceSource
	TextFace       text.Face
)

func (g *Game) HandleCommand(command string) {
	var err = g.LuaVM.DoString(command)
	if strings.Split(command, " ")[0] == "run" {
		err := g.LuaVM.DoFile(strings.TrimSpace(strings.Split(command, " ")[1]))
		if err != nil {
			log.Println("Error executing file:", err)
			g.AppendLine(err.Error(), false)
		}
		return
	}
	if err != nil {
		log.Println("Error executing command:", err)
		g.AppendLine(err.Error(), false)
	}
	g.AppendLine(g.Input.CurrentInputString, true)
}

func (g *Game) Update() error {
	if updateFn := g.LuaVM.GetGlobal("_update"); updateFn != lua.LNil {
		if err := g.LuaVM.CallByParam(lua.P{
			Fn:      updateFn,
			NRet:    0,
			Protect: true,
		}); err != nil {
			g.AppendLine("Lua error in _update: "+err.Error(), false)
		}
	}

	var inputChars []rune
	inputChars = ebiten.AppendInputChars(inputChars[:0])

	for _, r := range inputChars {
		switch r {
		case '\r', '\n':
			g.HandleCommand(g.Input.CurrentInputString)
			g.Input.CurrentInputString = ""
		default:
			g.Input.CurrentInputString += string(r)
		}
	}

	// on enter, just clear the "buffer" (if that's the word) and append the contents as a new line.
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.Input.Keys = []ebiten.Key{}
		g.HandleCommand(g.Input.CurrentInputString)
		g.Input.CurrentInputString = ""
	}

	// if backspace'd, remove one character
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if len(g.Input.CurrentInputString) > 0 {
			g.Input.CurrentInputString = g.Input.CurrentInputString[:len(g.Input.CurrentInputString)-1]
		}
	}

	// if escaped, switch tabs (POYO (yes me) THIS IS A TODO) (switching tabs here)
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.Navbar.CliEnabled = !g.Navbar.CliEnabled
	}

	// if input and CliEnabled, change the contents
	if len(g.LinearBuffer) > 0 && g.Navbar.CliEnabled {
		g.ModifyLine(len(g.LinearBuffer)-1, g.Input.CurrentInputString)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()

	if drawFn := g.LuaVM.GetGlobal("_draw"); drawFn != lua.LNil {
		if err := g.LuaVM.CallByParam(lua.P{
			Fn:      drawFn,
			NRet:    0,
			Protect: true,
		}); err != nil {
			g.AppendLine("Lua error in _draw: "+err.Error(), false)
		}
	}

	bufferImg := ebiten.NewImage(len(g.Screen.Buffer[0]), len(g.Screen.Buffer))

	pixels := make([]byte, len(g.Screen.Buffer)*len(g.Screen.Buffer[0])*4)
	for y := 0; y < len(g.Screen.Buffer); y++ {
		for x := 0; x < len(g.Screen.Buffer[y]); x++ {
			pixel := g.Screen.Buffer[y][x]
			idx := (y*len(g.Screen.Buffer[0]) + x) * 4
			pixels[idx] = pixel.R
			pixels[idx+1] = pixel.G
			pixels[idx+2] = pixel.B
			pixels[idx+3] = pixel.A
		}
	}

	bufferImg.WritePixels(pixels)
	screen.DrawImage(bufferImg, &ebiten.DrawImageOptions{})

	if g.Navbar.CliEnabled {
		lineHeight := g.Screen.FontSize + 1
		var totalLines int
		for _, line := range g.LinearBuffer {
			totalLines += len(line.Content)
		}

		y := g.Screen.Height - totalLines*lineHeight
		for _, line := range g.LinearBuffer {
			prefix := "- "
			if line.IsInput {
				prefix = "> "
			}
			for _, wrappedLine := range line.Content {
				img := ebiten.NewImage(g.Screen.Width, lineHeight)
				img.Fill(color.Black)
				text.Draw(img, prefix+wrappedLine, TextFace, &text.DrawOptions{})

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(0, float64(y))
				screen.DrawImage(img, op)

				y += lineHeight
				prefix = "  "
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.Screen.Width, g.Screen.Height
}

func main() {
	var screen = ScreenSpecs{
		Width:           128,
		Height:          128,
		UpscalingFactor: 4,
		Font:            "resources/cg-pixel-4x5-mono.otf",
		FontSize:        5,
		FontWidth:       4,
	}

	var navbar = Navbar{
		Tabs:         []string{"cli", "code", "sprite", "map"},
		CurrentTab:   1,
		NavbarColor:  color.RGBA{0x7F, 0x11, 0xE0, 0xff}, // purple-ish
		NavbarHeight: 10,
		CliEnabled:   true,
	}

	var lineBuffer = make([]*ebiten.Image, 1)
	lineBuffer[0] = ebiten.NewImage(screen.Width, 128)

	var game = Game{
		Navbar: navbar,
		Screen: screen,
		LuaVM:  lua.NewState(),
		Input: Input{
			CurrentInputString: "",
		},
	}

	game.setupLuaAPI()
	defer game.LuaVM.Close()
	game.AppendLine("", true)

	var font, err = os.ReadFile(screen.Font)
	if err != nil {
		log.Fatal(err)
	}
	TextFaceSource, err = text.NewGoTextFaceSource(bytes.NewReader(font))
	if err != nil {
		log.Fatal(err)
	}
	TextFace = &text.GoTextFace{
		Source: TextFaceSource,
		Size:   float64(game.Screen.FontSize),
	}

	ebiten.SetWindowSize(game.Screen.Width*game.Screen.UpscalingFactor, game.Screen.Height*game.Screen.UpscalingFactor)
	ebiten.SetWindowTitle("Dofi! :3")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
