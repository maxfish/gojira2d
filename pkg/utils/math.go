package utils

import (
	"errors"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

// CircleToPolygon approximate a circle shape with a regular polygon
func CircleToPolygon(center mgl64.Vec2, radius float64, numSegments int, startAngle float64) ([]mgl64.Vec2, error) {
	if radius <= 0 {
		return nil, errors.New("Radius cannot be <=0")
	}
	if numSegments < 3 {
		return nil, errors.New("numSegments must be >= 3")
	}
	point := mgl64.Rotate2D(startAngle).Mul2x1(mgl64.Vec2{radius, 0})
	vertices := make([]mgl64.Vec2, 0, numSegments*2)
	rotation := mgl64.Rotate2D((math.Pi * 2.0) / float64(numSegments))

	for index := 0; index < numSegments; index++ {
		p := point.Add(center)
		vertices = append(vertices, p)
		point = rotation.Mul2x1(point)
	}

	return vertices, nil
}

// GetBoundingBox returns the top left and the bottom right points of the 2D box bounding all the points passed.
func GetBoundingBox(points []mgl64.Vec2) (mgl64.Vec2, mgl64.Vec2) {
	var minX, minY, maxX, maxY float64
	minX = math.MaxFloat32
	minY = math.MaxFloat32
	maxX = -math.MaxFloat32
	maxY = -math.MaxFloat32
	for _, p := range points {
		if p.X() < minX {
			minX = p.X()
		}
		if p.X() > maxX {
			maxX = p.X()
		}
		if p.Y() < minY {
			minY = p.Y()
		}
		if p.Y() > maxY {
			maxY = p.Y()
		}
	}

	return mgl64.Vec2{minX, minY}, mgl64.Vec2{maxX, maxY}
}

func Mat4From64to32Bits(mat mgl64.Mat4) mgl32.Mat4 {
	return mgl32.Mat4{
		float32(mat[0]),
		float32(mat[1]),
		float32(mat[2]),
		float32(mat[3]),
		float32(mat[4]),
		float32(mat[5]),
		float32(mat[6]),
		float32(mat[7]),
		float32(mat[8]),
		float32(mat[9]),
		float32(mat[10]),
		float32(mat[11]),
		float32(mat[12]),
		float32(mat[13]),
		float32(mat[14]),
		float32(mat[15]),
	}
}
