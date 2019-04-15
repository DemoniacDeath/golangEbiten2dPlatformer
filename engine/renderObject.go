package engine

import (
	"github.com/DemoniacDeath/golangEbiten2dPlatformer/core"
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"log"
	"math"
)

type RenderObject struct {
	Texture *ebiten.Image
	flipped bool
}

func (ro *RenderObject) Render(
	screen *ebiten.Image,
	windowSize core.Size,
	position core.Vector,
	size core.Size,
	cameraPosition core.Vector,
	cameraSize core.Size,
) {
	renderPosition := position.
		Plus(core.Vector{X: -size.Width / 2, Y: -size.Height / 2}).
		Minus(cameraPosition).
		Minus(core.Vector{X: -cameraSize.Width / 2, Y: -cameraSize.Height / 2})
	textureWidth, textureHeight := ro.Texture.Size()

	geom := ebiten.GeoM{}
	geom.Scale(
		(windowSize.Width*(size.Width/cameraSize.Width))/float64(textureWidth),
		(windowSize.Height*(size.Height/cameraSize.Height))/float64(textureHeight),
	)
	if ro.flipped {
		geom.Scale(-1, 1)
	}
	var flipShift = 0.0
	if ro.flipped {
		flipShift = size.Width
	}
	geom.Translate(
		math.Round(windowSize.Width*((renderPosition.X+flipShift)/cameraSize.Width)),
		math.Round(windowSize.Height*(renderPosition.Y/cameraSize.Height)),
	)

	_ = screen.DrawImage(ro.Texture, &ebiten.DrawImageOptions{GeoM: geom})
}

func NewRenderObjectFromColor(color color.Color, size int) *RenderObject {
	texture, err := ebiten.NewImage(size, size, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	_ = texture.Fill(color)
	return &RenderObject{Texture: texture}
}
