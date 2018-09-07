package main

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/maxfish/gojira2d/pkg/app"
	g "github.com/maxfish/gojira2d/pkg/graphics"
)

func main() {
	app.Init(640, 480, "Quad")
	defer app.Terminate()

	app.SetFPSCounterVisible(true)
	app.Context.Camera2D.SetCentered(true)

	Quad := g.NewQuadPrimitive(mgl64.Vec3{0, 0, 0}, mgl64.Vec2{200, 200})
	Quad.SetAnchorToCenter()
	Quad.SetTexture(g.NewTextureFromFile("examples/assets/texture.png"))

	app.MainLoop(func(deltaTime float64) {
		// NOP
	}, func() {
		Quad.Draw(app.Context)
	})
}
