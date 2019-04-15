package game

import (
	"github.com/DemoniacDeath/golangEbiten2dPlatformer/engine"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type World struct {
	engine.BaseGameObject
}

func (w *World) HandleKeyboardInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		w.Context.Quit = true
	}

	w.BaseGameObject.HandleKeyboardInput()
}
