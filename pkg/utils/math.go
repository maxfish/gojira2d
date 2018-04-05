package utils

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

func CircleToPolygon(center mgl32.Vec2, radius float32, numSegments int, startAngle float32) ([]mgl32.Vec2) {
	point := mgl32.Rotate2D(startAngle).Mul2x1(mgl32.Vec2{radius, 0})
	vertices := make([]mgl32.Vec2, 0, numSegments*2)
	rotation := mgl32.Rotate2D((math.Pi * 2.0) / float32(numSegments))

	for index := 0; index < numSegments; index++ {
		p := point.Add(center)
		vertices = append(vertices, p)
		point = rotation.Mul2x1(point)
	}

	return vertices
}
