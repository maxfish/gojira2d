package ui

import (
	g "gojira2d/pkg/graphics"
)

// Font structure contains BmFont metadata and the texture
type Font struct {
	bm *BmFont
	tx *g.Texture
}

// FontProps holds various attributes that control text rendering
type FontProps struct {
	StrokeWidth float32
	StrokeEdge  float32
	Color       g.Color
}

var (
	// FontPropLarge should be used for large text
	FontPropLarge = FontProps{0.8, 0.02, g.Color{1, 1, 1, 1}}

	// FontPropSmall should be used for small text
	FontPropSmall = FontProps{0.5, 0.1, g.Color{1, 1, 1, 1}}
)

// NewFontFromFiles create Font structure from metadata and texture files
func NewFontFromFiles(bmpath, texpath string) *Font {
	f := &Font{}
	f.bm = NewBmFontFromFile(bmpath)
	f.tx = g.NewTextureFromFile(texpath)
	return f
}
