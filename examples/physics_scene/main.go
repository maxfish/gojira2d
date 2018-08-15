package main

import (
	"github.com/maxfish/gojira2d/pkg/app"
	"github.com/maxfish/gojira2d/pkg/physics"
)

const pixelsPerMeter float64 = 50

func main() {
	app.Init(800, 600, "Physics")
	defer app.Terminate()

	app.SetFPSCounterVisible(true)

	// Use a cartesian coordinates system
	app.Context.Camera2D.SetFlipVertical(true)
	app.Context.Camera2D.SetPosition(50, 0)

	// Load the scene from a Rube file
	scene := physics.NewB2DJsonSceneFromFile("examples/assets/test_scene.json")
	b2World := scene.World

	debugDraw := physics.NewBox2DDebugDraw(b2World, pixelsPerMeter)

	app.MainLoop(func(speed float64) {
		b2World.Step(1/scene.StepsPerSecond, scene.VelocityIterations, scene.PositionIterations)
		debugDraw.Update()
	}, func() {
		debugDraw.Draw(app.Context)
	})
}
