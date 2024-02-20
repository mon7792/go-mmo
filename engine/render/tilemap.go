package render

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mon7792/go-mmo/engine/asset"
	"github.com/mon7792/go-mmo/engine/tilemap"
)

type TilemapRender struct {
	spritesheet  *asset.SpriteSheet
	batch        *pixel.Batch
	tileToSprite map[tilemap.TileType]*pixel.Sprite
}

func NewTilemapRender(spritesheet *asset.SpriteSheet, tileToSprite map[tilemap.TileType]*pixel.Sprite) *TilemapRender {
	return &TilemapRender{
		spritesheet:  spritesheet,
		batch:        pixel.NewBatch(&pixel.TrianglesData{}, spritesheet.Picture()),
		tileToSprite: tileToSprite,
	}
}

func (tr *TilemapRender) Clear() {
	tr.batch.Clear()
}

func (tr *TilemapRender) Batch(tm *tilemap.TileMap) {
	for x := 0; x < tm.Width(); x++ {
		for y := 0; y < tm.Height(); y++ {
			tile, ok := tm.Get(x, y)
			if !ok {
				continue
			}
			pos := pixel.V(float64(x*tm.TileSize), float64(y*tm.TileSize))

			sprite, ok := tr.tileToSprite[tile.Type]
			if !ok {
				panic("unknownSprite: sprite not found")
			}

			mat := pixel.IM.Moved(pos)

			sprite.Draw(tr.batch, mat)
		}
	}
}

func (tr *TilemapRender) Draw(win *pixelgl.Window) {
	tr.batch.Draw(win)
}
