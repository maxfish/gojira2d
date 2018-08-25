package graphics

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/maxfish/gojira2d/pkg/utils"
)

const (
	// Float32Size is the size (in bytes) of a float32
	Float32Size = 4
)

// ModelMatrix matrix representing the primitive transformation
type ModelMatrix struct {
	mgl64.Mat4
	mat32       mgl32.Mat4
	size        mgl64.Mat4
	translation mgl64.Mat4
	rotation    mgl64.Mat4
	scale       mgl64.Mat4
	anchor      mgl64.Mat4
	dirty       bool
}

// Primitive2D a drawing primitive on the XY plane
type Primitive2D struct {
	Primitive
	position    mgl64.Vec3
	scale       mgl64.Vec2
	size        mgl64.Vec2
	anchor      mgl64.Vec2
	angle       float64
	flipX       bool
	flipY       bool
	color       Color
	modelMatrix ModelMatrix
}

// SetPosition sets the X,Y,Z position of the primitive. Z is used for the drawing order
func (p *Primitive2D) SetPosition(position mgl64.Vec3) {
	p.position = position
	p.modelMatrix.translation = mgl64.Translate3D(p.position.X(), p.position.Y(), p.position.Z())
	p.modelMatrix.dirty = true
}

// Position gets X,Y,Z of the primitive.
func (p *Primitive2D) Position() mgl64.Vec3 {
	return p.position
}

// SetAnchor sets the anchor point of the primitive, this will be the point placed at Position
func (p *Primitive2D) SetAnchor(anchor mgl64.Vec2) {
	p.anchor = anchor
	p.modelMatrix.anchor = mgl64.Translate3D(-p.anchor.X(), -p.anchor.Y(), 0)
	p.modelMatrix.dirty = true
}

// SetAnchorToCenter sets the anchor at the center of the primitive
func (p *Primitive2D) SetAnchorToCenter() {
	p.SetAnchor(mgl64.Vec2{p.size[0] / 2.0, p.size[1] / 2.0})
}

// Angle in radians
func (p *Primitive2D) Angle() float64 {
	return p.angle
}

// SetAngle sets the rotation angle around the Z axis
func (p *Primitive2D) SetAngle(radians float64) {
	p.angle = radians
	p.modelMatrix.rotation = mgl64.HomogRotate3DZ(p.angle)
	p.modelMatrix.dirty = true
}

// Size in pixels
func (p *Primitive2D) Size() mgl64.Vec2 {
	return mgl64.Vec2{p.size.X(), p.size.Y()}
}

// SetSize sets the size (in pixels) of the current primitive
func (p *Primitive2D) SetSize(size mgl64.Vec2) {
	p.size = size
	p.modelMatrix.size = mgl64.Scale3D(p.size.X(), p.size.Y(), 1)
	p.modelMatrix.dirty = true
}

// SetSizeFromTexture sets the size of the current primitive to the pixel size of the texture
func (p *Primitive2D) SetSizeFromTexture() {
	if p.texture == nil {
		return
	}
	p.SetSize(mgl64.Vec2{float64(p.texture.width), float64(p.texture.height)})
}

// SetScale sets the scaling factor on X and Y for the primitive. The scaling respects the anchor and the rotation
func (p *Primitive2D) SetScale(scale mgl64.Vec2) {
	p.scale = scale
	p.rebuildScaleMatrix()
}

// SetFlipX flips the primitive around the Y axis
func (p *Primitive2D) SetFlipX(flipX bool) {
	p.flipX = flipX
	p.rebuildScaleMatrix()
}

// SetFlipY flips the primitive around the X axis
func (p *Primitive2D) SetFlipY(flipY bool) {
	p.flipY = flipY
	p.rebuildScaleMatrix()
}

// SetColor sets the color passed to the shader
func (p *Primitive2D) SetColor(color Color) {
	p.color = color
}

// SetUniforms sets the shader's uniform variables
func (p *Primitive2D) SetUniforms() {
	p.shaderProgram.SetUniform("color", &p.color)
	p.shaderProgram.SetUniform("model", p.ModelMatrix32())
}

// Draw draws the primitive
func (p *Primitive2D) Draw(context *Context) {
	shaderID := p.shaderProgram.ID()
	context.BindTexture(p.texture)
	gl.UseProgram(shaderID)
	cameraMatrix := context.Camera2D.ProjectionMatrix32()
	p.shaderProgram.SetUniform("projection", &cameraMatrix)
	p.SetUniforms()
	gl.BindVertexArray(p.vaoId)
	gl.DrawArrays(p.arrayMode, 0, p.arraySize)
}

// DrawInBatch draws the primitive assuming that the correct texture and shader are already bound
func (p *Primitive2D) DrawInBatch(context *Context) {
	p.SetUniforms()
	gl.BindVertexArray(p.vaoId)
	gl.DrawArrays(p.arrayMode, 0, p.arraySize)
}

func (p *Primitive2D) rebuildMatrices() {
	p.modelMatrix.translation = mgl64.Translate3D(p.position.X(), p.position.Y(), p.position.Z())
	p.modelMatrix.anchor = mgl64.Translate3D(-p.anchor.X(), -p.anchor.Y(), 0)
	p.modelMatrix.rotation = mgl64.HomogRotate3DZ(p.angle)
	p.modelMatrix.size = mgl64.Scale3D(p.size.X(), p.size.Y(), 1)
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
	p.modelMatrix.scale = mgl64.Scale3D(scaleX, scaleY, 1)
	p.modelMatrix.dirty = true
}

func (p *Primitive2D) rebuildModelMatrix() {
	if p.modelMatrix.dirty {
		p.modelMatrix.Mat4 = p.modelMatrix.translation.Mul4(p.modelMatrix.rotation).Mul4(p.modelMatrix.scale).Mul4(p.modelMatrix.anchor).Mul4(p.modelMatrix.size)
		// updates the float32 version
		p.modelMatrix.mat32 = utils.Mat4From64to32Bits(p.modelMatrix.Mat4)
		p.modelMatrix.dirty = false
	}
}

// ModelMatrix returns the current model matrix
func (p *Primitive2D) ModelMatrix() *mgl64.Mat4 {
	p.rebuildModelMatrix()
	return &p.modelMatrix.Mat4
}

// ModelMatrix returns the current model matrix as mgl32.Mat4
func (p *Primitive2D) ModelMatrix32() *mgl32.Mat4 {
	p.rebuildModelMatrix()
	return &p.modelMatrix.mat32
}

// NewQuadPrimitive creates a rectangular primitive
func NewQuadPrimitive(position mgl64.Vec3, size mgl64.Vec2) *Primitive2D {
	q := &Primitive2D{
		position: position,
		size:     size,
		scale:    mgl64.Vec2{1, 1},
	}
	q.shaderProgram = NewShaderProgram(VertexShaderBase, "", FragmentShaderTexture)
	q.rebuildMatrices()
	q.arrayMode = gl.TRIANGLE_FAN
	q.arraySize = 4

	// Build the VAO
	q.SetVertices([]float32{0, 0, 0, 1, 1, 1, 1, 0})
	q.SetUVCoords([]float32{0, 0, 0, 1, 1, 1, 1, 0})
	return q
}

// NewRegularPolygonPrimitive creates a primitive from a regular polygon
func NewRegularPolygonPrimitive(center mgl64.Vec3, radius float64, numSegments int, filled bool) *Primitive2D {
	circlePoints, err := utils.CircleToPolygon(mgl64.Vec2{0, 0}, radius, numSegments, 0)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	q := &Primitive2D{
		position: center,
		size:     mgl64.Vec2{1, 1},
		scale:    mgl64.Vec2{1, 1},
	}
	q.shaderProgram = NewShaderProgram(VertexShaderBase, "", FragmentShaderSolidColor)
	q.rebuildMatrices()

	// Vertices
	vertices := make([]float32, 0, numSegments*2)
	for _, v := range circlePoints {
		vertices = append(vertices, float32(v[0]), float32(v[1]))
	}
	// Add one vertex for the last line
	vertices = append(vertices, float32(circlePoints[0][0]), float32(circlePoints[0][1]))

	if filled {
		q.arrayMode = gl.TRIANGLE_FAN
	} else {
		q.arrayMode = gl.LINE_STRIP
	}

	q.SetVertices(vertices)
	return q
}

// NewTriangles creates a primitive as a collection of triangles
func NewTriangles(
	vertices []float32,
	uvCoords []float32,
	texture *Texture,
	position mgl64.Vec3,
	size mgl64.Vec2,
	shaderProgram *ShaderProgram,
) *Primitive2D {
	p := &Primitive2D{
		position: position,
		size:     size,
		scale:    mgl64.Vec2{1, 1},
	}
	p.arrayMode = gl.TRIANGLES
	p.arraySize = int32(len(vertices) / 2)
	p.texture = texture
	p.shaderProgram = shaderProgram
	p.rebuildMatrices()
	gl.GenVertexArrays(1, &p.vaoId)
	gl.BindVertexArray(p.vaoId)
	p.SetVertices(vertices)
	p.SetUVCoords(uvCoords)
	gl.BindVertexArray(0)
	return p
}

// NewPolylinePrimitive creates a primitive from a sequence of points. The points coordinates are relative to the passed center
func NewPolylinePrimitive(center mgl64.Vec3, points []mgl64.Vec2, closed bool) *Primitive2D {
	primitive := &Primitive2D{
		position: center,
		size:     mgl64.Vec2{1, 1},
		scale:    mgl64.Vec2{1, 1},
	}
	primitive.shaderProgram = NewShaderProgram(VertexShaderBase, "", FragmentShaderSolidColor)
	primitive.rebuildMatrices()

	// Vertices
	var numVertices int32 = int32(len(points))
	vertices := make([]float32, 0, numVertices*2)
	for _, p := range points {
		vertices = append(vertices, float32(p[0]), float32(p[1]))
	}
	if closed {
		// Add the first point again to close the loop
		vertices = append(vertices, vertices[0], vertices[1])
		numVertices++
	}

	primitive.arrayMode = gl.LINE_STRIP
	primitive.arraySize = numVertices
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
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*Float32Size, gl.Ptr(vertices), gl.STATIC_DRAW)
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
	gl.BufferData(gl.ARRAY_BUFFER, len(uvCoords)*Float32Size, gl.Ptr(uvCoords), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}
