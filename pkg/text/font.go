package text

import (
	"log"
	"gojira2d/pkg/graphics"
	"github.com/go-gl/mathgl/mgl32"
)

type Font struct {
	bm *BmFont
	tx *graphics.Texture
}

func NewFontFromFiles(bmpath, texpath string) *Font {
	f := &Font{}
	f.bm = NewBmFontFromFile(bmpath)
	f.tx = graphics.NewTextureFromFile(texpath)
	return f
}

func (f *Font) RenderText(txt string) *graphics.Primitive2D {
	vertices := make([]float32, len(txt)*8+(len(txt)-1)*4)
	uvCoords := make([]float32, len(txt)*8+(len(txt)-1)*4)

	idx := 0
	end := len(txt)
	for ic, char := range txt {
		if bmc, ok := f.bm.Characters[string(char)]; ok {
			log.Printf("%+v", bmc)
			copy(vertices[idx:], []float32{
				float32(ic),1,
				float32(ic+1),1,
				float32(ic),0,
				float32(ic+1),0,
			})
			copy(uvCoords[idx:], []float32{0,1,1,1,0,0,1,0,})
			idx += 8

			if ic < end {
				copy(vertices[idx:], []float32{1,1,1,1})
				copy(uvCoords[idx:], []float32{1,1,1,1})
				idx += 4
			}
		} else {
			log.Printf(
				"ERR: char %v (%v) not found in font map",
				string(char), char,
			)
		}
	}

	p := graphics.NewTriangleStrip(
		vertices,
		uvCoords,
		f.tx,
		mgl32.Vec3{0,0,0},
		mgl32.Vec2{80,80},
	)
	p.SetAnchorToCenter()
	return p
}
