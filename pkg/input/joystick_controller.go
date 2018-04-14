package input

import (
	"fmt"

	"regexp"

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
	mapping         *GameControllerMapping
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

	c.findMapping()

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
	if int(axis) >= c.numAxes {
		return 0
	}
	axisIndex := c.axisFromMapping(axis)
	return c.axes[axisIndex]
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
	buttonIndex := c.buttonFromMapping(button)
	return c.buttonsPressed[buttonIndex]
}

func (c *JoystickController) ButtonReleased(button ControllerButton) bool {
	if !c.connected {
		return false
	}
	buttonIndex := c.buttonFromMapping(button)
	return c.buttonsReleased[buttonIndex]
}

func (c *JoystickController) ButtonDown(button ControllerButton) bool {
	if !c.connected {
		return false
	}
	buttonIndex := c.buttonFromMapping(button)
	if buttonIndex >= c.numButtons {
		return false
	}
	return c.buttonsDown[buttonIndex]
}

func (c *JoystickController) Description() string {
	return fmt.Sprintf("joystick:'%s' buttons:%d axes:%d", c.name, c.numButtons, c.numAxes)
}

func (c *JoystickController) SetMapping(mapping *GameControllerMapping) {
	c.mapping = mapping
}

func (c *JoystickController) findMapping() {
	for _, mapping := range GameControllerMappings {
		r, _ := regexp.Compile(mapping.nameRegEx)
		if r.MatchString(c.name) == true && len(mapping.buttons) == c.numButtons && len(mapping.axes) == c.numAxes {
			c.SetMapping(mapping)
			return
		}
	}
	//	Couldn't find a mapping, pick the XBox 360 one
	c.SetMapping(&MappingXBox360)
}

func (c *JoystickController) buttonFromMapping(index ControllerButton) int {
	if int(index) >= c.numButtons {
		return 0
	}
	return c.mapping.buttons[int(index)]
}

func (c *JoystickController) axisFromMapping(index ControllerAxis) int {
	if int(index) >= c.numAxes {
		return 0
	}
	return c.mapping.axes[int(index)]
}
