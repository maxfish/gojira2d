package graphics

import "github.com/go-gl/mathgl/mgl32"

type Matrix struct {
	Matrix mgl32.Mat4
	Dirty  bool
}

type Drawable interface {
	Texture() (*Texture)
	Shader() (*ShaderProgram)
	Draw(context *Context)
	drawInBatch(context *Context)
}

type Primitive struct {
	vaoId         uint32
	texture       *Texture
	shaderProgram *ShaderProgram
}

func (p *Primitive) SetTexture(texture *Texture) {
	p.texture = texture
}

func (p *Primitive) Texture() (*Texture) {
	return p.texture
}

func (p *Primitive) SetShader(shader *ShaderProgram) {
	p.shaderProgram = shader
}

func (p *Primitive) Shader() (*ShaderProgram) {
	return p.shaderProgram
}

func (p *Primitive) Draw(context *Context) {
}

func (p *Primitive) drawInBatch(context *Context) {
}

type Primitive2D struct {
	Primitive
	position          mgl32.Vec3
	scale             mgl32.Vec2
	size              mgl32.Vec2
	anchor            mgl32.Vec2
	angle             float32
	flipX             bool
	flipY             bool
	matrixSize        Matrix
	matrixTranslation Matrix
	matrixRotation    Matrix
	matrixScale       Matrix
	matrixAnchor      Matrix
}

func (p *Primitive2D) SetPosition(position mgl32.Vec3) {
	p.position = position
	p.matrixTranslation.Dirty = true
}

func (p *Primitive2D) SetAnchor(anchor mgl32.Vec2) {
	p.anchor = anchor
	p.matrixAnchor.Dirty = true
}

func (p *Primitive2D) SetAngle(radians float32) {
	p.angle = radians
	p.matrixRotation.Dirty = true
}

func (p *Primitive2D) SetScale(scale mgl32.Vec2) {
	p.scale = scale
	p.matrixScale.Dirty = true
}

func (p *Primitive2D) SetSize(size mgl32.Vec2) {
	p.size = size
	p.matrixSize.Dirty = true
}

func (p *Primitive2D) SetFlipX(flipX bool) {
	p.flipX = flipX
	p.matrixScale.Dirty = true
}

func (p *Primitive2D) SetFlipY(flipY bool) {
	p.flipY = flipY
	p.matrixScale.Dirty = true
}

func (p *Primitive2D) rebuildMatrices() {
	p.matrixTranslation.Matrix = mgl32.Translate3D(p.position.X(), p.position.Y(), p.position.Z())
	p.matrixScale.Matrix = mgl32.Scale3D(p.scale.X(), p.scale.Y(), 1)
	p.matrixRotation.Matrix = mgl32.HomogRotate3DZ(p.angle)
	p.matrixSize.Matrix = mgl32.Scale3D(p.size.X(), p.size.Y(), 1)
	p.matrixAnchor.Matrix = mgl32.Translate3D(-p.anchor.X(), -p.anchor.Y(), 0)
}
