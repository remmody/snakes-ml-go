package main

import (
	"log"
	"snakes-ml/config"
	"snakes-ml/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Window configuration from central config
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)
	ebiten.SetWindowTitle(config.WindowTitle)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Create and run game
	g := game.NewGame(config.WindowWidth, config.WindowHeight)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
