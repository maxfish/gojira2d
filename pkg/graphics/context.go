package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl64"
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
	Camera2D             *Camera2D
	viewMatrix           mgl64.Mat4
	currentTexture       *Texture
	currentShaderProgram *ShaderProgram
}

// BindTexture sets texture to be current texture if it isn't already
func (c *Context) BindTexture(texture *Texture) {
	if texture == nil {
		gl.BindTexture(gl.TEXTURE_2D, 0)
		c.currentTexture = nil
		return
	}

	gl.BindTexture(gl.TEXTURE_2D, texture.id)
	c.currentTexture = texture
}

// BindShader sets shader to be current shader if it isn't already
func (c *Context) BindShader(shader *ShaderProgram) {
	if c.currentShaderProgram == nil || shader.id != c.currentShaderProgram.id {
		gl.UseProgram(shader.id)
		c.currentShaderProgram = shader
	}
}
