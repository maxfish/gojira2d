package graphics

import (
	"log"

	"github.com/go-gl/gl/all-core/gl"
)

type FrameBuffer struct {
	Texture // Target texture
	fbo     uint32
}

func NewFrameBuffer(width int, height int) *FrameBuffer {
	fb := &FrameBuffer{}
	gl.GenFramebuffers(1, &fb.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, fb.fbo)

	// Create and bind an empty texture
	fb.Texture = *NewEmptyTexture(width, height)
	textureId := fb.Texture.id
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, textureId, 0)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		log.Print("FrameBuffer creation error")
		return nil
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	return fb
}

func (f *FrameBuffer) Width() int32 {
	return f.Texture.width
}

func (f *FrameBuffer) Height() int32 {
	return f.Texture.height
}

func (f *FrameBuffer) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.fbo)
}

func (f *FrameBuffer) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}
