package graphics

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Quad2DPrimitive struct {
	Primitive2D
}

func NewQuad2DPrimitive(position mgl32.Vec3, size mgl32.Vec2) (*Quad2DPrimitive) {
	q := Quad2DPrimitive{}
	q.position = position
	q.size = size
	q.scale = mgl32.Vec2{1, 1}
	q.generateBuffers()
	q.shaderProgram = NewShaderProgram(vertexShaderQuad, "", FragmentShaderTexture)
	q.invalidateMatrices()

	return &q
}

func (q *Quad2DPrimitive) Draw(context *Context) {
	shaderId := q.shaderProgram.Id()
	gl.BindTexture(gl.TEXTURE_2D, q.texture.Id())
	gl.UseProgram(shaderId)
	q.shaderProgram.SetUniform4fv("mProjection", &context.projectionMatrix)
	q.SetMatrices()
	gl.BindVertexArray(q.vaoId)
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
}

func (q *Quad2DPrimitive) EnqueueForDrawing(context *Context) {
	context.enqueueForDrawing(q)
}

func (q *Quad2DPrimitive) drawInBatch(context *Context) {
	q.SetMatrices()
	gl.BindVertexArray(q.vaoId)
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
}

func (q *Quad2DPrimitive) SetSizeFromTexture() {
	q.SetSize(mgl32.Vec2{float32(q.texture.width), float32(q.texture.height)})
}

func (q *Quad2DPrimitive) SetAnchorToCenter() {
	q.SetAnchor(mgl32.Vec2{q.size[0] / 2.0, q.size[1] / 2.0})
}

func (q *Quad2DPrimitive) generateBuffers() {
	gl.GenVertexArrays(1, &q.vaoId)
	gl.BindVertexArray(q.vaoId)

	// Vertices
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	vertices := []float32{0, 0, 0, 1, 1, 1, 1, 0}
	// 4 -> magic constant for float32 size
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	// Texture coordinates
	var vboUV uint32
	gl.GenBuffers(1, &vboUV)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboUV)
	//self._texture_coordinates = np.array([0, 0, 0, 1, 1, 1, 1, 0], dtype=np.int16)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	gl.BindVertexArray(0)
}

func (q *Quad2DPrimitive) invalidateMatrices() {
	q.matrixTranslation.Dirty = true
	q.matrixSize.Dirty = true
	q.matrixRotation.Dirty = true
	q.matrixScale.Dirty = true
	q.matrixAnchor.Dirty = true
}

func (q *Quad2DPrimitive) SetMatrices() {
	if q.matrixTranslation.Dirty {
		q.matrixTranslation.Matrix = mgl32.Translate3D(q.position.X(), q.position.Y(), q.position.Z())
		q.matrixTranslation.Dirty = false
		q.shaderProgram.SetUniform4fv("mTranslate", &q.matrixTranslation.Matrix)
	}
	if q.matrixScale.Dirty {
		scaleX := q.scale.X()
		if q.flipX {
			scaleX *= -1
		}
		scaleY := q.scale.Y()
		if q.flipY {
			scaleY *= -1
		}
		q.matrixScale.Matrix = mgl32.Scale3D(scaleX, scaleY, 1)
		q.shaderProgram.SetUniform4fv("mScale", &q.matrixScale.Matrix)
	}
	if q.matrixSize.Dirty {
		q.matrixSize.Matrix = mgl32.Scale3D(q.size.X(), q.size.Y(), 1)
		q.shaderProgram.SetUniform4fv("mSize", &q.matrixSize.Matrix)
	}
	if q.matrixRotation.Dirty {
		q.matrixRotation.Matrix = mgl32.HomogRotate3DZ(q.angle)
		q.shaderProgram.SetUniform4fv("mRotation", &q.matrixRotation.Matrix)
	}
	if q.matrixRotation.Dirty {
		q.matrixAnchor.Matrix = mgl32.Translate3D(-q.anchor.X(), -q.anchor.Y(), 0)
		q.shaderProgram.SetUniform4fv("mAnchor", &q.matrixAnchor.Matrix)
	}
}

const (
	vertexShaderQuad = `
        #version 410 core

        uniform mat4 mProjection;
        uniform mat4 mScale;
        uniform mat4 mTranslate;
        uniform mat4 mSize;
        uniform mat4 mRotation;
        uniform mat4 mAnchor;

        layout(location=0) in vec2 vertex;
        layout(location=1) in vec2 uv;

        out vec2 uv_out;

        void main() {
			mat4 mModel = mTranslate * mRotation * mScale * mAnchor * mSize;
            vec4 vertex_world = mModel * vec4(vertex, 0, 1);
            gl_Position = mProjection * vertex_world;
            uv_out = uv;
        }
        ` + "\x00"
)
