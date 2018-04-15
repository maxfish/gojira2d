package ui

import (
	g "github.com/maxfish/gojira2d/graphics"
)

// Font structure contains BmFont metadata and the texture
type Font struct {
	bm *BmFont
	tx *g.Texture
}

// FontRegistry is a dictionary of loaded fonts
var FontRegistry = make(map[string]*Font)

// NewFontFromFiles create Font structure from metadata and texture files
func NewFontFromFiles(name, bmpath, texpath string) *Font {
	if f, ok := FontRegistry[name]; ok {
		return f
	}

	f := &Font{}
	f.bm = NewBmFontFromFile(bmpath)
	f.tx = g.NewTextureFromFile(texpath)
	FontRegistry[name] = f
	return f
}
