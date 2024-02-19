package main

//go:generate packer --input images --stats

import (
	"fmt"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/mon7792/go-mmo/engine/asset"
	"github.com/mon7792/go-mmo/engine/render"
	"github.com/mon7792/go-mmo/engine/tilemap"
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

	// tilemap
	tileSize := 16
	mapSize := 100
	tiles := make([][]tilemap.Tile, mapSize, mapSize)

	grassSprite, err := spriteSheet.Get("grass.png")
	if err != nil {
		panic(err)
	}

	for x := range tiles {
		tiles[x] = make([]tilemap.Tile, mapSize, mapSize)
		for y := range tiles[x] {
			tiles[x][y] = tilemap.Tile{
				Type:   0,
				Sprite: grassSprite,
			}
		}
	}

	batch := pixel.NewBatch(&pixel.TrianglesData{}, grassSprite.Picture())
	tMap := tilemap.New(tiles, batch, tileSize)
	// rebatch
	tMap.Rebatch()

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
	// camera init
	camera := render.NewCamera(win, 0, 0)

	zoomSpeed := 1.0
	// game loop
	for !win.JustPressed(pixelgl.KeyEscape) {

		// make win blue
		win.Clear(colornames.Skyblue)

		// scroll
		scroll := win.MouseScroll()
		if scroll.Y != 0 {
			camera.Zoom += zoomSpeed * scroll.Y
		}

		// handle input
		for i := range hogs {
			hogs[i].HandleInput(win)
		}

		// camera
		camera.Position = hogs[0].Pos
		camera.Update()

		win.SetMatrix(camera.Matrix())

		// draw the sprites

		// draw the tilemap
		tMap.Draw(win)

		// draw the person
		for i := range hogs {
			hogs[i].Draw(win)
		}

		// camera
		win.SetMatrix(pixel.IM)

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
