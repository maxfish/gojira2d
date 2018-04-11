package graphics

import "github.com/go-gl/mathgl/mgl32"

type PostProcessingStep struct {
	quad        *Primitive2D
	frameBuffer *FrameBuffer
}

func NewPostProcessingStep(width int, height int) *PostProcessingStep {
	step := &PostProcessingStep{}
	step.frameBuffer = NewFrameBuffer(width, height)
	step.quad = NewQuadPrimitive(mgl32.Vec3{0, 0, 0}, mgl32.Vec2{float32(width), float32(height)})
	return step
}

func (p *PostProcessingStep) FrameBuffer() *FrameBuffer {
	return p.frameBuffer
}

func (p *PostProcessingStep) Quad() *Primitive2D {
	return p.quad
}
