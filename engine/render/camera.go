package render

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Camera struct {
	win      *pixelgl.Window
	Position pixel.Vec
	Zoom     float64
	Mat      pixel.Matrix
}

func NewCamera(win *pixelgl.Window, x, y float64) *Camera {
	return &Camera{
		win:      win,
		Position: pixel.V(x, y),
		Zoom:     1.0,
		Mat:      pixel.IM,
	}
}

func (c *Camera) Update() {
	screenCenter := c.win.Bounds().Center()
	movePos := pixel.V(
		math.Floor(-c.Position.X),
		math.Floor(-c.Position.Y)).Add(screenCenter)

	c.Mat = pixel.IM.Moved(movePos).Scaled(screenCenter, c.Zoom)
}

func (c *Camera) Matrix() pixel.Matrix {
	return c.Mat
}
