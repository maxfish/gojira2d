package graphics

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/maxfish/gojira2d/pkg/utils"
	"math"
)

const MinZoom float64 = 0.01
const MaxZoom float64 = 20

// Camera2D a Camera based on an orthogonal projection
type Camera2D struct {
	x                  float64
	y                  float64
	width              float64
	height             float64
	zoom               float64
	centered           bool
	flipVertical       bool
	near               float64
	far                float64
	projectionMatrix   mgl64.Mat4
	projectionMatrix32 mgl32.Mat4
	matrixDirty        bool
}

// NewCamera2D sets up an orthogonal projection camera
func NewCamera2D(width int, height int, zoom float64) *Camera2D {
	c := &Camera2D{
		width:  float64(width),
		height: float64(height),
		zoom:   zoom,
	}
	c.far = -2
	c.near = 2
	c.matrixDirty = true
	c.rebuildMatrix()

	return c
}

// ProjectionMatrix returns the projection matrix of the camera
func (c *Camera2D) ProjectionMatrix() mgl64.Mat4 {
	c.rebuildMatrix()
	return c.projectionMatrix
}

// ProjectionMatrix32 returns the projection matrix of the camera as mgl32.Mat4
func (c *Camera2D) ProjectionMatrix32() mgl32.Mat4 {
	c.rebuildMatrix()
	return c.projectionMatrix32
}

// SetPosition sets the current position of the camera. If the camera is centered, the center will be moving
func (c *Camera2D) SetPosition(x float64, y float64) {
	c.x = x
	c.y = y
	c.matrixDirty = true
}

// Translate move the camera position by the specified amount
func (c *Camera2D) Translate(x float64, y float64) {
	if c.flipVertical {
		y = -y
	}
	c.x += x
	c.y += y
	c.matrixDirty = true
}

// Zoom returns the current zoom level
func (c *Camera2D) Zoom() float64 {
	return c.zoom
}

// SetZoom sets the zoom factor
func (c *Camera2D) SetZoom(zoom float64) {
	zoom = mgl64.Clamp(zoom, MinZoom, MaxZoom)
	c.zoom = zoom
	c.matrixDirty = true
}

// SetCentered sets the center of the camera to the center of the screen
func (c *Camera2D) SetCentered(centered bool) {
	c.centered = centered
	c.matrixDirty = true
}

// SetFlipVertical sets the orientation of the vertical axis. Pass true to have a cartesian coordinate system
func (c *Camera2D) SetFlipVertical(flip bool) {
	c.flipVertical = flip
	c.matrixDirty = true
}

// SetVisibleArea configures the camera to make the specified area completely visible, position and zoom are changed accordingly
func (c *Camera2D) SetVisibleArea(x1 float32, y1 float32, x2 float32, y2 float32) {
	width := math.Abs(float64(x2 - x1))
	height := math.Abs(float64(y2 - y1))
	zoom := math.Min(float64(c.width/width), float64(c.height/height))
	c.SetZoom(zoom)

	x := math.Min(float64(x1), float64(x2))
	y := math.Min(float64(y1), float64(y2))
	if c.centered {
		c.SetPosition(x+width/2, y+height/2)
	} else {
		c.SetPosition(x, y)
	}
}

func (c *Camera2D) rebuildMatrix() {
	if !c.matrixDirty {
		return
	}
	var left, right, top, bottom float64

	if c.centered {
		halfWidth := c.width / 2 / c.zoom
		halfHeight := c.height / 2 / c.zoom
		left = -halfWidth
		right = halfWidth
		top = halfHeight
		bottom = -halfHeight
	} else {
		right = c.width / c.zoom
		top = c.height / c.zoom
	}

	left += c.x
	right += c.x
	top += c.y
	bottom += c.y

	if c.flipVertical {
		tmp := bottom
		bottom = top
		top = tmp
	}

	c.projectionMatrix = mgl64.Ortho(left, right, top, bottom, c.near, c.far)
	// updates the float32 version
	c.projectionMatrix32 = utils.Mat4From64to32Bits(c.projectionMatrix)
	c.matrixDirty = false
}
