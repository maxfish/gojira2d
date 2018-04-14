package graphics

import "github.com/go-gl/mathgl/mgl32"

type Color mgl32.Vec4

func NewColor(r, g, b, a float32) *Color {
	return &Color{r, g, b, a}
}

func (c *Color) Set(r, g, b, a float32) {
	c[0] = r
	c[1] = g
	c[2] = b
	c[3] = a
}

func (c *Color) R() float32 {
	return c[0]
}

func (c *Color) G() float32 {
	return c[1]
}

func (c *Color) B() float32 {
	return c[2]
}

func (c *Color) A() float32 {
	return c[3]
}
