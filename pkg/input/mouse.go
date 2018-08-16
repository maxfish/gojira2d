package input

import (
	"fmt"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/maxfish/gojira2d/pkg/app"
)

var (
	positionCallback glfw.CursorPosCallback
	buttonCallback   glfw.MouseButtonCallback

	connected   bool
	posX        int
	posY        int
	buttonsDown []bool
)

func registerMouseCallbacks(pCallback glfw.CursorPosCallback, bCallback glfw.MouseButtonCallback) {
	if positionCallback != nil {
		fmt.Printf("Error: Mouse callbacks have been already registered!")
	}
	positionCallback = pCallback
	buttonCallback = bCallback
	app.GetWindow().SetCursorPosCallback(positionCallback)
	app.GetWindow().SetMouseButtonCallback(buttonCallback)
}

func unregisterMouseCallbacks() {
	positionCallback = nil
	buttonCallback = nil
	app.GetWindow().SetPosCallback(nil)
	app.GetWindow().SetMouseButtonCallback(nil)
}

// Registers the callbacks and starts receiving the mouse's events
func ConnectMouse() {
	pCallback := func(w *glfw.Window, x float64, y float64) {
		posX = int(x)
		posY = int(y)
	}
	bCallback := func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if action == glfw.Press || action == glfw.Repeat {
			buttonsDown[button] = true
		} else {
			buttonsDown[button] = false
		}
	}

	buttonsDown = make([]bool, glfw.MouseButtonLast, glfw.MouseButtonLast)
	registerMouseCallbacks(pCallback, bCallback)
	connected = true
}

// Unregisters the callbacks and stops receiving the mouse's events
func DisconnectMouse() {
	connected = false
	unregisterMouseCallbacks()
	buttonsDown = nil
}

// Returns the coordinates of the cursor's position
func MousePosition() (int, int) {
	if !connected {
		return 0, 0
	}
	return posX, posY
}

// Returns the state of the specified mouse button
func MouseButton(index int) bool {
	if !connected {
		return false
	}
	return buttonsDown[index]
}
