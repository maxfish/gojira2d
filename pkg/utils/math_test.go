package utils

import (
	"math"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestCircleToPolygon(t *testing.T) {
	var tests = []struct {
		center            mgl32.Vec2
		radius            float32
		numSegments       int
		startAngle        float32
		expectingError    bool
		numReturnedPoints int
		returnedPoints    []mgl32.Vec2
	}{
		{mgl32.Vec2{100, 100}, 50, 4, 0, false, 4,
			[]mgl32.Vec2{{150, 100}, {100, 150}, {50, 100}, {100, 50}}},
		{mgl32.Vec2{100, 100}, 50, 4, 45, false, 4,
			[]mgl32.Vec2{{126.26, 142.54}, {57.45, 126.26}, {73.73, 57.45}, {142.54, 73.73}}},
		{mgl32.Vec2{1, 1}, 5, 2, 0, true, -1, nil},
		{mgl32.Vec2{1, 1}, -5, 10, 0, true, -1, nil},
	}

	for _, test := range tests {
		points, err := CircleToPolygon(test.center, test.radius, test.numSegments, test.startAngle)
		if err != nil && !test.expectingError {
			t.Errorf("Not expecting error %v", err)
		}
		if err == nil && test.expectingError {
			t.Errorf("Was expecting an error, none returned")
		}
		if err != nil {
			continue
		}
		if len(points) != test.numReturnedPoints {
			t.Errorf("Got %d points, expecting %d", len(points), test.numReturnedPoints)
		}
		for i, p := range points {
			check := math.Round(float64(p[0])) == math.Round(float64(test.returnedPoints[i][0]))
			check = check && math.Round(float64(p[1])) == math.Round(float64(test.returnedPoints[i][1]))
			if !check {
				t.Errorf("Point returned %v is different from the expected %v", p, test.returnedPoints[i])
				break
			}
		}
	}
}

func TestGetBoundingBox(t *testing.T) {
	var tests = []struct {
		points      []mgl32.Vec2
		topLeft     mgl32.Vec2
		bottomRight mgl32.Vec2
	}{
		{[]mgl32.Vec2{{5, 5}, {-10, -10}, {20, 20}}, mgl32.Vec2{-10, -10}, mgl32.Vec2{20, 20}},
		{[]mgl32.Vec2{{-100, -5}, {-80, 50}, {-4, 20}}, mgl32.Vec2{-100, -5}, mgl32.Vec2{-4, 50}},
		{[]mgl32.Vec2{{3, 5}, {3, 6}, {2, 6}, {2, 5}}, mgl32.Vec2{2, 5}, mgl32.Vec2{3, 6}},
	}

	for _, test := range tests {
		topLeft, bottomRight := GetBoundingBox(test.points)
		if topLeft != test.topLeft || bottomRight != test.bottomRight {
			t.Errorf("Got %v %v, expecting %v %v", topLeft, bottomRight, test.topLeft, test.bottomRight)
		}
	}
}

func TestMatrixConversion(t *testing.T) {
	matrix64 := mgl64.Mat4FromRows(mgl64.Vec4{6, -7.5, 8, -9}, mgl64.Vec4{10, -11, 12.5, -13}, mgl64.Vec4{14, -15, 16, -17.5}, mgl64.Vec4{18.5, -19, 20, -21})
	matrix32 := Mat4From64to32Bits(matrix64)
	matrix64From32 := mgl64.Mat4{
		float64(matrix32[0]), float64(matrix32[1]), float64(matrix32[2]), float64(matrix32[3]),
		float64(matrix32[4]), float64(matrix32[5]), float64(matrix32[6]), float64(matrix32[7]),
		float64(matrix32[8]), float64(matrix32[9]), float64(matrix32[10]), float64(matrix32[11]),
		float64(matrix32[12]), float64(matrix32[13]), float64(matrix32[14]), float64(matrix32[15]),
	}

	if !matrix64.ApproxEqual(matrix64From32) {
		t.Errorf("Mat4From64to32Bits failed\nexpected\n%s received\n%s", matrix32.String(), matrix64.String())
	}
}
