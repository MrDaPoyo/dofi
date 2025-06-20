package main

import (
	"bytes"
	"image/color"
	"log"
	"math"
	"os"

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
	if err != nil {
		g.AppendLine(err.Error(), false)
	}
	g.AppendLine(g.Input.CurrentInputString, true)
}

func (g *Game) Update() error {
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
		g.Input.CurrentInputString = ""
	}

	// if input, change the contents
	if len(g.LinearBuffer) > 0 {
		g.ModifyLine(len(g.LinearBuffer)-1, g.Input.CurrentInputString)
	}

	return nil
}

func (g *Game) wrapText(value string) []string {
	maxChars := int(math.Round(float64(g.Screen.Width)/float64(g.Screen.FontWidth))) - g.Screen.FontWidth*2 - g.Screen.FontWidth/2
	var lines []string
	for len(value) > maxChars {
		lines = append(lines, value[:maxChars])
		value = value[maxChars:]
	}
	lines = append(lines, value)
	return lines
}

func (g *Game) AppendLine(value string, input bool) {
	wrapped := g.wrapText(value)
	lineHeight := g.Screen.FontSize + 1
	newLineHeight := lineHeight * len(wrapped)
	removedHeights := 0

	for {
		totalHeight := 0
		for _, line := range g.LinearBuffer {
			totalHeight += lineHeight * len(line.Content)
		}

		if totalHeight+newLineHeight <= g.Screen.Height {
			break
		}

		if len(g.LinearBuffer) == 0 {
			break
		}
		removedHeights += lineHeight * len(g.LinearBuffer[0].Content)
		g.LinearBuffer = g.LinearBuffer[1:]
	}

	g.LinearBuffer = append(g.LinearBuffer, LinearBuffer{
		Content: wrapped,
		IsInput: input,
	})
}

func (g *Game) ModifyLine(index int, value string) {
	wrapped := g.wrapText(value)
	lineHeight := g.Screen.FontSize + 1
	
	// remove lines from top if there isn't enough space
	for {
		totalHeight := 0
		for _, line := range g.LinearBuffer {
			totalHeight += lineHeight * len(line.Content)
		}

		if totalHeight <= g.Screen.Height {
			break
		}

		if len(g.LinearBuffer) == 0 {
			break
		}

		g.LinearBuffer = g.LinearBuffer[1:]
	}
	
	if len(wrapped) > 0 && index < len(g.LinearBuffer) {
		g.LinearBuffer[index].Content = wrapped
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// var NavbarBackground = ebiten.NewImage(g.Screen.Width, g.Navbar.NavbarHeight)
	// NavbarBackground.Fill(g.Navbar.NavbarColor)
	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(0, 0)
	// screen.DrawImage(NavbarBackground, op)

	screen.Clear()

	y := 1 // global vertical position
	for _, line := range g.LinearBuffer {
		prefix := "- "
		if line.IsInput {
			prefix = "> "
		}
		for _, wrappedLine := range line.Content {
			img := ebiten.NewImage(g.Screen.Width, g.Screen.FontSize+1)
			img.Fill(color.Black)

			text.Draw(img, prefix+wrappedLine, TextFace, &text.DrawOptions{})

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(0, float64(y))
			screen.DrawImage(img, op)

			y += g.Screen.FontSize + 1
			prefix = "  " // indent
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
