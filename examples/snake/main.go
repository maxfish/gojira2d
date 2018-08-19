package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/maxfish/gojira2d/pkg/app"
	g "github.com/maxfish/gojira2d/pkg/graphics"
	"github.com/maxfish/gojira2d/pkg/input"
)

func main() {
	app.Init(640, 480, "Snake")
	defer app.Terminate()
	app.SetFPSCounterVisible(true)

	kbd := &input.KeyboardController{}
	kbd.Open(-1)

	snake := []*g.Primitive2D{
		g.NewRegularPolygonPrimitive(mgl32.Vec3{100, 100, 0}, 9, 4, true),
		g.NewRegularPolygonPrimitive(mgl32.Vec3{100, 80, 0}, 9, 4, true),
		g.NewRegularPolygonPrimitive(mgl32.Vec3{100, 60, 0}, 9, 4, true),
	}
	for _, s := range snake {
		s.SetColor(g.Color{1, 1, 1, 1})
		s.SetAngle(45 * math.Pi / 180)
	}

	directionX := 1
	directionY := 0
	moveDelay := 0.2
	timePassed := 0.0

	app.MainLoop(func(speed float64) {
		directionX, directionY = UpdateDirection(kbd, directionX, directionY)

		timePassed += speed
		if timePassed > float64(moveDelay) {
			var oldPos mgl32.Vec3
			for idx, part := range snake {
				if idx == 0 {
					oldPos = part.Position()
					part.SetPosition(
						mgl32.Vec3{
							oldPos.X() + float32(directionX)*20,
							oldPos.Y() + float32(directionY)*20,
							0,
						},
					)
				} else {
					tmpPos := part.Position()
					part.SetPosition(oldPos)
					oldPos = tmpPos
				}
			}
			timePassed = 0
		}
	}, func() {
		for _, p := range snake {
			p.Draw(app.Context)
		}
	})
}

// UpdateDirection ...
func UpdateDirection(kbd *input.KeyboardController, dx, dy int) (int, int) {
	kbd.Update()
	if kbd.ButtonPressed(input.ButtonDirPadUp) {
		if dy != 1 {
			dx = 0
			dy = -1
		}
	}
	if kbd.ButtonPressed(input.ButtonDirPadDown) {
		if dy != -1 {
			dx = 0
			dy = 1
		}
	}
	if kbd.ButtonPressed(input.ButtonDirPadLeft) {
		if dx != 1 {
			dx = -1
			dy = 0
		}
	}
	if kbd.ButtonPressed(input.ButtonDirPadRight) {
		if dx != -1 {
			dx = 1
			dy = 0
		}
	}
	return dx, dy
}
