package tilemap

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type TileType int8

type Tile struct {
	Type   TileType
	Sprite *pixel.Sprite
}

type TileMap struct {
	TileSize int // in pixel
	Tiles    [][]Tile
	batch    *pixel.Batch
}

func New(tiles [][]Tile, batch *pixel.Batch, tileSize int) *TileMap {
	return &TileMap{
		TileSize: tileSize,
		Tiles:    tiles,
		batch:    batch,
	}
}

func (tm *TileMap) Rebatch() {
	for x := range tm.Tiles {
		for y := range tm.Tiles[x] {
			tile := tm.Tiles[x][y]
			pos := pixel.V(float64(x*tm.TileSize), float64(y*tm.TileSize))

			mat := pixel.IM.Moved(pos)
			tile.Sprite.Draw(tm.batch, mat)
		}
	}
}

func (tm *TileMap) Draw(win *pixelgl.Window) {
	tm.batch.Draw(win)
}
