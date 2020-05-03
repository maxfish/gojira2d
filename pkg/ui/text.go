package ui

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/maxfish/gojira2d/pkg/graphics"
)

// Text is a UI element that just renders a string
type Text struct {
	drawable *graphics.Primitive2D
	position mgl64.Vec3
	size     mgl64.Vec2
	color    graphics.Color
	text     string
	font     *Font
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

func (t *Text) makeNewQuads() ([]float32, []float32) {
	var (
		vnum     = len(t.text) * charVertices
		idx      = 0
		cursorX  float32
		cursorY  float32
		lastChar int32
	)

	vertices := make([]float32, vnum)
	uvCoords := make([]float32, vnum)
	bmf := t.font.bm

	for _, char := range t.text {
		if char == 0x0a {
			cursorX = 0
			cursorY += float32(t.font.bm.lineHeight) + float32(t.font.bm.padding[1])
			lastChar = 0
			continue
		}
		bmc, ok := t.font.bm.Characters[char]
		if !ok {
			fmt.Printf("ERR: char %v (%v) not found in font map\n", string(char), char)
			continue
		}

		kerning, ok := bmc.kernings[lastChar]
		if !ok {
			kerning = 0
		}

		var scale = 1 / float32(bmf.lineHeight)

		deltaX := float32(bmc.offsetX+kerning+bmf.padding[2]) * scale
		deltaY := float32(bmc.offsetY+bmf.padding[0]) * scale
		width := float32(bmc.width) * scale
		height := float32(bmc.height) * scale

		copy(
			vertices[idx:],
			charQuad(
				cursorX+deltaX,
				cursorY+deltaY,
				width,
				height,
			),
		)
		copy(
			uvCoords[idx:],
			charQuad(
				float32(bmc.x)/float32(bmf.pageWidth),
				float32(bmc.y)/float32(bmf.pageHeight),
				float32(bmc.width)/float32(bmf.pageWidth),
				float32(bmc.height)/float32(bmf.pageHeight),
			),
		)
		idx += charVertices
		cursorX += float32(bmc.advanceX+bmf.padding[1]) * scale
		lastChar = char
	}

	return vertices, uvCoords
}

var textShaderProgram *graphics.ShaderProgram

// NewText creates a Primitive2D with character quads for given string
func NewText(
	txt string,
	font *Font,
	position mgl64.Vec3,
	size mgl64.Vec2,
	color graphics.Color,
) *Text {
	if textShaderProgram == nil {
		textShaderProgram = graphics.NewShaderProgram(
			graphics.VertexShaderBase, "", fragmentDistanceFieldFont,
		)
	}

	t := &Text{}
	t.text = txt
	t.position = position
	t.color = color
	t.size = size
	t.text = txt
	t.font = font

	charVertices, charUVCoords := t.makeNewQuads()
	t.drawable = graphics.NewTriangles(
		charVertices, charUVCoords, font.tx, position, size, textShaderProgram)

	return t
}

func (t *Text) uploadNewQuads() {
	charVertices, charUVCoords := t.makeNewQuads()
	t.drawable.SetVertices(charVertices)
	t.drawable.SetUVCoords(charUVCoords)
}

// SetText changes the rendered string and uploads new vertices/coordinates
func (t *Text) SetText(txt string) {
	if t.text != txt {
		t.text = txt
		t.uploadNewQuads()
	}
}

// SetColor ...
func (t *Text) SetColor(color graphics.Color) {
	t.color = color
}

// SetUniforms uploads relevant uniforms
func (t *Text) SetUniforms() {
	shaderProgram := t.Shader()
	shaderProgram.SetUniform("textColor", &t.color)
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
	gl.UseProgram(shaderProgram.ID())
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
		
		void main() {
			float dist = texture(tex, uv_out).a;
			float width = fwidth(dist);
			float alpha = smoothstep(0.5-width, 0.5+width, dist);
			color = vec4(vec3(textColor), alpha*textColor.a);
		}
        ` + "\x00"
)
