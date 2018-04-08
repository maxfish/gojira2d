package main

import (
	g "gojira2d/pkg/graphics"
	"github.com/go-gl/mathgl/mgl32"
	a "gojira2d/pkg/app"
	"math"
)

func main() {
	app := a.InitApp(640, 480, false, "Shapes")
	defer a.TerminateApp()

	primitives := []*g.Primitive2D{
		g.NewRegularPolygonPrimitive(mgl32.Vec3{100, 100, 0}, 50, 5, false),
		g.NewRegularPolygonPrimitive(mgl32.Vec3{250, 100, 0}, 50, 6, true),
		g.NewRegularPolygonPrimitive(mgl32.Vec3{400, 100, 0}, 50, 8, false),
		g.NewRegularPolygonPrimitive(mgl32.Vec3{550, 100, 0}, 50, 12, true),
		g.NewPolylinePrimitive(mgl32.Vec3{50, 420, 0}, []mgl32.Vec2{{60, 20}, {20, 20}, {20, 70}, {60, 70}, {60, 45}, {40, 45}}, false),
		g.NewPolylinePrimitive(mgl32.Vec3{110, 420, 0}, []mgl32.Vec2{{0, 0}, {0, 50}, {40, 50}}, false),
	}

	for _, p := range primitives {
		p.SetAnchorToCenter()
		p.SetColor(g.Color{1,1,1,1})
	}

	primitives[3].SetColor(g.Color{1,0,0,1})

	var animationAngle float32 = 0

	app.MainLoop(func(speed float64) {
		animationAngle += float32(speed)
		animationScale := float32(math.Abs(math.Sin(float64(animationAngle))))
		primitives[1].SetScale(mgl32.Vec2{animationScale, animationScale})
		primitives[2].SetAngle(animationAngle)
		primitives[3].SetAngle(animationAngle)
		primitives[3].SetScale(mgl32.Vec2{animationScale, animationScale})
	}, func() {
		for _, p := range primitives {
			p.EnqueueForDrawing(app.Context)
		}
	})
}
