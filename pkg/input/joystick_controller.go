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

var (
	JoystickControllers map[int]*JoystickController
)

func init() {
	// Map keeping track of the connected joysticks
	JoystickControllers = make(map[int]*JoystickController, MaxNumJoysticks)

	// Attach the status change callback
	glfw.SetJoystickCallback(func(joy, event int) {
		if glfw.MonitorEvent(event) == glfw.Connected {
			// The joystick was connected
			fmt.Printf("Joystick #%d: plugged in", joy)
			if JoystickControllers[joy] != nil {
				JoystickControllers[joy].Open(joy)
			}
		} else if glfw.MonitorEvent(event) == glfw.Disconnected {
			// The joystick was disconnected
			fmt.Printf("Joystick #%d: plugged out", joy)
			if JoystickControllers[joy] != nil && JoystickControllers[joy].Connected() {
				JoystickControllers[joy].pluggedOut()
			}
		}
	})
}

func (c *JoystickController) Open(deviceIndex int) bool {
	if c.connected {
		fmt.Printf("Joystick already open on device #%d", c.joystick)
		return true
	}

	if JoystickControllers[deviceIndex] != nil && JoystickControllers[deviceIndex] != c {
		fmt.Printf("Another joystick is associated to index #%d", deviceIndex)
		return false
	}

	JoystickControllers[deviceIndex] = c
	c.joystick = glfw.Joystick(deviceIndex)

	// The joystick is currently not connected but it might be plugged in later
	if !glfw.JoystickPresent(glfw.Joystick(deviceIndex)) {
		return true
	}

	c.connected = true
	c.name = glfw.GetJoystickName(c.joystick)

	// Get the num of buttons and axes
	c.numButtons = len(glfw.GetJoystickButtons(c.joystick))
	c.numAxes = len(glfw.GetJoystickAxes(c.joystick))
	// Build the slices
	c.buttonsDown = make([]bool, c.numButtons)
	c.buttonsPressed = make([]bool, c.numButtons)
	c.buttonsReleased = make([]bool, c.numButtons)

	fmt.Printf("Joystick #%d: opened. %s", c.joystick, c.Description())
	c.findMapping()

	return true
}

func (c *JoystickController) Close() {
	fmt.Printf("Joystick #%d: closed", c.joystick)
	JoystickControllers[int(c.joystick)] = nil
	c.pluggedOut()
}

func (c *JoystickController) pluggedOut() {
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

func (c *JoystickController) NumAxes() int {
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
	return fmt.Sprintf("name:'%s' buttons:%d axes:%d", c.name, c.numButtons, c.numAxes)
}

func (c *JoystickController) SetMapping(mapping *GameControllerMapping) {
	c.mapping = mapping
}

func (c *JoystickController) findMapping() {
	for _, mapping := range GameControllerMappings {
		r, _ := regexp.Compile(mapping.nameRegEx)
		if r.MatchString(c.name) == true && len(mapping.buttons) == c.numButtons && len(mapping.axes) == c.numAxes {
			c.SetMapping(mapping)
			fmt.Printf("Joystick #%d: mapping found", c.joystick)
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
