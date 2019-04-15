package game

import (
	"awesomeProject/core"
	"awesomeProject/engine"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"math"
)

type PlayerAnimationType int

const (
	PlayerAnimationIdle PlayerAnimationType = iota
	PlayerAnimationJump PlayerAnimationType = iota
	PlayerAnimationCrouch PlayerAnimationType = iota
	PlayerAnimationCrouchMove PlayerAnimationType = iota
	PlayerAnimationMove PlayerAnimationType = iota
)

type Player struct {
	engine.BaseGameObject
	speed  float64
	jumpSpeed float64
	power int
	maxPower int
	jumped bool
	health int
	dead bool
	won bool
	crouched bool
	originalSize core.Size
	animations map[PlayerAnimationType]*engine.Animation
}

func NewPlayer(baseGameObject *engine.BaseGameObject) *Player {
	player := &Player{BaseGameObject: *baseGameObject, animations: map[PlayerAnimationType]*engine.Animation{} }
	player.Physics = engine.NewPhysicsState(player)
	player.GetPhysics().Gravity = true
	player.GetPhysics().Still = false
	player.originalSize = baseGameObject.Frame.Size
	player.health = 100
	return player
}

func (p *Player) HandleKeyboardInput() {
	if p.dead {
		p.BaseGameObject.HandleKeyboardInput()
		return
	}

	sitDown := false
	moveLeft := false
	moveRight := false
	moveVector := core.Vector{}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) ||
		ebiten.IsKeyPressed(ebiten.KeyA) {
		moveVector.X -= p.speed
		moveLeft = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) ||
		ebiten.IsKeyPressed(ebiten.KeyD) {
		moveVector.X += p.speed
		moveRight = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) ||
		ebiten.IsKeyPressed(ebiten.KeyW) ||
		ebiten.IsKeyPressed(ebiten.KeySpace){
		if p.Physics == nil || !p.Physics.Gravity {
			moveVector.Y -= p.speed
		} else if !p.jumped {
			p.Physics.Velocity.Y = p.Physics.Velocity.Y - p.jumpSpeed
			p.jumped = true
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) ||
		ebiten.IsKeyPressed(ebiten.KeyS) ||
		ebiten.IsKeyPressed(ebiten.KeyControl) {
		if p.Physics == nil || !p.Physics.Gravity {
			moveVector.Y += p.speed
		} else {
			sitDown = true
		}
	}
	if sitDown && !p.crouched {
		p.Frame.Center.Y += p.originalSize.Height / 4
		p.Frame.Size.Height = p.originalSize.Height / 2
	} else if !sitDown && p.crouched {
		p.Frame.Center.Y -= p.originalSize.Height / 4
		p.Frame.Size.Height = p.originalSize.Height
	}
	p.crouched = sitDown

	if moveLeft && !moveRight {
		p.animations[PlayerAnimationMove].SetTurnedLeft(true)
		p.animations[PlayerAnimationCrouch].SetTurnedLeft(true)
		p.animations[PlayerAnimationCrouchMove].SetTurnedLeft(true)
	}
	if moveRight && !moveLeft {
		p.animations[PlayerAnimationMove].SetTurnedLeft(false)
		p.animations[PlayerAnimationCrouch].SetTurnedLeft(false)
		p.animations[PlayerAnimationCrouchMove].SetTurnedLeft(false)
	}

	switch true {
	case !moveLeft && !moveRight && !p.jumped && !p.crouched:
		p.Animation = p.animations[PlayerAnimationIdle]
	case !moveLeft && !moveRight && !p.jumped && p.crouched:
		p.Animation = p.animations[PlayerAnimationCrouch]
	case (moveLeft || moveRight) && !p.jumped && !p.crouched:
		p.Animation = p.animations[PlayerAnimationMove]
	case (moveLeft || moveRight) && !p.jumped && p.crouched:
		p.Animation = p.animations[PlayerAnimationCrouchMove]
	case p.jumped && p.crouched:
		p.Animation = p.animations[PlayerAnimationCrouch]
	case p.jumped && !p.crouched:
		p.Animation = p.animations[PlayerAnimationJump]
	}




	p.Frame.Center = p.Frame.Center.Plus(moveVector)



	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		if p.Physics != nil {
			p.Physics.Gravity = !p.Physics.Gravity
			if !p.Physics.Gravity {
				p.jumped = true
				p.Physics.Velocity = core.Vector{}
			}
		}
	}

	p.BaseGameObject.HandleKeyboardInput()
}

func (p *Player) HandleEnterCollision(collision engine.Collision) {
	switch collision.Collider.(type) {
	case *Consumable:
		p.power += 1
		collision.Collider.Remove()
		p.speed += 0.01
		p.jumpSpeed += 0.01
		if p.power >= p.maxPower {
			p.Win()
		}
	}
}

func (p *Player) HandleExitCollision(collider engine.GameObject) {
	if p.Physics != nil && len(p.Physics.Colliders) == 0 {
		p.jumped = true
	}
}

func (p *Player) HandleCollision(collision engine.Collision) {
	if math.Abs(collision.CollisionVector.X) > math.Abs(collision.CollisionVector.Y) {
		if collision.CollisionVector.Y > 0 && p.jumped && p.Physics != nil && p.Physics.Gravity {
			p.jumped = false
		}
	}
}

func (p *Player) DealDamage(damage int) {
	if !p.won {
		p.health -= damage
		if p.health < 0 {
			p.Die()
		}
	}
}

func (p *Player) Die() {
	p.dead = true
}

func (p *Player) Win() {
	p.won = true
}