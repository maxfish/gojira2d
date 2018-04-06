package graphics

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Drawable interface {
	Texture() (*Texture)
	Shader() (*ShaderProgram)
	Draw(context *Context)
	drawInBatch(context *Context)
}

type Context struct {
	projectionMatrix     mgl32.Mat4
	viewMatrix           mgl32.Mat4
	currentTexture       *Texture
	currentShaderProgram *ShaderProgram
	primitivesToDraw     map[uint32][]Drawable
}

func (c *Context) Clear() {
	gl.ClearColor(0.3, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}

func (c *Context) enqueueForDrawing(drawable Drawable) {
	if c.primitivesToDraw == nil {
		c.primitivesToDraw = make(map[uint32][]Drawable)
	}

	texture := drawable.Texture()
	var textureId uint32 = 0
	if texture != nil {
		textureId = texture.id
	}
	// Groups the primitives by their texture's id
	c.primitivesToDraw[textureId] = append(c.primitivesToDraw[textureId], drawable)
}

func (c *Context) RenderDrawableList() {
	for _, v := range c.primitivesToDraw {
		for _, drawable := range v {
			c.BindTexture(drawable.Texture())
			shader := drawable.Shader()
			c.BindShader(shader)
			// TODO this should be done only once per frame via uniform buffers
			shader.SetUniform4fv("mProjection", &c.projectionMatrix)
			drawable.(Drawable).drawInBatch(c)
		}
	}
}

func (c *Context) EraseDrawableList() {
	c.primitivesToDraw = make(map[uint32][]Drawable)
}

func (c *Context) BindTexture(texture *Texture) {
	if texture == nil {
		return
	}
	if c.currentTexture == nil || texture.id != c.currentTexture.id {
		gl.BindTexture(gl.TEXTURE_2D, texture.id)
		c.currentTexture = texture
	}
}

func (c *Context) BindShader(shader *ShaderProgram) {
	if c.currentShaderProgram == nil || shader.id != c.currentShaderProgram.id {
		gl.UseProgram(shader.id)
		c.currentShaderProgram = shader
	}
}

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
