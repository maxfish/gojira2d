package input

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type JoystickController struct {
	GameController
	connected       bool
	deviceIndex     glfw.Joystick
	name            string
	axes            []float32
	rawButtons      []byte
	buttonsPressed  []bool
	buttonsReleased []bool
	buttonsDown     []bool
}

func (c *JoystickController) Open(deviceIndex int) (bool) {
	if !glfw.JoystickPresent(glfw.Joystick(deviceIndex)) {
		c.Close()
		return false
	}

	c.connected = true
	c.deviceIndex = glfw.Joystick(deviceIndex)
	c.name = glfw.GetJoystickName(c.deviceIndex)
	c.Update()

	return true
}

func (c *JoystickController) Close() {
	c.connected = false
	c.deviceIndex = -1
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
	c.rawButtons = glfw.GetJoystickButtons(c.deviceIndex)
	if c.buttonsDown == nil {
		// builds the slices
		c.buttonsDown = make([]bool, len(c.rawButtons))
		c.buttonsPressed = make([]bool, len(c.rawButtons))
		c.buttonsReleased = make([]bool, len(c.rawButtons))
	}

	for i := 0; i < len(c.rawButtons); i++ {
		isDown := c.rawButtons[i] > 0
		c.buttonsPressed[i] = false
		c.buttonsReleased[i] = false
		if isDown {
			if !c.buttonsDown[i] {
				c.buttonsPressed[i] = true
			}
		} else {
			if c.buttonsDown[i] {
				c.buttonsReleased[i] = true
			}
		}
		c.buttonsDown[i] = isDown
	}

	// Axes
	c.axes = glfw.GetJoystickAxes(c.deviceIndex)
}

func (c *JoystickController) GetAxisValue(axisIndex int) (float32) {
	return c.axes[axisIndex]
}

func (c *JoystickController) GetAxisDigital(axisIndex int) (float32) {
	// TODO: Define dead zone...
	return 0
}

func (c *JoystickController) Connected() (bool) {
	return c.connected
}

func (c *JoystickController) NumButtons() (int) {
	return len(c.rawButtons)
}

func (c *JoystickController) NumAxis() (int) {
	return len(c.axes)
}

func (c *JoystickController) ButtonPressed(buttonId int) (bool) {
	if !c.connected {
		return false
	}
	return c.buttonsPressed[buttonId]
}

func (c *JoystickController) ButtonReleased(buttonId int) (bool) {
	if !c.connected {
		return false
	}
	return c.buttonsReleased[buttonId]
}

func (c *JoystickController) ButtonDown(buttonId int) (bool) {
	if !c.connected {
		return false
	}
	return c.buttonsDown[buttonId]
}
