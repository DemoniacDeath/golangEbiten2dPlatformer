package game

import "github.com/DemoniacDeath/golangEbiten2dPlatformer/engine"

type Consumable struct {
	engine.BaseGameObject
}

func NewConsumable(baseGameObject *engine.BaseGameObject) *Consumable {
	consumable := &Consumable{BaseGameObject: *baseGameObject}
	consumable.Physics = engine.NewPhysicsState(consumable)
	return consumable
}
