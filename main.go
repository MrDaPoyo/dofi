package main

import (
	"bytes"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Game struct {
	Tabs       []string
	CurrentTab int
	Screen     ScreenSpecs
}

type ScreenSpecs = struct {
	Width           int
	Height          int
	UpscalingFactor int
	Font            string
	FontSize        int
}

var (
	TextFaceSource *text.GoTextFaceSource
	TextFace       text.Face
)

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	text.Draw(screen, "Ellos World :333", TextFace, &text.DrawOptions{})
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
		FontSize:        8,
	}

	var game = Game{
		Tabs:       []string{"cli", "code", "sprite", "map"},
		CurrentTab: 1,
		Screen:     screen,
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
