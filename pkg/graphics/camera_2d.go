package graphics

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

// Camera2D a Camera based on an orthogonal projection
type Camera2D struct {
	x                float32
	y                float32
	width            float32
	height           float32
	zoom             float32
	centered         bool
	flipVertical     bool
	near             float32
	far              float32
	projectionMatrix mgl32.Mat4
	matrixDirty      bool
}

// NewCamera2D sets up an orthogonal projection camera
func NewCamera2D(width int, height int, zoom float32) *Camera2D {
	c := &Camera2D{
		width:  float32(width),
		height: float32(height),
		zoom:   zoom,
	}
	c.far = -2
	c.near = 2
	c.rebuildMatrix()

	return c
}

// ProjectionMatrix returns the projection matrix of the camera
func (c *Camera2D) ProjectionMatrix() mgl32.Mat4 {
	if c.matrixDirty {
		c.rebuildMatrix()
	}
	return c.projectionMatrix
}

// SetPosition sets the current position of the camera. If the camera is centered, the center will be moving
func (c *Camera2D) SetPosition(x float32, y float32) {
	c.x = x
	c.y = y
	c.matrixDirty = true
}

// SetZoom sets the zoom factor
func (c *Camera2D) SetZoom(zoom float32) {
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
	width := float32(math.Abs(float64(x2 - x1)))
	height := float32(math.Abs(float64(y2 - y1)))
	zoom := float32(math.Min(float64(c.width/width), float64(c.height/height)))
	c.SetZoom(zoom)
	if c.centered {
		c.SetPosition(x1+width/2, y1+height/2)
	} else {
		c.SetPosition(x1, y1)
	}
}

func (c *Camera2D) rebuildMatrix() {
	var left, right, top, bottom float32

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

	c.projectionMatrix = mgl32.Ortho(left, right, top, bottom, c.near, c.far)
	c.matrixDirty = false
}
