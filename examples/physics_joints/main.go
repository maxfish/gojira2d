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

	// Load the scene from a Rube file
	scene := physics.NewB2DJsonSceneFromFile("examples/assets/physics/scene_joints.json")
	b2World := scene.World

	// Set up the camera
	bb := scene.SceneBoundingBoxInPixels(pixelsPerMeter)
	app.Context.Camera2D.SetFlipVertical(true)
	app.Context.Camera2D.SetVisibleArea(float32(bb.LowerBound.X), float32(bb.LowerBound.Y), float32(bb.UpperBound.X), float32(bb.UpperBound.Y))

	debugDraw := physics.NewBox2DDebugDraw(b2World, pixelsPerMeter)

	app.MainLoop(func(speed float64) {
		b2World.Step(1/scene.StepsPerSecond, scene.VelocityIterations, scene.PositionIterations)
		debugDraw.Update()
	}, func() {
		debugDraw.Draw(app.Context)
	})
}
