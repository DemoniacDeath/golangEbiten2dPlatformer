package main

import (
	"github.com/DemoniacDeath/golangEbiten2dPlatformer/core"
	"github.com/DemoniacDeath/golangEbiten2dPlatformer/engine"
	"github.com/DemoniacDeath/golangEbiten2dPlatformer/game"
	"log"
)

func main() {
	if err := game.NewGame(
		engine.NewContext(
			engine.NewSettings(
				"Golang ebiten platformer",
				core.Size{
					Width:  800,
					Height: 600,
				},
			),
		),
	).
		Run(); err != nil {
		log.Fatal(err)
	}
}
