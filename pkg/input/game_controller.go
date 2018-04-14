package input

type ControllerButton int

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
)

type ControllerAxis int

const (
	AXIS_LEFT_X ControllerAxis = iota
	AXIS_LEFT_Y
	AXIS_RIGHT_X
	AXIS_RIGHT_Y
	AXIS_TRIGGER_LEFT
	AXIS_TRIGGER_RIGHT
)

type GameController interface {
	Connected() bool
	Open(deviceIndex int) bool
	Close()
	Update()
	NumButtons() int
	NumAxis() int
	ButtonPressed(button ControllerButton) bool
	ButtonReleased(button ControllerButton) bool
	ButtonDown(button ControllerButton) bool
	AxisValue(axis ControllerAxis) float32
	AxisDigitalValue(axis ControllerAxis) int
}
