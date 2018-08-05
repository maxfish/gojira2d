package graphics

import (
	"math"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestTransformations(t *testing.T) {
	p := &Primitive2D{}
	p.position = mgl32.Vec3{0, 0, 0}
	p.size = mgl32.Vec2{5, 5}
	p.scale = mgl32.Vec2{1, 1}
	p.rebuildMatrices()

	expected := mgl32.Mat4FromRows(mgl32.Vec4{5, 0, 0, 0}, mgl32.Vec4{0, 5, 0, 0}, mgl32.Vec4{0, 0, 1, 0}, mgl32.Vec4{0, 0, 0, 1})
	if !p.ModelMatrix().ApproxEqual(expected) {
		t.Errorf("Initial setup failed\nexpected\n%s received\n%s", expected.String(), p.modelMatrix.String())
	}

	p.SetScale(mgl32.Vec2{2, 3})
	expected = mgl32.Mat4FromRows(mgl32.Vec4{10, 0, 0, 0}, mgl32.Vec4{0, 15, 0, 0}, mgl32.Vec4{0, 0, 1, 0}, mgl32.Vec4{0, 0, 0, 1})
	if !p.ModelMatrix().ApproxEqual(expected) {
		t.Errorf("SetScale failed\nexpected\n%s received\n%s", expected.String(), p.modelMatrix.String())
	}

	p.SetSize(mgl32.Vec2{10, 10})
	expected = mgl32.Mat4FromRows(mgl32.Vec4{20, 0, 0, 0}, mgl32.Vec4{0, 30, 0, 0}, mgl32.Vec4{0, 0, 1, 0}, mgl32.Vec4{0, 0, 0, 1})
	if !p.ModelMatrix().ApproxEqual(expected) {
		t.Errorf("SetSize failed\nexpected\n%s received\n%s", expected.String(), p.modelMatrix.String())
	}

	// Texture is not set, nothing should change
	p.SetSizeFromTexture()
	if !p.ModelMatrix().ApproxEqual(expected) {
		t.Errorf("SetSizeFromTexture with no texture failed\nexpected\n%s received\n%s", expected.String(), p.modelMatrix.String())
	}

	p.texture = &Texture{0, 20, 20}
	p.SetSizeFromTexture()
	expected = mgl32.Mat4FromRows(mgl32.Vec4{40, 0, 0, 0}, mgl32.Vec4{0, 60, 0, 0}, mgl32.Vec4{0, 0, 1, 0}, mgl32.Vec4{0, 0, 0, 1})
	if !p.ModelMatrix().ApproxEqual(expected) {
		t.Errorf("SetSizeFromTexture failed\nexpected\n%s received\n%s", expected.String(), p.modelMatrix.String())
	}
	p.SetSize(mgl32.Vec2{10, 10})

	p.SetPosition(mgl32.Vec3{10, 10, -1})
	expected = mgl32.Mat4FromRows(mgl32.Vec4{20, 0, 0, 10}, mgl32.Vec4{0, 30, 0, 10}, mgl32.Vec4{0, 0, 1, -1}, mgl32.Vec4{0, 0, 0, 1})
	if !p.ModelMatrix().ApproxEqual(expected) {
		t.Errorf("SetPosition failed\nexpected\n%s received\n%s", expected.String(), p.modelMatrix.String())
	}

	p.SetAnchor(mgl32.Vec2{2, 2})
	expected = mgl32.Mat4FromRows(mgl32.Vec4{20, 0, 0, 6}, mgl32.Vec4{0, 30, 0, 4}, mgl32.Vec4{0, 0, 1, -1}, mgl32.Vec4{0, 0, 0, 1})
	if !p.ModelMatrix().ApproxEqual(expected) {
		t.Errorf("SetAnchor failed\nexpected\n%s received\n%s", expected.String(), p.modelMatrix.String())
	}

	p.SetAnchorToCenter()
	expected = mgl32.Mat4FromRows(mgl32.Vec4{20, 0, 0, 0}, mgl32.Vec4{0, 30, 0, -5}, mgl32.Vec4{0, 0, 1, -1}, mgl32.Vec4{0, 0, 0, 1})
	if !p.ModelMatrix().ApproxEqual(expected) {
		t.Errorf("SetAnchorToCenter failed\nexpected\n%s received\n%s", expected.String(), p.modelMatrix.String())
	}

	p.SetFlipX(true)
	p.SetFlipY(true)
	expected = mgl32.Mat4FromRows(mgl32.Vec4{-20, 0, 0, 20}, mgl32.Vec4{0, -30, 0, 25}, mgl32.Vec4{0, 0, 1, -1}, mgl32.Vec4{0, 0, 0, 1})
	if !p.ModelMatrix().ApproxEqual(expected) {
		t.Errorf("SetFlip failed\nexpected\n%s received\n%s", expected.String(), p.modelMatrix.String())
	}

	p.SetAngle(math.Pi / 4)
	expected = mgl32.Mat4FromRows(mgl32.Vec4{-14.142, 21.213, 0, 6.464}, mgl32.Vec4{-14.142, -21.213, 0, 27.677}, mgl32.Vec4{0, 0, 1, -1}, mgl32.Vec4{0, 0, 0, 1})
	if !p.ModelMatrix().ApproxEqualThreshold(expected, 0.0001) {
		t.Errorf("SetAngle failed\nexpected\n%s received\n%s", expected.String(), p.modelMatrix.String())
	}

	size := mgl32.Vec2{10, 10}
	if !p.Size().ApproxEqual(size) {
		t.Errorf("Size failed")
	}

	var angle float32 = 0.7853
	if !mgl32.FloatEqualThreshold(p.Angle(), angle, 0.0001) {
		t.Errorf("Angle failed\nexpected %f, received %f", angle, p.Angle())
	}

}
