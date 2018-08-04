package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Drawable ...
type Drawable interface {
	Texture() *Texture
	Shader() *ShaderProgram
	Draw(context *Context)
	DrawInBatch(context *Context)
}

// Context ...
type Context struct {
	projectionMatrix     mgl32.Mat4
	viewMatrix           mgl32.Mat4
	currentTexture       *Texture
	currentShaderProgram *ShaderProgram
}

// BindTexture sets texture to be current texture if it isn't already
func (c *Context) BindTexture(texture *Texture) {
	if texture == nil {
		return
	}
	if c.currentTexture == nil || texture.id != c.currentTexture.id {
		gl.BindTexture(gl.TEXTURE_2D, texture.id)
		c.currentTexture = texture
	}
}

// BindShader sets shader to be current chader if it isn't already
func (c *Context) BindShader(shader *ShaderProgram) {
	if c.currentShaderProgram == nil || shader.id != c.currentShaderProgram.id {
		gl.UseProgram(shader.id)
		c.currentShaderProgram = shader
	}
}

// SetOrtho2DProjection sets up orthogonal projection matrix
func (c *Context) SetOrtho2DProjection(windowWidth int, windowHeight int, screenScale float32, centered bool) {
	var left, right, top, bottom float32
	if centered {
		// 0,0 is placed at the center of the window
		halfWidth := float32(windowWidth) / 2 / screenScale
		halfHeight := float32(windowHeight) / 2 / screenScale
		left = -halfWidth
		right = halfWidth
		top = halfHeight
		bottom = -halfHeight
	} else {
		left = 0
		right = float32(windowWidth)
		top = float32(windowHeight)
		bottom = 0
	}
	c.projectionMatrix = mgl32.Ortho(left, right, top, bottom, 1, -1)
}
