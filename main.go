package main

import (
	"bytes"
	"image/color"
	"log"
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
}

type Navbar = struct {
	Tabs         []string
	CurrentTab   int
	NavbarColor  color.RGBA
	NavbarHeight int
}

type LinearBuffer = struct {
	Image   *ebiten.Image
	Content string
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
		g.Input.CurrentInputString = ""
		g.AppendLine(g.Input.CurrentInputString, true)
	}

	// if backspace'd, remove one character
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if len(g.Input.CurrentInputString) > 0 {
			g.Input.CurrentInputString = g.Input.CurrentInputString[:len(g.Input.CurrentInputString)-1]
		}
	}

	// if escaped, switch tabs (POYO (yes me) THIS IS A TODO)
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.Input.CurrentInputString = ""
	}

	// if input, change the contents
	if len(g.LinearBuffer) > 0 {
		g.ModifyLine(len(g.LinearBuffer)-1, g.Input.CurrentInputString)
	}

	return nil
}

func (g *Game) AppendLine(value string, input bool) {
	var line = LinearBuffer{
		Image:   ebiten.NewImage(g.Screen.Width, g.Screen.FontSize+1),
		Content: value,
		IsInput: input,
	}
	g.LinearBuffer = append(g.LinearBuffer, line)
	if len(g.LinearBuffer) > (g.Screen.Height/g.Screen.FontSize - 4) {
		g.LinearBuffer = g.LinearBuffer[1:] // slice and push older lines upwards
	}
}

func (g *Game) ModifyLine(index int, value string) {
	g.LinearBuffer[index].Content = value
}

func (g *Game) Draw(screen *ebiten.Image) {
	// var NavbarBackground = ebiten.NewImage(g.Screen.Width, g.Navbar.NavbarHeight)
	// NavbarBackground.Fill(g.Navbar.NavbarColor)
	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(0, 0)
	// screen.DrawImage(NavbarBackground, op)

	screen.Clear()
	for index := range g.LinearBuffer {
		var image = g.LinearBuffer[index]
		// h := image.Image.Bounds().Dy()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, float64(index*g.Screen.FontSize+index+1))
		image.Image.Fill(color.Black)
		text.Draw(image.Image, "> "+image.Content, TextFace, &text.DrawOptions{})
		screen.DrawImage(image.Image, op)
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

	game.AppendLine("", false)

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
