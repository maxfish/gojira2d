package input

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

// ControllerButton type representing a button of the input device
type ControllerButton int

// ControllerAxis type representing an axis of the input device
type ControllerAxis int

// Constants for all the buttons and axes
const (
	ButtonA ControllerButton = iota
	ButtonB
	ButtonX
	ButtonY
	ButtonBack
	ButtonGuide
	ButtonStart
	ButtonLeftStick
	ButtonRightStick
	ButtonLeftShoulder
	ButtonRightShoulder
	ButtonDirPadUp
	ButtonDirPadDown
	ButtonDirPadLeft
	ButtonDirPadRight
	AxisLeftX ControllerAxis = iota
	AxisLeftY
	AxisRightX
	AxisRightY
	AxisTriggerLeft
	AxisTriggerRight

	// MaxNumJoysticks max num of joysticks allowed by glfw
	MaxNumJoysticks = glfw.JoystickLast
)

// A set of basic mappings
var (
	GameControllerMappings []*GameControllerMapping
	MappingXBox360         GameControllerMapping
	MappingPS4             GameControllerMapping
	MappingKeyboard        GameControllerMapping
)

// GameController represents a physical input device
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
	AxisValue(axis ControllerAxis) float64
	AxisDigitalValue(axis ControllerAxis) int
	SetMapping(mapping *GameControllerMapping)
	Description() string
}

// GameControllerMapping a axes/buttons mapping for a specific device
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
