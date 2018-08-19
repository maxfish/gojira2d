package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/maxfish/gojira2d/pkg/app"
	g "github.com/maxfish/gojira2d/pkg/graphics"
)

func main() {
	app.Init(640, 480, "Shapes")
	defer app.Terminate()

	app.SetFPSCounterVisible(true)
	primitives := []*g.Primitive2D{
		g.NewRegularPolygonPrimitive(mgl64.Vec3{100, 100, 0}, 50, 5, false),
		g.NewRegularPolygonPrimitive(mgl64.Vec3{250, 100, 0}, 50, 6, true),
		g.NewRegularPolygonPrimitive(mgl64.Vec3{400, 100, 0}, 50, 8, false),
		g.NewRegularPolygonPrimitive(mgl64.Vec3{550, 100, 0}, 50, 12, true),
		g.NewPolylinePrimitive(mgl64.Vec3{50, 420, 0}, []mgl64.Vec2{{20, -20}, {-20, -20}, {-20, 40}, {20, 40}, {20, 10}, {0, 10}}, false),
		g.NewPolylinePrimitive(mgl64.Vec3{110, 420, 0}, []mgl64.Vec2{{-20, -20}, {-20, 40}, {20, 40}}, false),
	}

	for _, p := range primitives {
		p.SetColor(g.Color{1, 1, 1, 1})
	}

	primitives[3].SetColor(g.Color{1, 0, 0, 1})

	var animationAngle float64 = 0

	app.MainLoop(func(speed float64) {
		animationAngle += speed
		animationScale := math.Abs(math.Sin(float64(animationAngle)))
		primitives[1].SetScale(mgl64.Vec2{animationScale, animationScale})
		primitives[2].SetAngle(animationAngle)
		primitives[3].SetAngle(animationAngle)
		primitives[3].SetScale(mgl64.Vec2{animationScale, animationScale})
	}, func() {
		for _, p := range primitives {
			p.Draw(app.Context)
		}
	})
}
