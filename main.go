package main

import (
	"log"

	"./core"
	"./engine"
	"./game"
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
