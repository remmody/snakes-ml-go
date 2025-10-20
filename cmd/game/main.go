package main

import (
	"log"
	"snakes-ml/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("AI Snake Game - Deep Q-Learning")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := game.NewGame(1280, 720)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
