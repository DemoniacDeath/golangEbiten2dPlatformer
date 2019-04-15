package engine

import (
	"github.com/hajimehoshi/ebiten"
	"image"
)

const speedScale = 100000

type Animation struct {
	frames []*RenderObject
	speed int
	startTick int64
	turnedLeft bool

}

func (a *Animation) Animate(ticks int64) *RenderObject {
	if a.startTick == 0 || ticks - a.startTick >= int64(len(a.frames) * a.speed * speedScale) {
		a.startTick = ticks
	}
	return a.frames[int(ticks - a.startTick) / (a.speed * speedScale)]
}

func (a *Animation) SetTurnedLeft(value bool) {
	if a.turnedLeft != value {
		for _, frame := range a.frames {
			frame.flipped = value
		}
	}
	a.turnedLeft = value
}

func NewAnimationWithSingleRenderObject(renderObject *RenderObject) *Animation {
	return &Animation{
		frames: []*RenderObject{renderObject},
		speed: 1,
	}
}

func NewAnimationWithSpeedAndImage(speed int, sourceImage *ebiten.Image, width int, height int, numberOfFrames int) *Animation {
	animation := &Animation{
		frames: make([]*RenderObject, numberOfFrames),
		speed: speed,
	}
	for i := 0; i < numberOfFrames; i++ {
		frameImage, err := ebiten.NewImageFromImage(
			sourceImage.SubImage(
				image.Rect(
					0, i * height,
					width, (i+1) * height - 1,
				),
			),
			ebiten.FilterDefault,
		)
		if err != nil {
			panic(err)
		}
		animation.frames[i] = &RenderObject{Texture: frameImage}
	}
	return animation
}