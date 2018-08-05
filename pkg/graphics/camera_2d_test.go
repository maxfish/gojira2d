package graphics

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestCenteredCamera2D(t *testing.T) {
	c := NewCamera2D(100, 100, 1, true)

	c.SetPosition(10, 10)
	expected := mgl32.Mat4FromRows(mgl32.Vec4{0.02, 0, 0, -0.2}, mgl32.Vec4{0, -0.02, 0, 0.2}, mgl32.Vec4{0, 0, 0.5, 0}, mgl32.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("SetPosition failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}

	c.SetZoom(5)
	expected = mgl32.Mat4FromRows(mgl32.Vec4{0.1, 0, 0, -1}, mgl32.Vec4{0, -0.1, 0, 1}, mgl32.Vec4{0, 0, 0.5, 0}, mgl32.Vec4{0, 0, 0, 1})
	if !c.ProjectionMatrix().ApproxEqual(expected) {
		t.Errorf("SetZoom failed\nexpected\n%s received\n%s", expected.String(), c.projectionMatrix.String())
	}
}

func TestCamera2D(t *testing.T) {
	c := NewCamera2D(100, 100, 10, false)

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
}
