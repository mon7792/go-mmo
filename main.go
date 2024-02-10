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

	hogSprites, err := spriteSheet.Get("hedge-hog-mv-1.png")
	if err != nil {
		panic(err)
	}
	hogPosition := win.Bounds().Center()

	// game loop
	for !win.JustPressed(pixelgl.KeyEscape) {

		// make win blue
		win.Clear(colornames.Skyblue)

		// check for inputs
		if win.Pressed(pixelgl.KeyLeft) {
			hogPosition.X -= 2.0
		}
		if win.Pressed(pixelgl.KeyRight) {
			hogPosition.X += 2.0
		}
		if win.Pressed(pixelgl.KeyUp) {
			hogPosition.Y += 2.0
		}
		if win.Pressed(pixelgl.KeyDown) {
			hogPosition.Y -= 2.0
		}

		// draw the sprites
		hogSprites.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(hogPosition))
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
