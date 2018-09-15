package graphics

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

func TestCamera2D(t *testing.T) {
	c := NewCamera2D(100, 100, 10)

	c.SetPosition(10, 10)
	expected := mgl64.Mat4FromRows(mgl64.Vec4{0.2, 0, 0, -3}, mgl64.Vec4{0, -0.2, 0, 3}, mgl64.Vec4{0, 0, 0.5, 0}, mgl64.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("SetPosition failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}
	// Checks the 32bit version
	expected32 := mgl32.Mat4FromRows(mgl32.Vec4{0.2, 0, 0, -3}, mgl32.Vec4{0, -0.2, 0, 3}, mgl32.Vec4{0, 0, 0.5, 0}, mgl32.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix32().ApproxEqual(expected32) {
		t.Errorf("ProjectionMatrix32 failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}

	c.SetZoom(5)
	expected = mgl64.Mat4FromRows(mgl64.Vec4{0.1, 0, 0, -2}, mgl64.Vec4{0, -0.1, 0, 2}, mgl64.Vec4{0, 0, 0.5, 0}, mgl64.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("SetZoom failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}

	// Centered
	c.SetCentered(true)
	expected = mgl64.Mat4FromRows(mgl64.Vec4{0.1, 0, 0, -1}, mgl64.Vec4{0, -0.1, 0, 1}, mgl64.Vec4{0, 0, 0.5, 0}, mgl64.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("SetCentered failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}

	// Flipped and centered
	c.SetFlipVertical(true)
	expected = mgl64.Mat4FromRows(mgl64.Vec4{0.1, 0, 0, -1}, mgl64.Vec4{0, 0.1, 0, -1}, mgl64.Vec4{0, 0, 0.5, 0}, mgl64.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("SetFlipVertical failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}

	// Flipped and not centered
	c.SetCentered(false)
	expected = mgl64.Mat4FromRows(mgl64.Vec4{0.1, 0, 0, -2}, mgl64.Vec4{0, 0.1, 0, -2}, mgl64.Vec4{0, 0, 0.5, 0}, mgl64.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("SetFlipVertical failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}

	// SetVisibleArea
	c = NewCamera2D(100, 100, 10)

	// not centered
	c.SetCentered(false)
	c.SetVisibleArea(50, 50, 150, 100)
	expected = mgl64.Mat4FromRows(mgl64.Vec4{0.02, 0, 0, -2}, mgl64.Vec4{0, -0.02, 0, 2}, mgl64.Vec4{0, 0, 0.5, 0}, mgl64.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("SetVisibleArea failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}
	// x1,y1 swapped with x2,y2
	c.SetVisibleArea(150, 100, 50, 50)
	expected = mgl64.Mat4FromRows(mgl64.Vec4{0.02, 0, 0, -2}, mgl64.Vec4{0, -0.02, 0, 2}, mgl64.Vec4{0, 0, 0.5, 0}, mgl64.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("SetVisibleArea failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}

	// centered
	c.SetCentered(true)
	c.SetVisibleArea(50, 50, 150, 100)
	expected = mgl64.Mat4FromRows(mgl64.Vec4{0.02, 0, 0, -2}, mgl64.Vec4{0, -0.02, 0, 1.5}, mgl64.Vec4{0, 0, 0.5, 0}, mgl64.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("SetVisibleArea failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}

	// Translate
	c.Translate(50, 50)
	expected = mgl64.Mat4FromRows(mgl64.Vec4{0.02, 0, 0, -3}, mgl64.Vec4{0, -0.02, 0, 2.5}, mgl64.Vec4{0, 0, 0.5, 0}, mgl64.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("Translate failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}
	c.SetFlipVertical(true)
	c.Translate(50, 50)
	expected = mgl64.Mat4FromRows(mgl64.Vec4{0.02, 0, 0, -4}, mgl64.Vec4{0, 0.02, 0, -1.5}, mgl64.Vec4{0, 0, 0.5, 0}, mgl64.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("Translate failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}

	// Zoom
	c.SetZoom(15)
	var expectedZoom float64 = 15
	if !mgl64.FloatEqual(c.Zoom(), expectedZoom) {
		t.Errorf("Zoom() failed\nexpected\n%f received\n%f", expectedZoom, c.Zoom())
	}
}

func TestCamera2DCoordsConversion(t *testing.T) {
	c := NewCamera2D(100, 100, 10)

	screen := mgl64.Vec2{37, 25}
	worldExpected := mgl64.Vec3{3.7, 7.5, 0}
	world := c.ScreenToWorld(screen)
	if !world.ApproxEqual(worldExpected) {
		t.Errorf("ScreenToWorld() failed\nexpected\n%f received\n%f", worldExpected, world)
	}

	screen2 := c.WorldToScreen(world)
	if !screen2.ApproxEqual(screen) {
		t.Errorf("WorldToScreen() failed\nexpected\n%f received\n%f", screen, screen2)
	}

	c.SetFlipVertical(true)
	screen = mgl64.Vec2{137, 85}
	worldExpected = mgl64.Vec3{13.7, 8.5, 0}
	world = c.ScreenToWorld(screen)
	if !world.ApproxEqual(worldExpected) {
		t.Errorf("ScreenToWorld() failed\nexpected\n%f received\n%f", worldExpected, world)
	}

	screen2 = c.WorldToScreen(world)
	if !screen2.ApproxEqual(screen) {
		t.Errorf("WorldToScreen() failed\nexpected\n%f received\n%f", screen, screen2)
	}

}
