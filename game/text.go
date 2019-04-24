package game

import (
	"../engine"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image/color"
)

type UIText struct {
	engine.BaseGameObject
	text  string
	face  font.Face
	color color.Color
}

func NewUIText(baseGameObject *engine.BaseGameObject) *UIText {
	return &UIText{BaseGameObject: *baseGameObject}
}

func NewUITextWith(baseGameObject *engine.BaseGameObject, text string, face font.Face, color color.Color) *UIText {
	uiText := &UIText{BaseGameObject: *baseGameObject, text: text, face: face, color: color}
	uiText.regenerate()
	return uiText
}

func (t *UIText) regenerate() {
	if t.face == nil || t.color == nil || len(t.text) == 0 {
		return
	}
	width := t.Frame.Size.Width * 4
	height := t.Frame.Size.Height * 4
	x := t.Frame.Center.X
	y := t.Frame.Center.Y + t.Frame.Size.Height*2

	image, _ := ebiten.NewImage(int(width), int(height), ebiten.FilterDefault)

	text.Draw(image, t.text, t.face, int(x), int(y), t.color)

	t.RenderObject = &engine.RenderObject{Texture: image}
}

func (t *UIText) SetText(text string) {
	t.text = text
	t.regenerate()
}

func (t *UIText) SetFontFace(face font.Face) {
	t.face = face
	t.regenerate()
}

func (t *UIText) SetColor(color color.Color) {
	t.color = color
	t.regenerate()
}
