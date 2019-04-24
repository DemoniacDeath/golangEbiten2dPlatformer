package engine

import (
	"../core"
	"github.com/hajimehoshi/ebiten"
)

type GameObject interface {
	HandleKeyboardInput()
	HandleEnterCollision(Collision)
	HandleCollision(Collision)
	HandleExitCollision(GameObject)
	Animate(int64)
	Render(*ebiten.Image, core.Vector, core.Vector, core.Size)
	ProcessPhysics()
	DetectCollisions()
	AddChild(GameObject)
	Clean()
	GlobalPosition() core.Vector
	collectColliders(*[]GameObject)
	setParent(GameObject)
	isRemoved() bool
	GetPhysics() *PhysicsState
	GetBase() *BaseGameObject
	GetFrame() core.Rect
	Remove()
	DealDamage(damage int)
	SetVisible(bool)
}

type BaseGameObject struct {
	children     map[GameObject]bool
	parent       GameObject
	Context      *Context
	Frame        core.Rect
	RenderObject *RenderObject
	Animation    *Animation
	Physics      *PhysicsState
	visible      bool
	removed      bool
}

func NewBaseGameObject(context *Context, frame core.Rect) *BaseGameObject {
	return &BaseGameObject{
		Context:  context,
		Frame:    frame,
		visible:  true,
		removed:  false,
		children: map[GameObject]bool{},
	}
}

func (o *BaseGameObject) HandleKeyboardInput() {
	for child := range o.children {
		child.HandleKeyboardInput()
	}
}

func (o *BaseGameObject) HandleEnterCollision(collision Collision) {}
func (o *BaseGameObject) HandleCollision(collision Collision)      {}
func (o *BaseGameObject) HandleExitCollision(collider GameObject)  {}

func (o *BaseGameObject) Animate(now int64) {
	if o.Animation != nil {
		o.RenderObject = o.Animation.Animate(now)
	}
	for child := range o.children {
		child.Animate(now)
	}
}

func (o *BaseGameObject) ProcessPhysics() {
	if o.Physics != nil {
		o.Physics.Change()
	}
	for child := range o.children {
		child.ProcessPhysics()
	}
}

func (o *BaseGameObject) AddChild(child GameObject) {
	o.children[child] = true
	child.setParent(o)
}

func (o *BaseGameObject) setParent(parent GameObject) {
	o.parent = parent
}

func (o *BaseGameObject) Clean() {
	if o.Physics != nil {
		o.Physics.Clean()
	}
	for child := range o.children {
		if child.isRemoved() {
			delete(o.children, child)
		}
	}
	for child := range o.children {
		child.Clean()
	}
}

func (o *BaseGameObject) isRemoved() bool {
	return o.removed
}

func (o *BaseGameObject) Render(screen *ebiten.Image, localBasis core.Vector, cameraPosition core.Vector, cameraSize core.Size) {
	if o.visible && o.RenderObject != nil {
		o.RenderObject.Render(
			screen,
			o.Context.Settings.WindowSize,
			o.Frame.Center.Plus(localBasis),
			o.Frame.Size,
			cameraPosition,
			cameraSize,
		)
	}

	for child := range o.children {
		child.Render(screen, o.Frame.Center.Plus(localBasis), cameraPosition, cameraSize)
	}
}

func (o *BaseGameObject) GlobalPosition() core.Vector {
	if o.parent != nil {
		return o.Frame.Center.Plus(o.parent.GlobalPosition())
	} else {
		return o.Frame.Center
	}
}

func (o *BaseGameObject) DetectCollisions() {
	var allColliders []GameObject
	o.collectColliders(&allColliders)

	for i := range allColliders {
		for j := range allColliders {
			physics := allColliders[i].GetPhysics()
			if physics != nil {
				physics.DetectCollisions(allColliders[j].GetPhysics())
			}
		}
	}
}

func (o *BaseGameObject) collectColliders(allColliders *[]GameObject) {
	if o.Physics != nil {
		*allColliders = append(*allColliders, o)
	}

	for child := range o.children {
		child.collectColliders(allColliders)
	}
}

func (o *BaseGameObject) GetPhysics() *PhysicsState {
	return o.Physics
}

func (o *BaseGameObject) GetBase() *BaseGameObject {
	return o
}

func (o *BaseGameObject) GetFrame() core.Rect {
	return o.Frame
}

func (o *BaseGameObject) Remove() {
	o.removed = true
}

func (o *BaseGameObject) DealDamage(damage int) {}

func (o *BaseGameObject) SetVisible(v bool) {
	o.visible = v
}
