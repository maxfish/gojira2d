package input

type GameController interface {
	Connected() (bool)
	Open(deviceIndex int) (bool)
	Close()
	Update()
	NumButtons() (int)
	NumAxis() (int)
	ButtonPressed(buttonId int) (bool)
	ButtonReleased(buttonId int) (bool)
	ButtonDown(buttonId int) (bool)
	GetAxisValue(axisIndex int) (float32)
	GetAxisDigital(axisIndex int) (float32)
}
