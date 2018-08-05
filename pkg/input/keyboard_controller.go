package input

import (
	"fmt"
	"log"

	"github.com/go-gl/glfw/v3.2/glfw"
)

// KeyboardController A GameController that uses the keyboard to simulate a joystick
type KeyboardController struct {
	GameController
	connected       bool
	numButtons      int
	numAxes         int
	buttonsPressed  []bool
	buttonsReleased []bool
	buttonsDown     []bool
	buttonsRaw      []bool
	axes            []float32
	mapping         *GameControllerMapping
	keyMapping      map[glfw.Key]int
}

// Open initializes the keyboard. The parameter is ignored
func (c *KeyboardController) Open(_ int) bool {
	if !IsKeyboardFree() {
		log.Print("The keyboard is already in use")
		return false
	}

	c.connected = true
	c.numButtons = 15 // Xbox360
	c.numAxes = 2     // Only the left stick
	c.SetMapping(&MappingKeyboard)

	// Build the slices
	c.buttonsDown = make([]bool, c.numButtons)
	c.buttonsPressed = make([]bool, c.numButtons)
	c.buttonsReleased = make([]bool, c.numButtons)
	c.buttonsRaw = make([]bool, c.numButtons)
	c.axes = make([]float32, c.numAxes)

	RegisterKeyCallback(func(w *glfw.Window, key glfw.Key, scanCode int, action glfw.Action, mods glfw.ModifierKey) {
		if index, ok := c.keyMapping[key]; ok {
			if action == glfw.Press {
				c.buttonsRaw[index] = true
			} else if action == glfw.Release {
				c.buttonsRaw[index] = false
			}
		}
	})

	return true
}

// Close disables the keyboard callback and resets all the data
func (c *KeyboardController) Close() {
	UnregisterKeyCallback()
	c.connected = false
	c.buttonsDown = nil
	c.buttonsPressed = nil
	c.buttonsReleased = nil
	c.buttonsRaw = nil
	c.axes = nil
}

// Update gets the input from the keyboard and sets the game controller accordingly
func (c *KeyboardController) Update() {
	if !c.connected {
		return
	}

	// Buttons
	for i, button := range c.buttonsRaw {
		isDown := button
		c.buttonsPressed[i] = isDown && !c.buttonsDown[i]
		c.buttonsReleased[i] = !isDown && c.buttonsDown[i]
		c.buttonsDown[i] = isDown
	}

	// Axes
	if c.buttonsDown[BUTTON_DIR_PAD_LEFT] {
		c.axes[0] = -1
	} else if c.buttonsDown[BUTTON_DIR_PAD_RIGHT] {
		c.axes[0] = 1
	} else {
		c.axes[0] = 0
	}
	if c.buttonsDown[BUTTON_DIR_PAD_UP] {
		c.axes[1] = -1
	} else if c.buttonsDown[BUTTON_DIR_PAD_DOWN] {
		c.axes[1] = 1
	} else {
		c.axes[1] = 0
	}
}

// AxisValue returns the current value, from -1 to 1, of the axis
func (c *KeyboardController) AxisValue(axis ControllerAxis) float32 {
	if int(axis) >= c.numAxes {
		return 0
	}
	return c.axes[axis]
}

// AxisDigitalValue returns a digital value for the axis
func (c *KeyboardController) AxisDigitalValue(axis ControllerAxis) int {
	return 0
}

// Connected returns if this controller is connected and initialized
func (c *KeyboardController) Connected() bool {
	return c.connected
}

// NumButtons number of the buttons this controller is simulating
func (c *KeyboardController) NumButtons() int {
	return c.numButtons
}

// NumAxis number of the axes this controller is simulating
func (c *KeyboardController) NumAxis() int {
	return c.numAxes
}

// ButtonPressed checks if a button has been pressed since the last frame
func (c *KeyboardController) ButtonPressed(button ControllerButton) bool {
	if !c.connected {
		return false
	}
	return c.buttonsPressed[button]
}

// ButtonReleased checks if a button has been released since the last frame
func (c *KeyboardController) ButtonReleased(button ControllerButton) bool {
	if !c.connected {
		return false
	}
	return c.buttonsReleased[button]
}

// ButtonDown checks if a button is currently pressed
func (c *KeyboardController) ButtonDown(button ControllerButton) bool {
	if !c.connected {
		return false
	}
	return c.buttonsDown[button]
}

// Description identification string of this controller
func (c *KeyboardController) Description() string {
	return fmt.Sprintf("joystick:'Keyboard' buttons:%d axes:0", c.numButtons)
}

// SetMapping maps the keys to the virtual controller inputs
func (c *KeyboardController) SetMapping(mapping *GameControllerMapping) {
	c.mapping = mapping
	c.keyMapping = make(map[glfw.Key]int)
	for i, key := range c.mapping.buttons {
		c.keyMapping[glfw.Key(key)] = i
	}
}
