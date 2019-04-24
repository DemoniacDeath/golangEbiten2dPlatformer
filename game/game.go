package game

import (
	"github.com/golang/freetype/truetype"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"time"

	"../core"
	"../engine"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Game struct {
	context  *engine.Context
	camera   *Camera
	world    *World
	player   *Player
	ui       *engine.BaseGameObject
	firstRun bool
	textures map[TextureType]*ebiten.Image
}

func NewGame(context *engine.Context) *Game {
	return &Game{context: context, firstRun: true}
}

func (g *Game) Run() error {
	g.init()

	return ebiten.Run(g.update,
		int(g.context.Settings.WindowSize.Width),
		int(g.context.Settings.WindowSize.Height),
		1,
		g.context.Settings.Title)
}

type TextureType int

const (
	TextureBrick  TextureType = iota
	TextureIdle   TextureType = iota
	TextureMove   TextureType = iota
	TextureJump   TextureType = iota
	TextureCrouch TextureType = iota
)

func (g *Game) init() {

	textureSources := map[TextureType]string{
		TextureBrick:  "game/textures/brick.png",
		TextureIdle:   "game/textures/idle.png",
		TextureMove:   "game/textures/move.png",
		TextureJump:   "game/textures/jump.png",
		TextureCrouch: "game/textures/crouch.png",
	}
	g.textures = make(map[TextureType]*ebiten.Image, len(textureSources))

	for textureType, textureSource := range textureSources {
		texture, _, err := ebitenutil.NewImageFromFile(textureSource, ebiten.FilterDefault)
		if err != nil {
			panic(err)
		}
		g.textures[textureType] = texture
	}

	textureSize := 8
	gridSquareSize := float64(10)
	itemChance := 0.16

	g.world = &World{
		BaseGameObject: *engine.NewBaseGameObject(
			g.context,
			core.NewRect(
				0,
				0,
				float64(g.context.Settings.WindowSize.Width/2),
				float64(g.context.Settings.WindowSize.Height/2),
			),
		),
	}

	g.ui = engine.NewBaseGameObject(g.context, core.NewRect(
		0,
		0,
		float64(g.context.Settings.WindowSize.Width/4),
		float64(g.context.Settings.WindowSize.Height/4),
	))

	g.camera = NewCamera(
		engine.NewBaseGameObject(
			g.context,
			core.NewRect(
				0,
				0,
				float64(g.context.Settings.WindowSize.Width/2),
				float64(g.context.Settings.WindowSize.Height/2),
			),
		),
	)

	g.player = NewPlayer(
		engine.NewBaseGameObject(
			g.context,
			core.NewRect(
				0,
				0,
				gridSquareSize,
				gridSquareSize*2,
			),
		),
	)

	g.player.RenderObject = &engine.RenderObject{Texture: g.textures[TextureIdle]}
	g.player.speed = 1.3
	g.player.jumpSpeed = 2.5
	g.player.Physics.GravityForce = 0.1

	g.player.AddChild(g.camera)

	g.world.AddChild(g.player)

	frame := NewFrame(
		engine.NewBaseGameObject(
			g.context,
			core.NewRect(
				0,
				0,
				float64(g.context.Settings.WindowSize.Width/2),
				float64(g.context.Settings.WindowSize.Height/2),
			),
		),
		gridSquareSize,
	)

	black := engine.NewRenderObjectFromColor(color.Black, textureSize)

	frame.Floor.RenderObject = black
	frame.LeftWall.RenderObject = black
	frame.RightWall.RenderObject = black
	frame.Ceiling.RenderObject = black

	g.world.AddChild(frame)

	core.Srand()
	count := int(g.world.Frame.Size.Width * g.world.Frame.Size.Height * itemChance / (gridSquareSize * gridSquareSize))
	powerCount := count / 2
	g.player.maxPower = powerCount
	x := int(g.world.Frame.Size.Width/gridSquareSize - 2)
	y := int(g.world.Frame.Size.Height/gridSquareSize - 2)
	var rndX, rndY int
	takenX := make([]int, count)
	takenY := make([]int, count)
	for i := 0; i < count; i++ {
		for {
			taken := false
			rndX = core.Rand(0, x)
			rndY = core.Rand(0, y)
			for j := 0; j <= i; j++ {
				if rndX == takenX[j] && rndY == takenY[j] {
					taken = true
					break
				}
			}
			if !taken {
				break
			}
		}

		takenX[i] = rndX
		takenY[i] = rndY

		rect := core.NewRect(
			g.world.Frame.Size.Width/2-gridSquareSize*1.5-float64(rndX)*gridSquareSize,
			g.world.Frame.Size.Height/2-gridSquareSize*1.5-float64(rndY)*gridSquareSize,
			gridSquareSize,
			gridSquareSize)

		if powerCount > 0 {
			gameObject := NewConsumable(engine.NewBaseGameObject(g.context, rect))
			gameObject.RenderObject = engine.NewRenderObjectFromColor(color.RGBA{G: 0xff, A: 0xff}, 1)
			g.world.AddChild(gameObject)
			powerCount--
		} else {
			gameObject := NewSolid(engine.NewBaseGameObject(g.context, rect))
			gameObject.RenderObject = &engine.RenderObject{Texture: g.textures[TextureBrick]}
			g.world.AddChild(gameObject)
		}
	}

	var healthBarHolder = engine.NewBaseGameObject(g.context, core.NewRect(-g.world.Frame.Size.Width/4+16, -g.world.Frame.Size.Height/4+2.5, 30, 3))
	healthBarHolder.RenderObject = engine.NewRenderObjectFromColor(color.Black, 1)
	g.ui.AddChild(healthBarHolder)
	var powerBarHolder = engine.NewBaseGameObject(g.context, core.NewRect(g.world.Frame.Size.Width/4-16, -g.world.Frame.Size.Height/4+2.5, 30, 3))
	powerBarHolder.RenderObject = engine.NewRenderObjectFromColor(color.Black, 1)
	g.ui.AddChild(powerBarHolder)

	g.player.healthBar = NewUIBar(engine.NewBaseGameObject(g.context, core.Rect{core.Vector{}, core.Size{29, 2}}))
	g.player.healthBar.RenderObject = engine.NewRenderObjectFromColor(color.RGBA{R: 0xff, A: 0xff}, 1)
	healthBarHolder.AddChild(g.player.healthBar)
	g.player.healthBar.SetValue(100)
	g.player.powerBar = NewUIBar(engine.NewBaseGameObject(g.context, core.Rect{core.Vector{}, core.Size{29, 2}}))
	g.player.powerBar.RenderObject = engine.NewRenderObjectFromColor(color.RGBA{G: 0xff, A: 0xff}, 1)
	powerBarHolder.AddChild(g.player.powerBar)
	g.player.powerBar.SetValue(0)

	font, err := ioutil.ReadFile("game/fonts/Scratch_.ttf")
	if err != nil {
		log.Fatal(err)
	}
	tt, err := truetype.Parse(font)
	if err != nil {
		log.Fatal(err)
	}
	face := truetype.NewFace(tt, &truetype.Options{
		Size: 40,
		DPI:  72,
	})

	g.player.winText = NewUITextWith(
		engine.NewBaseGameObject(
			g.context,
			core.Rect{
				core.Vector{},
				core.Size{
					100,
					20,
				},
			},
		),
		"Congratulations! You won!",
		face,
		color.RGBA{G: 0xff, A: 0xff},
	)
	g.player.winText.SetVisible(false)
	g.ui.AddChild(g.player.winText)

	g.player.deathText = NewUITextWith(
		engine.NewBaseGameObject(
			g.context,
			core.Rect{
				core.Vector{},
				core.Size{
					100,
					20,
				},
			},
		),
		"You died! Game Over!",
		face,
		color.RGBA{R: 0xff, A: 0xff},
	)
	g.player.deathText.SetVisible(false)
	g.ui.AddChild(g.player.deathText)
}

func (g *Game) firstRunInit() {
	g.player.animations[PlayerAnimationIdle] = engine.NewAnimationWithSingleRenderObject(&engine.RenderObject{Texture: g.textures[TextureIdle]})
	g.player.animations[PlayerAnimationJump] = engine.NewAnimationWithSingleRenderObject(&engine.RenderObject{Texture: g.textures[TextureJump]})
	g.player.animations[PlayerAnimationCrouch] = engine.NewAnimationWithSingleRenderObject(&engine.RenderObject{Texture: g.textures[TextureCrouch]})
	g.player.animations[PlayerAnimationCrouchMove] = engine.NewAnimationWithSingleRenderObject(&engine.RenderObject{Texture: g.textures[TextureCrouch]})
	g.player.animations[PlayerAnimationMove] = engine.NewAnimationWithSpeedAndImage(
		1000,
		g.textures[TextureMove],
		40,
		80,
		6,
	)
}

func (g *Game) update(screen *ebiten.Image) error {

	if g.firstRun {
		g.firstRunInit()
	}

	if g.context.Quit {
		os.Exit(0)
	}

	g.world.HandleKeyboardInput()

	g.world.Clean()

	g.world.ProcessPhysics()

	g.world.DetectCollisions()

	g.world.Animate(time.Now().UnixNano())

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	_ = screen.Fill(color.White)

	g.world.Render(screen, g.world.Frame.Center, g.camera.GlobalPosition(), g.camera.Frame.Size)

	g.ui.Render(screen, g.ui.Frame.Center, core.Vector{}, g.camera.originalSize)

	if g.firstRun {
		g.firstRun = false
	}

	return nil
}
