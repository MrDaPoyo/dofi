package main

import (
	"bytes"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/yuin/gopher-lua"
)

type Game struct {
	Screen ScreenSpecs
	LuaVM  *lua.LState
	Navbar Navbar
	Input  Input
}

type ScreenSpecs = struct {
	Width           int
	Height          int
	UpscalingFactor int
	Font            string
	FontSize        int
	LineBuffer      []*ebiten.Image
}

type Navbar = struct {
	Tabs         []string
	CurrentTab   int
	NavbarColor  color.RGBA
	NavbarHeight int
}

type LinearBuffer = struct {
	Id      int
	Content string
}

type Input = struct {
	Keys               []ebiten.Key
	Inputs             LinearBuffer
	CurrentInputString string
	CursorY            int
}

var (
	TextFaceSource *text.GoTextFaceSource
	TextFace       text.Face
)

func (g *Game) Update() error {
	log.Println(g.Input.CurrentInputString)
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.Input.Keys = []ebiten.Key{}
		g.Screen.LineBuffer = append(g.Screen.LineBuffer, ebiten.NewImage(g.Screen.Width, len(g.Screen.LineBuffer)*g.Screen.FontSize+2))
	} else {
		g.Input.Keys = inpututil.AppendJustPressedKeys(g.Input.Keys[:0])
		for _, k := range g.Input.Keys {
			if k == ebiten.KeySpace {
				g.Input.CurrentInputString += " "
			} else if k == ebiten.KeyBackspace {
				if len(g.Input.CurrentInputString) > 0 {
					g.Input.CurrentInputString = g.Input.CurrentInputString[:len(g.Input.CurrentInputString)-1]
				}
			} else if k != ebiten.KeyEscape || k != ebiten.KeyAlt {
				g.Input.CurrentInputString += k.String()
			}
		}
	}
	if len(g.Screen.LineBuffer) > 14 {
		g.Screen.LineBuffer = g.Screen.LineBuffer[0:14]
		g.Input.CursorY = g.Screen.Height - g.Screen.FontSize
	}
	// g.LuaVM.DoString(`print("hello")`)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// var NavbarBackground = ebiten.NewImage(g.Screen.Width, g.Navbar.NavbarHeight)
	// NavbarBackground.Fill(g.Navbar.NavbarColor)
	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(0, 0)
	// screen.DrawImage(NavbarBackground, op)
	//
	screen.Clear()
	for _, image := range g.Screen.LineBuffer {
		h := image.Bounds().Dy()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, float64(h))
		text.Draw(g.Screen.LineBuffer[len(g.Screen.LineBuffer)-1], "> "+g.Input.CurrentInputString, TextFace, &text.DrawOptions{})
		screen.DrawImage(image, op)
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
		Font:            "resources/font.ttf",
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
	screen.LineBuffer = lineBuffer

	var game = Game{
		Navbar: navbar,
		Screen: screen,
		LuaVM:  lua.NewState(),
		Input: Input{
			CursorY:     0,
			CurrentInputString: "",
		},
	}

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
