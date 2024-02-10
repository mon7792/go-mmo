package asset

import (
	"encoding/json"
	"errors"
	"image"
	_ "image/png"
	"io"
	"io/fs"

	"github.com/faiface/pixel"
	"github.com/unitoftime/packer"
)

var errSpriteNotFound = errors.New("sprite not found")

type Load interface {
	Open(path string) (fs.File, error)
	Sprites(path string) (*pixel.Sprite, error)
	SpriteSheet(path string) (*SpriteSheet, error)
}

type load struct {
	filesystem fs.FS
}

func Newload(filesystem fs.FS) Load {
	return &load{
		filesystem: filesystem,
	}
}

func (l *load) Open(path string) (fs.File, error) {
	return l.filesystem.Open(path)
}

func (l *load) Json(path string, data interface{}) (err error) {
	var f fs.File
	f, err = l.filesystem.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if fsErr := f.Close(); fsErr != nil {
			err = fsErr
		}
	}()

	var jsonData []byte
	jsonData, err = io.ReadAll(f)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, &data)
}

func (l *load) Image(path string) (img image.Image, err error) {
	var f fs.File
	f, err = l.filesystem.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if fsErr := f.Close(); fsErr != nil {
			err = fsErr
		}
	}()
	img, _, err = image.Decode(f)
	if err != nil {
		return nil, err
	}
	return
}

func (l *load) Sprites(path string) (ps *pixel.Sprite, err error) {
	var f fs.File
	f, err = l.filesystem.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if fsErr := f.Close(); fsErr != nil {
			err = fsErr
		}
	}()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	pic := pixel.PictureDataFromImage(img)
	return pixel.NewSprite(pic, pic.Bounds()), nil
}

func (l *load) SpriteSheet(path string) (*SpriteSheet, error) {
	// load the json
	serializedSpriteSheet := packer.SerializedSpritesheet{}
	err := l.Json(path, &serializedSpriteSheet)
	if err != nil {
		return nil, err
	}
	// load the image
	img, err := l.Image(serializedSpriteSheet.ImageName)
	if err != nil {
		return nil, err
	}
	pic := pixel.PictureDataFromImage(img)
	// create the sprites
	bounds := pic.Bounds()
	lookup := make(map[string]*pixel.Sprite)
	for k, v := range serializedSpriteSheet.Frames {
		rect := pixel.R(
			v.Frame.X,
			bounds.H()-v.Frame.Y,
			v.Frame.X+v.Frame.W,
			bounds.H()-(v.Frame.Y+v.Frame.H)).Norm()

		lookup[k] = pixel.NewSprite(pic, rect)
	}
	return NewSpriteSheet(pic, lookup), nil
}

// SpriteSheet is a collection of sprites from a single image
type SpriteSheet struct {
	picture pixel.Picture
	lookup  map[string]*pixel.Sprite
}

func NewSpriteSheet(picture pixel.Picture, lookup map[string]*pixel.Sprite) *SpriteSheet {
	return &SpriteSheet{
		picture: picture,
		lookup:  lookup,
	}
}

func (s *SpriteSheet) Get(name string) (*pixel.Sprite, error) {
	sprite, ok := s.lookup[name]
	if !ok {
		return nil, errSpriteNotFound
	}
	return sprite, nil
}
