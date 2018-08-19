package graphics

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestCamera2D(t *testing.T) {
	c := NewCamera2D(100, 100, 10)

	c.SetPosition(10, 10)
	expected := mgl32.Mat4FromRows(mgl32.Vec4{0.2, 0, 0, -3}, mgl32.Vec4{0, -0.2, 0, 3}, mgl32.Vec4{0, 0, 0.5, 0}, mgl32.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("SetPosition failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}

	c.SetZoom(5)
	expected = mgl32.Mat4FromRows(mgl32.Vec4{0.1, 0, 0, -2}, mgl32.Vec4{0, -0.1, 0, 2}, mgl32.Vec4{0, 0, 0.5, 0}, mgl32.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("SetZoom failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}

	// Centered
	c.SetCentered(true)
	expected = mgl32.Mat4FromRows(mgl32.Vec4{0.1, 0, 0, -1}, mgl32.Vec4{0, -0.1, 0, 1}, mgl32.Vec4{0, 0, 0.5, 0}, mgl32.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("SetCentered failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}

	// Flipped and centered
	c.SetFlipVertical(true)
	expected = mgl32.Mat4FromRows(mgl32.Vec4{0.1, 0, 0, -1}, mgl32.Vec4{0, 0.1, 0, -1}, mgl32.Vec4{0, 0, 0.5, 0}, mgl32.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("SetFlipVertical failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}

	// Flipped and not centered
	c.SetCentered(false)
	expected = mgl32.Mat4FromRows(mgl32.Vec4{0.1, 0, 0, -2}, mgl32.Vec4{0, 0.1, 0, -2}, mgl32.Vec4{0, 0, 0.5, 0}, mgl32.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("SetFlipVertical failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}

	// SetVisibleArea
	c = NewCamera2D(100, 100, 10)

	// not centered
	c.SetCentered(false)
	c.SetVisibleArea(50, 50, 150, 100)
	expected = mgl32.Mat4FromRows(mgl32.Vec4{0.02, 0, 0, -2}, mgl32.Vec4{0, -0.02, 0, 2}, mgl32.Vec4{0, 0, 0.5, 0}, mgl32.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("SetVisibleArea failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}

	// centered
	c.SetCentered(true)
	c.SetVisibleArea(50, 50, 150, 100)
	expected = mgl32.Mat4FromRows(mgl32.Vec4{0.02, 0, 0, -2}, mgl32.Vec4{0, -0.02, 0, 1.5}, mgl32.Vec4{0, 0, 0.5, 0}, mgl32.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("SetVisibleArea failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}
}
