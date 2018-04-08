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

type FontProps struct {
	StrokeWidth float32
	StrokeEdge  float32
	Color       graphics.Color
}

var (
	FontPropLarge = FontProps{0.8, 0.02, graphics.Color{1, 1, 1, 1}}
	FontPropSmall = FontProps{0.5, 0.1, graphics.Color{1, 1, 1, 1}}
)

func NewFontFromFiles(bmpath, texpath string) *Font {
	f := &Font{}
	f.bm = NewBmFontFromFile(bmpath)
	f.tx = graphics.NewTextureFromFile(texpath)
	return f
}

const CharVertices = 12

func charQuad(offsetX, offsetY, width, height float32) []float32 {
	q := [CharVertices]float32{
		offsetX, offsetY + height,         // bl
		offsetX + width, offsetY + height, // br
		offsetX, offsetY,                  // tl
		offsetX, offsetY,                  // tl
		offsetX + width, offsetY + height, // br
		offsetX + width, offsetY,          // tr
	}
	return q[:]
}

func (f *Font) RenderText(
	txt string,
	position mgl32.Vec3,
	size mgl32.Vec2,
	fp FontProps,
) *graphics.Primitive2D {
	var (
		vnum = len(txt)*CharVertices
		vertices = make([]float32, vnum)
		uvCoords = make([]float32, vnum)
		idx = 0
		cursorX float32 = 0
		cursorY float32 = 0
		scale = 1.0/float32(f.bm.lineHeight)
		uvscalex = 1.0/float32(f.bm.pageWidth)
		uvscaley = 1.0/float32(f.bm.pageHeight)
		lastChar int32 = 0
	)

	for _, char := range txt {
		if char == 0x0a {
			cursorX = 0
			cursorY += 1
			lastChar = 0
			continue
		}
		bmc, ok := f.bm.Characters[string(char)]
		if !ok {
			log.Printf(
				"ERR: char %v (%v) not found in font map",
				string(char), char,
			)
			continue
		}

		kerning := 0
		kl, ok := f.bm.Kernings[lastChar]
		if ok {
			kc, ok := kl[char]
			if ok {
				kerning = kc
			}
		}

		copy(
			vertices[idx:],
			charQuad(
				cursorX+float32(bmc.offsetX+kerning)*scale,
				cursorY+float32(bmc.offsetY)*scale,
				float32(bmc.width)*scale,
				float32(bmc.height)*scale,
			),
		)
		copy(
			uvCoords[idx:],
			charQuad(
				float32(bmc.x)*uvscalex,
				float32(bmc.y)*uvscaley,
				float32(bmc.width)*uvscalex,
				float32(bmc.height)*uvscaley,
			),
		)
		idx += CharVertices
		cursorX += float32(bmc.advanceX)*scale
		lastChar = char
	}

	shaderProgram := graphics.NewShaderProgram(
		graphics.VertexShaderPrimitive2D,
		"",
		FragmentDistanceFieldFont,
	)

	// TODO: this should be done in draw/drawInBatch?
	//shaderProgram.SetUniformV4fv("textColor", &fp.Color)
	//shaderProgram.SetUniformV2f("widthEdge", fp.StrokeWidth, fp.StrokeEdge)

	return graphics.NewTriangles(vertices, uvCoords, f.tx, position, size, shaderProgram)
}

var (
	FragmentDistanceFieldFont = `
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
