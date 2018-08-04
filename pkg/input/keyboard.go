package input

import (
	"log"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/maxfish/gojira2d/pkg/app"
)

var (
	keyCallbackFunc glfw.KeyCallback
)

func RegisterKeyCallback(callback glfw.KeyCallback) {
	if keyCallbackFunc != nil {
		log.Panic("A keyboard key-callback is already registered!")
	}
	keyCallbackFunc = callback
	app.GetWindow().SetKeyCallback(callback)
}

func UnregisterKeyCallback() {
	keyCallbackFunc = nil
	app.GetWindow().SetKeyCallback(nil)
}

func IsKeyboardFree() bool {
	return keyCallbackFunc == nil
}
