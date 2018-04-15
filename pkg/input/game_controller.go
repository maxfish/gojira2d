package input

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type ControllerButton int
type ControllerAxis int

const (
	BUTTON_A ControllerButton = iota
	BUTTON_B
	BUTTON_X
	BUTTON_Y
	BUTTON_BACK
	BUTTON_GUIDE
	BUTTON_START
	BUTTON_LEFT_STICK
	BUTTON_RIGHT_STICK
	BUTTON_LEFT_SHOULDER
	BUTTON_RIGHT_SHOULDER
	BUTTON_DIR_PAD_UP
	BUTTON_DIR_PAD_DOWN
	BUTTON_DIR_PAD_LEFT
	BUTTON_DIR_PAD_RIGHT
	AXIS_LEFT_X ControllerAxis = iota
	AXIS_LEFT_Y
	AXIS_RIGHT_X
	AXIS_RIGHT_Y
	AXIS_TRIGGER_LEFT
	AXIS_TRIGGER_RIGHT

	MAX_NUM_JOYSTICKS = glfw.JoystickLast
)

var (
	GameControllerMappings []*GameControllerMapping
	MappingXBox360         GameControllerMapping
	MappingPS4             GameControllerMapping
	MappingKeyboard        GameControllerMapping
)

type GameController interface {
	Connected() bool
	Open(deviceIndex int) bool
	Close()
	Update()
	NumButtons() int
	NumAxes() int
	ButtonPressed(button ControllerButton) bool
	ButtonReleased(button ControllerButton) bool
	ButtonDown(button ControllerButton) bool
	AxisValue(axis ControllerAxis) float32
	AxisDigitalValue(axis ControllerAxis) int
	SetMapping(mapping *GameControllerMapping)
	Description() string
}

type GameControllerMapping struct {
	nameRegEx string
	buttons   []int
	axes      []int
}

func (g *GameControllerMapping) set(nameRegEx string, buttons []int, axes []int) {
	g.nameRegEx = nameRegEx
	g.buttons = buttons
	g.axes = axes
}

func init() {
	// Keyboard as controller. This mapping should not be used by a JoystickController
	MappingKeyboard.set("<None>",
		[]int{
			int(glfw.KeyA), int(glfw.KeyS), int(glfw.KeyD), int(glfw.KeyF),
			int(glfw.Key1), int(glfw.Key2), int(glfw.Key3),
			int(glfw.KeyQ), int(glfw.KeyR), int(glfw.KeyW), int(glfw.KeyE),
			int(glfw.KeyUp), int(glfw.KeyDown), int(glfw.KeyLeft), int(glfw.KeyRight),
		},
		[]int{})

	// List of the all the mappings but the keyboard ones
	GameControllerMappings = make([]*GameControllerMapping, 0, 10)

	// PS4 controller with USB cable (MacOS)
	MappingPS4.set(".*Wireless Controller.*", []int{1, 2, 0, 3, 8, 12, 9, 10, 11, 4, 5, 14, 16, 17, 15, 13, 6, 7}, []int{0, 1, 2, 3, 4, 5})
	GameControllerMappings = append(GameControllerMappings, &MappingPS4)

	// XBox 360 wired controller (MacOS)
	MappingXBox360.set(".*Xbox 360.*", []int{11, 12, 13, 14, 5, 10, 4, 6, 7, 8, 9, 0, 1, 2, 3}, []int{0, 1, 2, 3, 4, 5})
	GameControllerMappings = append(GameControllerMappings, &MappingXBox360)
}
