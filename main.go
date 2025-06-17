package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

const SCREEN_WIDTH = 128;
const SCREEN_HEIGHT = 128;
const UPSCALING_FACTOR = 4;

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func main() {
	ebiten.SetWindowSize(SCREEN_WIDTH * UPSCALING_FACTOR, SCREEN_HEIGHT * UPSCALING_FACTOR)
	ebiten.SetWindowTitle("Dofi! :3")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}