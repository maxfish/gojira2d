package main

import (
	a "github.com/maxfish/gojira2d/pkg/app"
	"github.com/maxfish/gojira2d/pkg/physics"
)

const pixelsPerMeter float64 = 50

func main() {
	a.Init(800, 600, "Physics")
	defer a.Terminate()

	// Use a cartesian coordinates system
	a.Context.Camera2D.SetFlipVertical(true)
	a.Context.Camera2D.SetPosition(50, 0)

	// Load the scene from a Rube file
	scene := physics.NewB2DJsonSceneFromFile("examples/assets/test_scene.json")
	b2World := scene.World

	debugDraw := physics.NewBox2DDebugDraw(b2World, pixelsPerMeter)

	a.MainLoop(func(speed float64) {
		b2World.Step(1/scene.StepsPerSecond, scene.VelocityIterations, scene.PositionIterations)
		debugDraw.Update()
	}, func() {
		debugDraw.Draw(a.Context)
	})
}
