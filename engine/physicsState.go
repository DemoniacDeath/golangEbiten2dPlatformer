package engine

import (
	"github.com/DemoniacDeath/golangEbiten2dPlatformer/core"
	"math"
)

type PhysicsState struct {
	gameObject   GameObject
	Velocity     core.Vector
	Gravity      bool
	Still        bool
	GravityForce float64
	Colliders    map[GameObject]bool
}

func NewPhysicsState(gameObject GameObject) *PhysicsState {
	return &PhysicsState{gameObject: gameObject, Still: true, Colliders: map[GameObject]bool{}}
}

func (p *PhysicsState) Clean() {
	for collider := range p.Colliders {
		if collider.isRemoved() {
			delete(p.Colliders, collider)
		}
	}
}

func (p *PhysicsState) Change() {
	if p.Gravity {
		p.Velocity = p.Velocity.Plus(core.Vector{Y: p.GravityForce})
	}
	p.gameObject.GetBase().Frame.Center = p.gameObject.GetBase().Frame.Center.Plus(p.Velocity)
}

func (p *PhysicsState) DetectCollisions(c *PhysicsState) {
	if c == nil {
		return
	}
	if p.Still && c.Still {
		return
	}

	x1 := p.gameObject.GlobalPosition().X - p.gameObject.GetFrame().Size.Width/2
	x2 := c.gameObject.GlobalPosition().X - c.gameObject.GetFrame().Size.Width/2
	X1 := x1 + p.gameObject.GetFrame().Size.Width
	X2 := x2 + c.gameObject.GetFrame().Size.Width
	y1 := p.gameObject.GlobalPosition().Y - p.gameObject.GetFrame().Size.Height/2
	y2 := c.gameObject.GlobalPosition().Y - c.gameObject.GetFrame().Size.Height/2
	Y1 := y1 + p.gameObject.GetFrame().Size.Height
	Y2 := y2 + c.gameObject.GetFrame().Size.Height

	dX1 := X1 - x2
	dX2 := x1 - X2
	dY1 := Y1 - y2
	dY2 := y1 - Y2

	alreadyCollided := p.Colliders[c.gameObject] || c.Colliders[p.gameObject]

	if dX1 > 0 &&
		dX2 < 0 &&
		dY1 > 0 &&
		dY2 < 0 {
		overlapX := dX2
		if math.Abs(dX1) < math.Abs(dX2) {
			overlapX = dX1
		}
		overlapY := dY2
		if math.Abs(dY1) < math.Abs(dY2) {
			overlapY = dY1
		}
		overlapArea := core.Vector{X: overlapX, Y: overlapY}
		if !alreadyCollided {
			p.Colliders[c.gameObject] = true
			c.Colliders[p.gameObject] = true

			p.gameObject.HandleEnterCollision(Collision{c.gameObject, overlapArea})
			c.gameObject.HandleEnterCollision(Collision{p.gameObject, overlapArea.Times(-1)})
		}

		p.gameObject.HandleCollision(Collision{c.gameObject, overlapArea})
		c.gameObject.HandleCollision(Collision{p.gameObject, overlapArea.Times(-1)})
	} else if alreadyCollided {
		delete(p.Colliders, c.gameObject)
		delete(c.Colliders, p.gameObject)
		p.gameObject.HandleExitCollision(c.gameObject)
		c.gameObject.HandleExitCollision(p.gameObject)
	}
}
