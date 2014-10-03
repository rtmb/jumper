package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

const SCALE = 3
const DIRECTION_RIGHT = 0
const DIRECTION_LEFT = 1
const STOP = -1
const RUNSPEED = SCALE * 5
const JUMPSPEED = -(SCALE * 8)
const GRAVITY = SCALE * 0.6
const TOPSPEED = SCALE * 8

type GameData struct {
	Spr      SpriteManager
	Lvl      Level
	Ply      Player
	renderer *sdl.Renderer
}

type GameObject interface {
	Update()
	Interp()
	Draw()
}

type Drawer interface {
	Draw()
}

type Updater interface {
	Update()
}

func main() {

	if 0 != sdl.Init(sdl.INIT_EVERYTHING) {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL: %s\n", sdl.GetError())
		os.Exit(2)
	}
	/*
		window := sdl.CreateWindow("goplot", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
			1600, 900, sdl.WINDOW_SHOWN|sdl.WINDOW_FULLSCREEN)
	*/
	window := sdl.CreateWindow("goplot", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if window == nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", sdl.GetError())
		os.Exit(2)
	}
	renderer := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if renderer == nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", sdl.GetError())
		os.Exit(2)
	}

	gd := Game_Init(renderer)

	gameOver := false

	dt := 0.03
	accumulator := 0.0

	currentTime := float64(sdl.GetTicks()) / 1000.0

	frameTime := 0.0

	for {
		/*
			begin mainloop
			implemented as described in
			http://gafferongames.com/game-physics/fix-your-timestep/
		*/
		newTime := float64(sdl.GetTicks()) / 1000.0
		frameTime = newTime - currentTime
		if frameTime > 0.25 {
			frameTime = 0.25
		}
		currentTime = newTime
		accumulator += frameTime
		for {
			if accumulator < dt {
				break
			}
			gd.Update()
			accumulator -= dt
		}

		alpha := accumulator / dt
		gd.Ply.Interpolate(alpha)

		gd.Draw()
		/*
			end mainloop
		*/

		if gameOver {
			break
		}

		var event sdl.Event
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				_ = t
				gameOver = true
				break
			}
		}
	}
	gd.Spr.TearDown()
	sdl.Quit()
}

func Game_Init(renderer *sdl.Renderer) GameData {
	spr := Init_from_json(GetDataPath()+"sprites.json", renderer)
	lvl := DummyLevel(spr, renderer)
	ply := Init_player(spr, renderer, lvl)
	return GameData{spr, lvl, ply, renderer}
}

func (gd *GameData) Draw() {
	gd.renderer.Clear()
	gd.Lvl.Draw()
	gd.Ply.Draw()
	gd.renderer.Present()
}

func (gd *GameData) Update() {
	gd.handleKeys()
	gd.Ply.Update()
}

func (gd *GameData) handleKeys() {
	keystate := sdl.GetKeyboardState()
	if keystate[sdl.GetScancodeFromName("LEFT")] == 1 {
		gd.Ply.SetDirection(DIRECTION_LEFT)
	}
	if keystate[sdl.GetScancodeFromName("RIGHT")] == 1 {
		gd.Ply.SetDirection(DIRECTION_RIGHT)
	}
	if keystate[sdl.GetScancodeFromName("RIGHT")]+keystate[sdl.GetScancodeFromName("LEFT")] == 0 {
		gd.Ply.SetDirection(STOP)
	}
	if keystate[sdl.GetScancodeFromName("SPACE")] == 1 {
		gd.Ply.Jump()
	}

}

func GetDataPath() string {
	return os.Getenv("GOPATH") + "/src/github.com/rtmb/jumper/data/"
}
