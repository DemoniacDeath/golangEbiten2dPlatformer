package game

import (
	"../core"
	"../engine"
	"github.com/hajimehoshi/ebiten"
)

type Camera struct {
	engine.BaseGameObject
	originalSize core.Size
}

func NewCamera(baseGameObject *engine.BaseGameObject) *Camera {
	return &Camera{BaseGameObject: *baseGameObject, originalSize: baseGameObject.Frame.Size.Div(2)}
}

func (w *Camera) HandleKeyboardInput() {
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		w.Frame.Size = w.originalSize.Times(2)
	} else {
		w.Frame.Size = w.originalSize
	}

	w.BaseGameObject.HandleKeyboardInput()
}
