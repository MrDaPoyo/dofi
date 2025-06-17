package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.Screen.Width, g.Screen.Height
}

func main() {
	var screen = ScreenSpecs{
		Width:           128,
		Height:          128,
		UpscalingFactor: 4,
	}
	var game = Game{
		Tabs:       []string{"cli", "code", "sprite", "map"},
		CurrentTab: 1,
		Screen:     screen,
	}
	ebiten.SetWindowSize(game.Screen.Width * game.Screen.UpscalingFactor, game.Screen.Height * game.Screen.UpscalingFactor)
	ebiten.SetWindowTitle("Dofi! :3")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
