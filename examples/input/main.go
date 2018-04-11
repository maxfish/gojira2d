package main

import (
	a "gojira2d/pkg/app"
	"gojira2d/pkg/input"
)

func main() {
	app := a.InitApp(640, 480, true, "Input")
	defer a.TerminateApp()

	joy := &input.JoystickController{}
	joy.Open(0)

	app.MainLoop(func(speed float64) {
		joy.Update()
	}, func() {
		// TODO: print some text with info on the joystick's status
	})
}
