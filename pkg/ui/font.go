package ui

import (
	g "gojira2d/pkg/graphics"
)

// Font structure contains BmFont metadata and the texture
type Font struct {
	bm *BmFont
	tx *g.Texture
}

// NewFontFromFiles create Font structure from metadata and texture files
func NewFontFromFiles(bmpath, texpath string) *Font {
	f := &Font{}
	f.bm = NewBmFontFromFile(bmpath)
	f.tx = g.NewTextureFromFile(texpath)
	return f
}
