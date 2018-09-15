package utils

import (
	"fmt"
	"github.com/maxfish/gojira2d/pkg/graphics"
)

// TextureStorage handles multiple textures, each one with its own id
type TextureStorage struct {
	path     string
	textures map[string]*graphics.Texture
}

// MakeTextureStorage creates a new initialized instance
func MakeTextureStorage(path string) *TextureStorage {
	return &TextureStorage{
		path:     path,
		textures: make(map[string]*graphics.Texture),
	}
}

// LoadTexture loads a texture from file and assigns an id to it
func (ts *TextureStorage) LoadTexture(fileName string, textureId string) {
	fullPath := fmt.Sprintf("%s/%s", ts.path, fileName)
	t := graphics.NewTextureFromFile(fullPath)
	ts.textures[textureId] = t
}

// TextureForId returns the texture associated with the id
func (ts *TextureStorage) TextureForId(textureId string) *graphics.Texture {
	return ts.textures[textureId]
}
