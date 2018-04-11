package graphics

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	maxPostProcessingSteps = 16
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
	clearColor           Color
	primitivesToDraw     map[uint32][]Drawable
	postProcessingSteps  []*PostProcessingStep
}

func (c *Context) BeginRendering() {
	if c.postProcessingSteps != nil && len(c.postProcessingSteps) > 0 {
		c.postProcessingSteps[0].frameBuffer.Bind()
	}
	c.Clear()
}

func (c *Context) EndRendering() {
	c.renderDrawableList()
	c.eraseDrawableList()

	if c.postProcessingSteps != nil && len(c.postProcessingSteps) > 0 {
		if len(c.postProcessingSteps) > 1 {
			for i := 1; i < len(c.postProcessingSteps); i++ {
				c.postProcessingSteps[i].frameBuffer.Bind()
				c.postProcessingSteps[i-1].quad.Draw(c)
			}

			c.postProcessingSteps[len(c.postProcessingSteps)-1].frameBuffer.Unbind()
			c.postProcessingSteps[len(c.postProcessingSteps)-1].quad.Draw(c)
		}
	}
}

func (c *Context) Clear() {
	gl.ClearColor(c.clearColor.G(), c.clearColor.G(), c.clearColor.B(), c.clearColor.A())
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}

func (c *Context) SetClearColor(color Color) {
	c.clearColor = color
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

func (c *Context) renderDrawableList() {
	for _, v := range c.primitivesToDraw {
		for _, drawable := range v {
			c.BindTexture(drawable.Texture())
			shader := drawable.Shader()
			c.BindShader(shader)
			// TODO this should be done only once per frame via uniform buffers
			shader.SetUniform("mProjection", &c.projectionMatrix)
			drawable.(Drawable).drawInBatch(c)
		}
	}
}

func (c *Context) eraseDrawableList() {
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

func (c *Context) AddPostProcessingStep(step *PostProcessingStep) {
	if c.postProcessingSteps == nil {
		c.postProcessingSteps = make([]*PostProcessingStep, maxPostProcessingSteps)
	}
	c.postProcessingSteps = append(c.postProcessingSteps, step)
}
