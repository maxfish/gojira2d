package ui

import (
	"gojira2d/pkg/graphics"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Text is a UI element that just renders a string
type Text struct {
	drawable  *graphics.Primitive2D
	position  mgl32.Vec3
	size      mgl32.Vec2
	fontProps FontProps
	text      string
	font      *Font
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

func charQuads(txt string, font *Font) ([]float32, []float32) {
	var (
		vnum     = len(txt) * charVertices
		idx      = 0
		cursorX  float32
		cursorY  float32
		lastChar int32
	)

	vertices := make([]float32, vnum)
	uvCoords := make([]float32, vnum)

	for _, char := range txt {
		if char == 0x0a {
			cursorX = 0
			cursorY++
			lastChar = 0
			continue
		}
		bmc, ok := font.bm.Characters[char]
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

	return vertices, uvCoords
}

var textShaderProgram *graphics.ShaderProgram

// NewText creates a Primitive2D with character quads for given string
func (f *Font) NewText(
	txt string,
	position mgl32.Vec3,
	size mgl32.Vec2,
	fp FontProps,
) *Text {
	if textShaderProgram == nil {
		textShaderProgram = graphics.NewShaderProgram(
			graphics.VertexShaderPrimitive2D, "", fragmentDistanceFieldFont,
		)
	}

	text := &Text{}
	charVertices, charUVCoords := charQuads(txt, f)

	text.drawable = graphics.NewTriangles(
		charVertices, charUVCoords, f.tx, position, size, textShaderProgram)
	text.position = position
	text.fontProps = fp
	text.size = size
	text.text = txt
	text.font = f

	return text
}

// SetText changes the rendered string
func (t *Text) SetText(txt string) {
	charVertices, charUVCoords := charQuads(txt, t.font)
	t.drawable.SetVertices(charVertices)
	t.drawable.SetUVCoords(charUVCoords)
}

// EnqueueForDrawing see Drawable.EnqueueForDrawing
func (t *Text) EnqueueForDrawing(context *graphics.Context) {
	context.EnqueueForDrawing(t)
}

// SetUniforms uploads relevant uniforms
func (t *Text) SetUniforms() {
	shaderProgram := t.Shader()
	shaderProgram.SetUniform("textColor", &t.fontProps.Color)
	shaderProgram.SetUniform(
		"widthEdge",
		&mgl32.Vec2{t.fontProps.StrokeWidth, t.fontProps.StrokeEdge},
	)
}

// Drawable implementation

// Texture returns drawable texture
func (t *Text) Texture() *graphics.Texture {
	return t.drawable.Texture()
}

// Shader returns shader program
func (t *Text) Shader() *graphics.ShaderProgram {
	return t.drawable.Shader()
}

// Draw runs all the necessary routines to make drawable appear on screen
func (t *Text) Draw(context *graphics.Context) {
	shaderProgram := t.Shader()
	gl.UseProgram(shaderProgram.Id())
	t.SetUniforms()
	t.drawable.Draw(context)
}

// DrawInBatch is like Draw() but without setting up texture and shader
func (t *Text) DrawInBatch(context *graphics.Context) {
	t.SetUniforms()
	t.drawable.DrawInBatch(context)
}

// Drawable end

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
          float alpha = 1.0 - smoothstep(widthEdge.x, widthEdge.x+widthEdge.y, distance);
          color = vec4(vec3(textColor),alpha);
        }
        ` + "\x00"
)
