package input

type GameController interface {
	Connected() (bool)
	NumButtons() (int)
	NumAxis() (int)
	NumBalls() (int)
	ButtonPressed(buttonId int) (bool)
	ButtonReleased(buttonId int) (bool)
	ButtonDown(buttonId int) (bool)
	Open(deviceIndex int) (bool)
	Close()
	Update()
	GetAxisValue(axisIndex int) (float32)
	GetAxisDigital(axisIndex int) (float32)
}
