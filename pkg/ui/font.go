package ui

import (
	"gojira2d/pkg/graphics"
	"log"

	"github.com/go-gl/mathgl/mgl32"
)

// Font structure contains BmFont metadata and the texture
type Font struct {
	bm *BmFont
	tx *graphics.Texture
}

// FontProps holds various attributes that control text rendering
type FontProps struct {
	StrokeWidth float32
	StrokeEdge  float32
	Color       graphics.Color
}

var (
	// FontPropLarge should be used for large text
	FontPropLarge = FontProps{0.8, 0.02, graphics.Color{1, 1, 1, 1}}

	// FontPropSmall should be used for small text
	FontPropSmall = FontProps{0.5, 0.1, graphics.Color{1, 1, 1, 1}}
)

// NewFontFromFiles create Font structure from metadata and texture files
func NewFontFromFiles(bmpath, texpath string) *Font {
	f := &Font{}
	f.bm = NewBmFontFromFile(bmpath)
	f.tx = graphics.NewTextureFromFile(texpath)
	return f
}

const charVertices = 12

func charQuad(offsetX, offsetY, width, height float32) []float32 {
	q := [charVertices]float32{
		offsetX, offsetY + height, // bl
		offsetX + width, offsetY + height, // br
		offsetX, offsetY, // tl
		offsetX, offsetY, // tl
		offsetX + width, offsetY + height, // br
		offsetX + width, offsetY, // tr
	}
	return q[:]
}

// RenderText creates a Primitive2D with character quads for given string
func (f *Font) RenderText(
	txt string,
	position mgl32.Vec3,
	size mgl32.Vec2,
	fp FontProps,
) *graphics.Primitive2D {
	var (
		vnum     = len(txt) * charVertices
		vertices = make([]float32, vnum)
		uvCoords = make([]float32, vnum)
		idx      = 0
		cursorX  float32
		cursorY  float32
		lastChar int32
	)

	for _, char := range txt {
		if char == 0x0a {
			cursorX = 0
			cursorY++
			lastChar = 0
			continue
		}
		bmc, ok := f.bm.Characters[char]
		if !ok {
			log.Printf(
				"ERR: char %v (%v) not found in font map",
				string(char), char,
			)
			continue
		}

		kerning, ok := bmc.f32kernings[lastChar]
		if !ok {
			kerning = 0
		}

		copy(
			vertices[idx:],
			charQuad(
				cursorX+bmc.f32offsetX+kerning,
				cursorY+bmc.f32offsetY,
				bmc.f32lineWidth,
				bmc.f32lineHeight,
			),
		)
		copy(
			uvCoords[idx:],
			charQuad(
				bmc.f32x,
				bmc.f32y,
				bmc.f32width,
				bmc.f32height,
			),
		)
		idx += charVertices
		cursorX += bmc.f32advanceX
		lastChar = char
	}

	shaderProgram := graphics.NewShaderProgram(
		graphics.VertexShaderPrimitive2D, "", fragmentDistanceFieldFont,
	)

	// TODO: this should be done in draw/drawInBatch?
	//shaderProgram.SetUniformV4fv("textColor", &fp.Color)
	//shaderProgram.SetUniformV2f("widthEdge", fp.StrokeWidth, fp.StrokeEdge)

	return graphics.NewTriangles(vertices, uvCoords, f.tx, position, size, shaderProgram)
}

var (
	fragmentDistanceFieldFont = `
        #version 410 core

        in vec2 uv_out;
        out vec4 color;

        uniform sampler2D tex;
		uniform vec4 textColor;
		uniform vec2 widthEdge;

        void main() { 
			float distance = 1.0 - texture(tex, uv_out).a;
			vec2 widthEdge = vec2(0.5,0.1);
			float alpha = 1.0 - smoothstep(widthEdge.x, widthEdge.x+widthEdge.y, distance);
            color = vec4(0,0,0,alpha);
        }
        ` + "\x00"
)
