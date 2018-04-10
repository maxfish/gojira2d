package graphics

type Primitive struct {
	vaoId         uint32
	vboVertices   uint32
	vboUVCoords   uint32
	arrayMode     uint32
	arraySize     int32
	texture       *Texture
	shaderProgram *ShaderProgram
}

func (p *Primitive) SetTexture(texture *Texture) {
	p.texture = texture
}

func (p *Primitive) Texture() *Texture {
	return p.texture
}

func (p *Primitive) SetShader(shader *ShaderProgram) {
	p.shaderProgram = shader
}

func (p *Primitive) Shader() *ShaderProgram {
	return p.shaderProgram
}

func (p *Primitive) Draw(context *Context) {
}

func (p *Primitive) drawInBatch(context *Context) {
}
