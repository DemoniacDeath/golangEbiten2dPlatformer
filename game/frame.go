package game

import (
	"github.com/DemoniacDeath/golangEbiten2dPlatformer/core"
	"github.com/DemoniacDeath/golangEbiten2dPlatformer/engine"
)

type Frame struct {
	engine.BaseGameObject
	Floor     *Solid
	LeftWall  *Solid
	RightWall *Solid
	Ceiling   *Solid
}

func NewFrame(b *engine.BaseGameObject, width float64) *Frame {
	frame := &Frame{
		BaseGameObject: *b,
		Floor: NewSolid(
			engine.NewBaseGameObject(
				b.Context,
				core.NewRect(
					0,
					b.Frame.Size.Height/2-width/2,
					b.Frame.Size.Width,
					width,
				),
			),
		),
		LeftWall: NewSolid(
			engine.NewBaseGameObject(
				b.Context,
				core.NewRect(
					-b.Frame.Size.Width/2+width/2,
					0,
					width,
					b.Frame.Size.Height-width*2,
				),
			),
		),
		RightWall: NewSolid(
			engine.NewBaseGameObject(
				b.Context,
				core.NewRect(
					b.Frame.Size.Width/2-width/2,
					0,
					width,
					b.Frame.Size.Height-width*2,
				),
			),
		),
		Ceiling: NewSolid(
			engine.NewBaseGameObject(
				b.Context,
				core.NewRect(
					0,
					-b.Frame.Size.Height/2+width/2,
					b.Frame.Size.Width,
					width,
				),
			),
		),
	}
	frame.AddChild(frame.Floor)
	frame.AddChild(frame.LeftWall)
	frame.AddChild(frame.RightWall)
	frame.AddChild(frame.Ceiling)
	return frame
}
