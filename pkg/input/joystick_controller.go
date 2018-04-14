package input

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type JoystickController struct {
	GameController
	connected       bool
	joystick        glfw.Joystick
	name            string
	numButtons      int
	numAxes         int
	axes            []float32
	buttonsPressed  []bool
	buttonsReleased []bool
	buttonsDown     []bool
}

func (c *JoystickController) Open(deviceIndex int) bool {
	if c.connected {
		if glfw.Joystick(deviceIndex) == c.joystick {
			// Device already connected
			return true
		} else {
			// We are opening another device
			c.Close()
		}
	}

	if !glfw.JoystickPresent(glfw.Joystick(deviceIndex)) {
		c.Close()
		return false
	}

	c.connected = true
	c.joystick = glfw.Joystick(deviceIndex)
	c.name = glfw.GetJoystickName(c.joystick)

	// Get the num of buttons and axes
	c.numButtons = len(glfw.GetJoystickButtons(c.joystick))
	c.numAxes = len(glfw.GetJoystickAxes(c.joystick))
	// Build the slices
	c.buttonsDown = make([]bool, c.numButtons)
	c.buttonsPressed = make([]bool, c.numButtons)
	c.buttonsReleased = make([]bool, c.numButtons)

	return true
}

func (c *JoystickController) Close() {
	c.connected = false
	c.joystick = -1
	c.name = ""
	c.buttonsDown = nil
	c.buttonsPressed = nil
	c.buttonsReleased = nil
}

func (c *JoystickController) Update() {
	if !c.connected {
		return
	}

	// Buttons
	for i, button := range glfw.GetJoystickButtons(c.joystick) {
		isDown := button > 0
		c.buttonsPressed[i] = isDown && !c.buttonsDown[i]
		c.buttonsReleased[i] = !isDown && c.buttonsDown[i]
		c.buttonsDown[i] = isDown
	}

	// Axes
	c.axes = glfw.GetJoystickAxes(c.joystick)
}

func (c *JoystickController) AxisValue(axis ControllerAxis) float32 {
	return c.axes[axis]
}

func (c *JoystickController) AxisDigitalValue(axis ControllerAxis) int {
	// TODO: Define dead zone...
	return 0
}

func (c *JoystickController) Connected() bool {
	return c.connected
}

func (c *JoystickController) NumButtons() int {
	return c.numButtons
}

func (c *JoystickController) NumAxis() int {
	return c.numAxes
}

func (c *JoystickController) ButtonPressed(button ControllerButton) bool {
	if !c.connected {
		return false
	}
	return c.buttonsPressed[button]
}

func (c *JoystickController) ButtonReleased(button ControllerButton) bool {
	if !c.connected {
		return false
	}
	return c.buttonsReleased[button]
}

func (c *JoystickController) ButtonDown(button ControllerButton) bool {
	if !c.connected {
		return false
	}
	return c.buttonsDown[button]
}
