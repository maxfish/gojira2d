package input

type ControllerAxis int
type ControllerButton int

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
	GetAxisValue(axis ControllerAxis) float32
	GetAxisDigital(axis ControllerAxis) float32
}
