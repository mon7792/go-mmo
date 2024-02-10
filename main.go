package main

//go:generate packer --input images --stats

import (
	"fmt"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/mon7792/go-mmo/engine/asset"
)

func run() {
	// all the game code goes here
	fmt.Println("Hello, Pixel!")

	// create a window
	// window
	cfg := pixelgl.WindowConfig{
		Title:     "Hello, Pixel!",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(false)

	// make win blue
	win.Clear(colornames.Skyblue)

	// load sprites
	load := asset.Newload(os.DirFS("./"))
	spriteSheet, err := load.SpriteSheet("packed.json")
	if err != nil {
		panic(err)
	}

	// 1st hog
	hogSprites1, err := spriteSheet.Get("hedge-hog-mv-1.png")
	if err != nil {
		panic(err)
	}
	hogPosition1 := win.Bounds().Center()

	// 2nd hog
	hogSprites2, err := spriteSheet.Get("hedge-hog-mv-2.png")
	if err != nil {
		panic(err)
	}
	hogPosition2 := win.Bounds().Center()

	// create a person list
	var hogs []*Person

	hogs = append(hogs, NewPerson(hogSprites1, hogPosition1, KeyBind{
		Up:    pixelgl.KeyUp,
		Down:  pixelgl.KeyDown,
		Left:  pixelgl.KeyLeft,
		Right: pixelgl.KeyRight,
	}))

	hogs = append(hogs, NewPerson(hogSprites2, hogPosition2, KeyBind{
		Up:    pixelgl.KeyW,
		Down:  pixelgl.KeyS,
		Left:  pixelgl.KeyA,
		Right: pixelgl.KeyD,
	}))

	// game loop
	for !win.JustPressed(pixelgl.KeyEscape) {

		// make win blue
		win.Clear(colornames.Skyblue)

		// handle input
		for i := range hogs {
			hogs[i].HandleInput(win)
		}
		// draw the sprites
		for i := range hogs {
			hogs[i].Draw(win)
		}

		// update the window
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

type KeyBind struct {
	Up, Down, Left, Right pixelgl.Button
}

type Person struct {
	Sprite  *pixel.Sprite
	Pos     pixel.Vec
	KeyBind KeyBind
}

func NewPerson(sprite *pixel.Sprite, pos pixel.Vec, keyBind KeyBind) *Person {
	return &Person{
		Sprite:  sprite,
		Pos:     pos,
		KeyBind: keyBind,
	}
}

func (p *Person) Draw(win *pixelgl.Window) {
	p.Sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(p.Pos))
}

func (p *Person) HandleInput(win *pixelgl.Window) {
	if win.Pressed(p.KeyBind.Up) {
		p.Pos.Y += 2.0
	}
	if win.Pressed(p.KeyBind.Down) {
		p.Pos.Y -= 2.0
	}
	if win.Pressed(p.KeyBind.Left) {
		p.Pos.X -= 2.0
	}
	if win.Pressed(p.KeyBind.Right) {
		p.Pos.X += 2.0
	}
}
