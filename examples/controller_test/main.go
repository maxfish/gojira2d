package main

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/maxfish/gojira2d/pkg/app"
	"github.com/maxfish/gojira2d/pkg/graphics"
	"github.com/maxfish/gojira2d/pkg/input"
)

func main() {
	app.Init(640, 480, "Controller Test")
	defer app.Terminate()

	var joy input.GameController

	// Tries getting a joystick...
	joy = &input.JoystickController{}
	joy.Open(0)
	if !joy.Connected() {
		joy.Close()
		// falls back to the keyboard
		keyboard := &input.KeyboardController{}
		keyboard.Open(-1)
		joy = keyboard
	}
	parts := make([]joystickPart, 0, 32)

	// 2 x Sticks
	for i := 0; i < 2; i++ {
		parts = append(parts, NewStick(mgl64.Vec3{200 + float64(i)*240, 240, 0}, i*2, i*2+1))
	}
	// 2 x Triggers
	for i := 0; i < 2; i++ {
		parts = append(parts, NewTrigger(mgl64.Vec3{200 + float64(i)*240, 95, 0}, 4+i))
	}
	// Buttons
	for i := 0; i < joy.NumButtons(); i++ {
		pos := posForButton(i)
		button := NewButton(pos, i)
		parts = append(parts, button)
	}

	app.MainLoop(func(deltaTimeMs float64) {
		joy.Update()
		for _, p := range parts {
			p.Update(joy)
		}
	}, func() {
		for _, p := range parts {
			p.Draw(app.Context)
		}
	})
}

type joystickPart interface {
	Draw(*graphics.Context)
	Update(input.GameController)
}

type Button struct {
	buttonIndex  int
	shape        *graphics.Primitive2D
	shapePressed *graphics.Primitive2D
	pressed      bool
}

func NewButton(position mgl64.Vec3, buttonIndex int) *Button {
	b := &Button{}
	b.buttonIndex = buttonIndex
	position = position.Add(mgl64.Vec3{0, 0, 1})
	b.shape = graphics.NewRegularPolygonPrimitive(position, 12, 16, false)
	b.shape.SetAnchorToCenter()
	b.shape.SetColor(graphics.Color{0.5, 0.5, 0.5, 1})
	position = position.Add(mgl64.Vec3{0, 0, -2})
	b.shapePressed = graphics.NewRegularPolygonPrimitive(position, 12, 16, true)
	b.shapePressed.SetAnchorToCenter()
	b.shapePressed.SetColor(graphics.Color{0.5, 1, 0.5, 1})
	return b
}

func posForButton(buttonIndex int) mgl64.Vec3 {
	var xLeft float64 = 200
	var xCenter float64 = 320
	var xRight float64 = 440
	var y float64 = 320
	var yCenter float64 = 240
	switch input.ControllerButton(buttonIndex) {
	case input.ButtonA:
		return mgl64.Vec3{xRight, y + 60, 0}
	case input.ButtonB:
		return mgl64.Vec3{xRight + 30, y + 30, 0}
	case input.ButtonX:
		return mgl64.Vec3{xRight - 30, y + 30, 0}
	case input.ButtonY:
		return mgl64.Vec3{xRight, y, 0}
	case input.ButtonDirPadDown:
		return mgl64.Vec3{xLeft, y + 60, 0}
	case input.ButtonDirPadRight:
		return mgl64.Vec3{xLeft + 30, y + 30, 0}
	case input.ButtonDirPadLeft:
		return mgl64.Vec3{xLeft - 30, y + 30, 0}
	case input.ButtonDirPadUp:
		return mgl64.Vec3{xLeft, y, 0}
	case input.ButtonBack:
		return mgl64.Vec3{xCenter - 30, y - 20, 0}
	case input.ButtonGuide:
		return mgl64.Vec3{xCenter, y - 20, 0}
	case input.ButtonStart:
		return mgl64.Vec3{xCenter + 30, y - 20, 0}
	case input.ButtonLeftShoulder:
		return mgl64.Vec3{xLeft, yCenter - 90, 0}
	case input.ButtonRightShoulder:
		return mgl64.Vec3{xRight, yCenter - 90, 0}
	case input.ButtonLeftStick:
		return mgl64.Vec3{xLeft, yCenter, 0}
	case input.ButtonRightStick:
		return mgl64.Vec3{xRight, yCenter, 0}
	}

	return mgl64.Vec3{float64(buttonIndex) * 30, 450, 0}
}

func (b *Button) Update(controller input.GameController) {
	b.pressed = controller.ButtonDown(input.ControllerButton(b.buttonIndex))
}

func (b *Button) Draw(ctx *graphics.Context) {
	if b.pressed {
		b.shapePressed.Draw(ctx)
	} else {
		b.shape.Draw(ctx)
	}
}

type Stick struct {
	axisIndexX input.ControllerAxis
	axisIndexY input.ControllerAxis
	shape      *graphics.Primitive2D
	knob       *graphics.Primitive2D
	position   mgl64.Vec3
	knobPos    mgl64.Vec2
}

func NewStick(position mgl64.Vec3, axisIndexX int, axisIndexY int) *Stick {
	b := &Stick{}
	b.position = position
	b.axisIndexX = input.ControllerAxis(axisIndexX)
	b.axisIndexY = input.ControllerAxis(axisIndexY)
	b.shape = graphics.NewRegularPolygonPrimitive(position, 36, 24, false)
	b.shape.SetAnchorToCenter()
	b.shape.SetColor(graphics.Color{0.3, 0.3, 0.3, 1})
	b.knob = graphics.NewRegularPolygonPrimitive(position, 36, 24, true)
	b.knob.SetAnchorToCenter()
	b.knob.SetColor(graphics.Color{0.8, 0.8, 0.8, 1})
	return b
}

func (s *Stick) Update(controller input.GameController) {
	x := controller.AxisValue(s.axisIndexX) * 30
	y := controller.AxisValue(s.axisIndexY) * 30
	absPos := s.position
	absPos = absPos.Add(mgl64.Vec3{x, y, 0})
	s.knob.SetPosition(absPos)
}

func (s *Stick) Draw(ctx *graphics.Context) {
	s.shape.Draw(ctx)
	s.knob.Draw(ctx)
}

type Trigger struct {
	axisIndex input.ControllerAxis
	shape     *graphics.Primitive2D
	knob      *graphics.Primitive2D
	position  mgl64.Vec3
	amount    float64
}

func NewTrigger(position mgl64.Vec3, axisIndex int) *Trigger {
	t := &Trigger{}
	t.position = position
	t.axisIndex = input.ControllerAxis(axisIndex)
	t.shape = graphics.NewRegularPolygonPrimitive(position, 24, 20, false)
	t.shape.SetAnchorToCenter()
	t.shape.SetColor(graphics.Color{0.3, 0.3, 0.3, 1})
	t.knob = graphics.NewRegularPolygonPrimitive(position, 24, 20, true)
	t.knob.SetAnchorToCenter()
	return t
}

func (t *Trigger) Update(controller input.GameController) {
	amount := (controller.AxisValue(t.axisIndex) + 1) / 2.0
	scale := amount
	t.knob.SetScale(mgl64.Vec2{scale, scale})
	t.knob.SetColor(graphics.Color{0.2, float32(amount), 0.2, 1})
}

func (t *Trigger) Draw(ctx *graphics.Context) {
	t.shape.Draw(ctx)
	t.knob.Draw(ctx)
}
