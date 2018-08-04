package graphics

import "github.com/go-gl/mathgl/mgl32"

// Color is a Vec4
type Color mgl32.Vec4

// NewColor creates a new color from the RBGA components
func NewColor(r, g, b, a float32) *Color {
	return &Color{r, g, b, a}
}

// Set the color using new RBGA components
func (c *Color) Set(r, g, b, a float32) {
	c[0] = r
	c[1] = g
	c[2] = b
	c[3] = a
}

// R returns the red component of the color
func (c *Color) R() float32 {
	return c[0]
}

// G returns the green component of the color
func (c *Color) G() float32 {
	return c[1]
}

// B returns the blue component of the color
func (c *Color) B() float32 {
	return c[2]
}

// A returns the alpha component of the color
func (c *Color) A() float32 {
	return c[3]
}
