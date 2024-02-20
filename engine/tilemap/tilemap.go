package tilemap

import (
	"github.com/faiface/pixel"
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

func (tm *TileMap) Width() int {
	return len(tm.Tiles)
}

func (tm *TileMap) Height() int {
	return len(tm.Tiles[0])
}

func (tm *TileMap) Get(x, y int) (Tile, bool) {
	if x < 0 || x >= tm.Width() || y < 0 || y >= tm.Height() {
		return Tile{}, false
	}
	return tm.Tiles[x][y], true
}
