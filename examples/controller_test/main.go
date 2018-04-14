package main

import (
	a "gojira2d/pkg/app"
	"gojira2d/pkg/graphics"
	"gojira2d/pkg/input"

	"github.com/go-gl/mathgl/mgl32"
)

func main() {
	app := a.InitApp(640, 480, false, "Controller Test")
	defer a.TerminateApp()

	var joy input.GameController

	// Tries connecting a joystick...
	joy = &input.JoystickController{}
	ok := joy.Open(0)
	if ok {
		println(joy.Description())
	} else {
		// falls back to the keyboard
		keyboard := &input.KeyboardController{}
		keyboard.SetWindow(app.Window)
		keyboard.Open(-1)
		println(keyboard.Description())
		joy = keyboard
	}
	parts := make([]joystickPart, 0, 32)

	// 2 x Sticks
	for i := 0; i < 2; i++ {
		parts = append(parts, NewStick(mgl32.Vec3{200 + float32(i)*240, 240, 0}, i*2, i*2+1))
	}
	// 2 x Triggers
	for i := 0; i < 2; i++ {
		parts = append(parts, NewTrigger(mgl32.Vec3{200 + float32(i)*240, 95, 0}, 4+i))
	}
	// Buttons
	for i := 0; i < joy.NumButtons(); i++ {
		pos := posForButton(i)
		button := NewButton(pos, i)
		parts = append(parts, button)
	}

	app.MainLoop(func(speed float64) {
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

func NewButton(position mgl32.Vec3, buttonIndex int) *Button {
	b := &Button{}
	b.buttonIndex = buttonIndex
	position = position.Add(mgl32.Vec3{0, 0, 1})
	b.shape = graphics.NewRegularPolygonPrimitive(position, 12, 16, false)
	b.shape.SetAnchorToCenter()
	b.shape.SetColor(graphics.Color{0.5, 0.5, 0.5, 1})
	position = position.Add(mgl32.Vec3{0, 0, -2})
	b.shapePressed = graphics.NewRegularPolygonPrimitive(position, 12, 16, true)
	b.shapePressed.SetAnchorToCenter()
	b.shapePressed.SetColor(graphics.Color{0.5, 1, 0.5, 1})
	return b
}

func posForButton(buttonIndex int) mgl32.Vec3 {
	var xLeft float32 = 200
	var xCenter float32 = 320
	var xRight float32 = 440
	var y float32 = 320
	var yCenter float32 = 240
	switch input.ControllerButton(buttonIndex) {
	case input.BUTTON_A:
		return mgl32.Vec3{xRight, y + 60, 0}
	case input.BUTTON_B:
		return mgl32.Vec3{xRight + 30, y + 30, 0}
	case input.BUTTON_X:
		return mgl32.Vec3{xRight - 30, y + 30, 0}
	case input.BUTTON_Y:
		return mgl32.Vec3{xRight, y, 0}
	case input.BUTTON_DIR_PAD_DOWN:
		return mgl32.Vec3{xLeft, y + 60, 0}
	case input.BUTTON_DIR_PAD_RIGHT:
		return mgl32.Vec3{xLeft + 30, y + 30, 0}
	case input.BUTTON_DIR_PAD_LEFT:
		return mgl32.Vec3{xLeft - 30, y + 30, 0}
	case input.BUTTON_DIR_PAD_UP:
		return mgl32.Vec3{xLeft, y, 0}
	case input.BUTTON_BACK:
		return mgl32.Vec3{xCenter - 30, y - 20, 0}
	case input.BUTTON_GUIDE:
		return mgl32.Vec3{xCenter, y - 20, 0}
	case input.BUTTON_START:
		return mgl32.Vec3{xCenter + 30, y - 20, 0}
	case input.BUTTON_LEFT_SHOULDER:
		return mgl32.Vec3{xLeft, yCenter - 90, 0}
	case input.BUTTON_RIGHT_SHOULDER:
		return mgl32.Vec3{xRight, yCenter - 90, 0}
	case input.BUTTON_LEFT_STICK:
		return mgl32.Vec3{xLeft, yCenter, 0}
	case input.BUTTON_RIGHT_STICK:
		return mgl32.Vec3{xRight, yCenter, 0}
	}

	return mgl32.Vec3{float32(buttonIndex) * 30, 450, 0}
}

func (b *Button) Update(controller input.GameController) {
	b.pressed = controller.ButtonDown(input.ControllerButton(b.buttonIndex))
}

func (b *Button) Draw(ctx *graphics.Context) {
	if b.pressed {
		b.shapePressed.EnqueueForDrawing(ctx)
	} else {
		b.shape.EnqueueForDrawing(ctx)
	}
}

type Stick struct {
	axisIndexX input.ControllerAxis
	axisIndexY input.ControllerAxis
	shape      *graphics.Primitive2D
	knob       *graphics.Primitive2D
	position   mgl32.Vec3
	knobPos    mgl32.Vec2
}

func NewStick(position mgl32.Vec3, axisIndexX int, axisIndexY int) *Stick {
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
	absPos = absPos.Add(mgl32.Vec3{x, y, 0})
	s.knob.SetPosition(absPos)
}

func (s *Stick) Draw(ctx *graphics.Context) {
	s.shape.EnqueueForDrawing(ctx)
	s.knob.EnqueueForDrawing(ctx)
}

type Trigger struct {
	axisIndex input.ControllerAxis
	shape     *graphics.Primitive2D
	knob      *graphics.Primitive2D
	position  mgl32.Vec3
	amount    float32
}

func NewTrigger(position mgl32.Vec3, axisIndex int) *Trigger {
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
	amount := (controller.AxisValue(t.axisIndex) + 1) / 2
	scale := amount
	t.knob.SetScale(mgl32.Vec2{scale, scale})
	t.knob.SetColor(graphics.Color{0.2, amount, 0.2, 1})
}

func (t *Trigger) Draw(ctx *graphics.Context) {
	t.shape.EnqueueForDrawing(ctx)
	t.knob.EnqueueForDrawing(ctx)
}
