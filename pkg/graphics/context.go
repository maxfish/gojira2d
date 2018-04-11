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
	clearColor           Color
	primitivesToDraw     map[uint32][]Drawable
}

// Clear clears the screen using Context.clearColor
func (c *Context) Clear() {
	gl.ClearColor(c.clearColor.R(), c.clearColor.G(), c.clearColor.B(), c.clearColor.A())
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	// Do we need to do this every Clear?
	// gl.Enable(gl.BLEND)
	// gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}

// SetClearColor changes OpenGL background clear color
func (c *Context) SetClearColor(color Color) {
	c.clearColor = color
}

// EnqueueForDrawing adds a drawable to drawing list
func (c *Context) EnqueueForDrawing(drawable Drawable) {
	if c.primitivesToDraw == nil {
		c.primitivesToDraw = make(map[uint32][]Drawable)
	}

	texture := drawable.Texture()
	var textureID uint32
	if texture != nil {
		textureID = texture.id
	}
	// Groups the primitives by their texture's id
	c.primitivesToDraw[textureID] = append(c.primitivesToDraw[textureID], drawable)
}

// RenderDrawableList binds shader and texture for each primitive and calls DrawInBatch
func (c *Context) RenderDrawableList() {
	for _, v := range c.primitivesToDraw {
		for _, drawable := range v {
			c.BindTexture(drawable.Texture())
			shader := drawable.Shader()
			c.BindShader(shader)
			// TODO this should be done only once per frame via uniform buffers
			shader.SetUniform("mProjection", &c.projectionMatrix)
			drawable.DrawInBatch(c)
		}
	}
}

// EraseDrawableList resets primitivesToDraw to empty list
func (c *Context) EraseDrawableList() {
	c.primitivesToDraw = make(map[uint32][]Drawable)
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
