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
	}

	for _, test := range tests {
		topLeft, bottomRight := GetBoundingBox(test.points)
		if topLeft != test.topLeft || bottomRight != test.bottomRight {
			t.Errorf("Got %v %v, expecting %v %v", topLeft, bottomRight, test.topLeft, test.bottomRight)
		}
	}
}
