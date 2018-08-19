package input

import (
	"fmt"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/maxfish/gojira2d/pkg/app"
)

var (
	positionCallback glfw.CursorPosCallback
	buttonCallback   glfw.MouseButtonCallback
	scrollCallback   glfw.ScrollCallback

	connected     bool
	posX          float64
	posY          float64
	deltaX        float64
	deltaY        float64
	scrollOffsetX float64
	scrollOffsetY float64
	buttonsDown   []bool
)

func registerMouseCallbacks(pCallback glfw.CursorPosCallback, bCallback glfw.MouseButtonCallback, sCallback glfw.ScrollCallback) {
	if positionCallback != nil {
		fmt.Printf("Error: Mouse callbacks have been already registered!")
	}
	positionCallback = pCallback
	buttonCallback = bCallback
	scrollCallback = sCallback
	app.GetWindow().SetCursorPosCallback(positionCallback)
	app.GetWindow().SetMouseButtonCallback(buttonCallback)
	app.GetWindow().SetScrollCallback(scrollCallback)
}

func unregisterMouseCallbacks() {
	positionCallback = nil
	buttonCallback = nil
	app.GetWindow().SetPosCallback(nil)
	app.GetWindow().SetMouseButtonCallback(nil)
	app.GetWindow().SetScrollCallback(nil)
}

// ConnectMouse Registers the callbacks and starts receiving the mouse's events
func ConnectMouse() {
	pCallback := func(w *glfw.Window, x float64, y float64) {
		deltaX = x - posX
		deltaY = y - posY
		posX = x
		posY = y
	}
	bCallback := func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if action == glfw.Press || action == glfw.Repeat {
			buttonsDown[button] = true
		} else {
			buttonsDown[button] = false
		}
	}
	sCallback := func(w *glfw.Window, xoff float64, yoff float64) {
		scrollOffsetX = xoff
		scrollOffsetY = yoff
	}

	buttonsDown = make([]bool, glfw.MouseButtonLast, glfw.MouseButtonLast)
	registerMouseCallbacks(pCallback, bCallback, sCallback)
	connected = true
}

// DisconnectMouse Unregisters the callbacks and stops receiving the mouse's events
func DisconnectMouse() {
	connected = false
	unregisterMouseCallbacks()
	buttonsDown = nil
}

// MousePosition Returns the coordinates of the cursor's position
func MousePosition() (float64, float64) {
	if !connected {
		ConnectMouse()
	}
	return posX, posY
}

// MouseDelta Returns the latest movement of the cursor
func MouseDelta() (float64, float64) {
	if !connected {
		ConnectMouse()
	}

	dX := deltaX
	dY := deltaY

	// The event has been consumed
	deltaX = 0
	deltaY = 0
	return dX, dY
}

// MouseButton Returns the state of the specified mouse button
func MouseButton(index int) bool {
	if !connected {
		ConnectMouse()
	}
	return buttonsDown[index]
}

// MouseScroll Returns the last wheel's scroll offsets
func MouseScroll() (float64, float64) {
	if !connected {
		ConnectMouse()
	}
	x, y := scrollOffsetX, scrollOffsetY

	// The scroll event has been consumed
	scrollOffsetX = 0
	scrollOffsetY = 0

	return x, y
}
