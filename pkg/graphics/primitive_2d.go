package graphics

import (
	"github.com/maxfish/pkg/utils"

	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const FLOAT32_SIZE = 4

type ModelMatrix struct {
	mgl32.Mat4
	size        mgl32.Mat4
	translation mgl32.Mat4
	rotation    mgl32.Mat4
	scale       mgl32.Mat4
	anchor      mgl32.Mat4
	dirty       bool
}

type Primitive2D struct {
	Primitive
	position    mgl32.Vec3
	scale       mgl32.Vec2
	size        mgl32.Vec2
	anchor      mgl32.Vec2
	angle       float32
	flipX       bool
	flipY       bool
	color       Color
	modelMatrix ModelMatrix
}

func (p *Primitive2D) SetPosition(position mgl32.Vec3) {
	p.position = position
	p.modelMatrix.translation = mgl32.Translate3D(p.position.X(), p.position.Y(), p.position.Z())
	p.modelMatrix.dirty = true
}

func (p *Primitive2D) SetAnchor(anchor mgl32.Vec2) {
	p.anchor = anchor
	p.modelMatrix.anchor = mgl32.Translate3D(-p.anchor.X(), -p.anchor.Y(), 0)
	p.modelMatrix.dirty = true
}

func (p *Primitive2D) SetAngle(radians float32) {
	p.angle = radians
	p.modelMatrix.rotation = mgl32.HomogRotate3DZ(p.angle)
	p.modelMatrix.dirty = true
}

func (p *Primitive2D) SetSize(size mgl32.Vec2) {
	p.size = size
	p.modelMatrix.size = mgl32.Scale3D(p.size.X(), p.size.Y(), 1)
	p.modelMatrix.dirty = true
}

func (p *Primitive2D) SetScale(scale mgl32.Vec2) {
	p.scale = scale
	p.rebuildScaleMatrix()
}

func (p *Primitive2D) SetFlipX(flipX bool) {
	p.flipX = flipX
	p.rebuildScaleMatrix()
}

func (p *Primitive2D) SetFlipY(flipY bool) {
	p.flipY = flipY
	p.rebuildScaleMatrix()
}

func (p *Primitive2D) SetColor(color Color) {
	p.color = color
}

func (p *Primitive2D) SetUniforms() {
	p.shaderProgram.SetUniform("color", &p.color)
	p.shaderProgram.SetUniform("mModel", p.ModelMatrix())
}

func (p *Primitive2D) SetSizeFromTexture() {
	p.SetSize(mgl32.Vec2{float32(p.texture.width), float32(p.texture.height)})
}

func (p *Primitive2D) SetAnchorToCenter() {
	p.SetAnchor(mgl32.Vec2{p.size[0] / 2.0, p.size[1] / 2.0})
}

func (p *Primitive2D) EnqueueForDrawing(context *Context) {
	context.EnqueueForDrawing(p)
}

func (p *Primitive2D) Draw(context *Context) {
	shaderId := p.shaderProgram.Id()
	gl.BindTexture(gl.TEXTURE_2D, p.texture.Id())
	gl.UseProgram(shaderId)
	p.shaderProgram.SetUniform("mProjection", &context.projectionMatrix)
	p.SetUniforms()
	gl.BindVertexArray(p.vaoId)
	gl.DrawArrays(p.arrayMode, 0, p.arraySize)
}

// Texture and shaders are already bound when this is called
func (p *Primitive2D) DrawInBatch(context *Context) {
	p.SetUniforms()
	gl.BindVertexArray(p.vaoId)
	gl.DrawArrays(p.arrayMode, 0, p.arraySize)
}

func (p *Primitive2D) rebuildMatrices() {
	p.modelMatrix.translation = mgl32.Translate3D(p.position.X(), p.position.Y(), p.position.Z())
	p.modelMatrix.anchor = mgl32.Translate3D(-p.anchor.X(), -p.anchor.Y(), 0)
	p.modelMatrix.rotation = mgl32.HomogRotate3DZ(p.angle)
	p.modelMatrix.size = mgl32.Scale3D(p.size.X(), p.size.Y(), 1)
	p.rebuildScaleMatrix()

	p.modelMatrix.dirty = true
}

func (p *Primitive2D) rebuildScaleMatrix() {
	scaleX := p.scale.X()
	if p.flipX {
		scaleX *= -1
	}
	scaleY := p.scale.Y()
	if p.flipY {
		scaleY *= -1
	}
	p.modelMatrix.scale = mgl32.Scale3D(scaleX, scaleY, 1)
	p.modelMatrix.dirty = true
}

func (p *Primitive2D) ModelMatrix() *mgl32.Mat4 {
	if p.modelMatrix.dirty {
		p.modelMatrix.Mat4 = p.modelMatrix.translation.Mul4(p.modelMatrix.rotation).Mul4(p.modelMatrix.scale).Mul4(p.modelMatrix.anchor).Mul4(p.modelMatrix.size)
		//p.modelMatrix.Mat4 = p.modelMatrix.translation.Mul4(p.modelMatrix.size)
	}
	return &p.modelMatrix.Mat4
}

func NewQuadPrimitive(position mgl32.Vec3, size mgl32.Vec2) *Primitive2D {
	q := &Primitive2D{}
	q.position = position
	q.size = size
	q.scale = mgl32.Vec2{1, 1}
	q.shaderProgram = NewShaderProgram(VertexShaderPrimitive2D, "", FragmentShaderTexture)
	q.rebuildMatrices()

	q.arrayMode = gl.TRIANGLE_FAN
	q.arraySize = 4

	// Build the VAO
	q.SetVertices([]float32{0, 0, 0, 1, 1, 1, 1, 0})
	q.SetUVCoords([]float32{0, 0, 0, 1, 1, 1, 1, 0})
	return q
}

func NewRegularPolygonPrimitive(position mgl32.Vec3, radius float32, numSegments int, filled bool) *Primitive2D {
	circlePoints, err := utils.CircleToPolygon(mgl32.Vec2{0.5, 0.5}, 0.5, numSegments, 0)
	if err != nil {
		log.Panic(err)
		return nil
	}

	q := &Primitive2D{}
	q.position = position
	q.size = mgl32.Vec2{radius * 2, radius * 2}
	q.scale = mgl32.Vec2{1, 1}
	q.shaderProgram = NewShaderProgram(VertexShaderPrimitive2D, "", FragmentShaderSolidColor)
	q.rebuildMatrices()

	// Vertices
	vertices := make([]float32, 0, numSegments*2)
	for _, v := range circlePoints {
		vertices = append(vertices, v[0], v[1])
	}
	// Add one vertex for the last line
	vertices = append(vertices, circlePoints[0][0], circlePoints[0][1])

	if filled {
		q.arrayMode = gl.TRIANGLE_FAN
	} else {
		q.arrayMode = gl.LINE_STRIP
	}

	q.SetVertices(vertices)
	return q
}

func NewTriangles(
	vertices []float32,
	uvCoords []float32,
	texture *Texture,
	position mgl32.Vec3,
	size mgl32.Vec2,
	shaderProgram *ShaderProgram,
) *Primitive2D {
	p := &Primitive2D{}
	p.arrayMode = gl.TRIANGLES
	p.arraySize = int32(len(vertices) / 2)
	p.texture = texture
	p.shaderProgram = shaderProgram
	p.position = position
	p.scale = mgl32.Vec2{1, 1}
	p.size = size
	p.rebuildMatrices()
	gl.GenVertexArrays(1, &p.vaoId)
	gl.BindVertexArray(p.vaoId)
	p.SetVertices(vertices)
	p.SetUVCoords(uvCoords)
	gl.BindVertexArray(0)
	return p
}

func NewPolylinePrimitive(position mgl32.Vec3, points []mgl32.Vec2, closed bool) *Primitive2D {
	topLeft, bottomRight := utils.GetBoundingBox(points)

	primitive := &Primitive2D{}
	primitive.position = position
	primitive.size = bottomRight.Sub(topLeft)
	primitive.scale = mgl32.Vec2{1, 1}
	primitive.shaderProgram = NewShaderProgram(VertexShaderPrimitive2D, "", FragmentShaderSolidColor)
	primitive.rebuildMatrices()

	// Vertices
	vertices := make([]float32, 0, len(points)*2)
	for _, p := range points {
		// The vertices coordinates are relative to the top left and are scaled by size
		vertices = append(vertices, (p[0]-topLeft[0])/primitive.size.X(), (p[1]-topLeft[1])/primitive.size.Y())
	}
	if closed {
		// Add the first point again to close the loop
		vertices = append(vertices, vertices[0], vertices[1])
	}

	primitive.arrayMode = gl.LINE_STRIP
	primitive.arraySize = int32(len(vertices) / 2)
	primitive.SetVertices(vertices)
	return primitive
}

// SetVertices uploads new set of vertices into opengl buffer
func (p *Primitive2D) SetVertices(vertices []float32) {
	if p.vaoId == 0 {
		gl.GenVertexArrays(1, &p.vaoId)
	}
	gl.BindVertexArray(p.vaoId)
	if p.vboVertices == 0 {
		gl.GenBuffers(1, &p.vboVertices)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, p.vboVertices)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	p.arraySize = int32(len(vertices) / 2)
	gl.BindVertexArray(0)
}

// SetUVCoords uploads new UV coordinates
func (p *Primitive2D) SetUVCoords(uvCoords []float32) {
	if p.vaoId == 0 {
		gl.GenVertexArrays(1, &p.vaoId)
	}
	gl.BindVertexArray(p.vaoId)
	if p.vboUVCoords == 0 {
		gl.GenBuffers(1, &p.vboUVCoords)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, p.vboUVCoords)
	gl.BufferData(gl.ARRAY_BUFFER, len(uvCoords)*FLOAT32_SIZE, gl.Ptr(uvCoords), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}

const (
	VertexShaderPrimitive2D = `
        #version 410 core

        uniform mat4 mModel;
        uniform mat4 mProjection;

        layout(location=0) in vec2 vertex;
        layout(location=1) in vec2 uv;

        out vec2 uv_out;

        void main() {
            vec4 vertex_world = mModel * vec4(vertex, 0, 1);
            gl_Position = mProjection * vertex_world;
            uv_out = uv;
        }
        ` + "\x00"
)
