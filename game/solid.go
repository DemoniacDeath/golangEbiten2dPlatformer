package game

import (
	"github.com/DemoniacDeath/golangEbiten2dPlatformer/engine"
	"math"
)

type Solid struct {
	engine.BaseGameObject
}

func NewSolid(b *engine.BaseGameObject) *Solid {
	solid := &Solid{BaseGameObject: *b}
	solid.Physics = engine.NewPhysicsState(solid)
	return solid
}

func (o *Solid) HandleEnterCollision(collision engine.Collision) {
	if collision.Collider.GetPhysics() != nil && collision.Collider.GetPhysics().Velocity.Y > 5 {
		collision.Collider.DealDamage(int(math.Round(collision.Collider.GetPhysics().Velocity.Y * 10)))
	}
}

func (o *Solid) HandleCollision(collision engine.Collision) {
	if math.Abs(collision.CollisionVector.X) < math.Abs(collision.CollisionVector.Y) {
		collision.Collider.GetBase().Frame.Center.X += collision.CollisionVector.X
		if collision.Collider.GetPhysics() != nil {
			collision.Collider.GetPhysics().Velocity.X = 0
		}
	} else {
		collision.Collider.GetBase().Frame.Center.Y += collision.CollisionVector.Y
		if collision.Collider.GetPhysics() != nil {
			collision.Collider.GetPhysics().Velocity.Y = 0
		}
	}
}
