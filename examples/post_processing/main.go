package main

import (
	a "gojira2d/pkg/app"
	g "gojira2d/pkg/graphics"

	"github.com/go-gl/mathgl/mgl32"
)

func main() {
	app := a.InitApp(640, 480, true, "Post processing")
	defer a.TerminateApp()

	Quad := g.NewQuadPrimitive(mgl32.Vec3{0, 0, 0}, mgl32.Vec2{200, 200})
	Quad.SetAnchorToCenter()
	Quad.SetTexture(g.NewTextureFromFile("examples/assets/texture.png"))

	step1 := g.NewPostProcessingStep(640, 480)
	shader1 := g.NewShaderProgramFromFiles("", "", "examples/post_processing/grayscale.frag")
	step1.Quad().SetShader(shader1)
	app.Context.AddPostProcessingStep(step1)
	step2 := g.NewPostProcessingStep(640, 480)
	shader2 := g.NewShaderProgramFromFiles("", "", "examples/post_processing/scanline.frag")
	step1.Quad().SetShader(shader2)
	app.Context.AddPostProcessingStep(step2)

	app.MainLoop(func(speed float64) {
		// NOP
	}, func() {
		Quad.EnqueueForDrawing(app.Context)
	})
}
