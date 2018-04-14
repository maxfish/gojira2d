package utils

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

func CircleToPolygon(center mgl32.Vec2, radius float32, numSegments int, startAngle float32) []mgl32.Vec2 {
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

func GetBoundingBox(points []mgl32.Vec2) (mgl32.Vec2, mgl32.Vec2) {
	var minX, minY, maxX, maxY float32
	minX = math.MaxFloat32
	minY = math.MaxFloat32
	maxX = -math.MaxFloat32
	maxY = -math.MaxFloat32
	for _, p := range points {
		if p.X() < minX {
			minX = p.X()
		}
		if p.Y() > maxX {
			maxX = p.X()
		}
		if p.Y() < minY {
			minY = p.Y()
		}
		if p.Y() > maxY {
			maxY = p.Y()
		}
	}

	return mgl32.Vec2{minX, minY}, mgl32.Vec2{maxX, maxY}
}
