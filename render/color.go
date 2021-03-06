package render

import (
	"fmt"
	"image/color"

	"github.com/johanhenriksson/goworld/math/byte4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
)

// Predefined Colors
var (
	White       = Color{1, 1, 1, 1}
	Black       = Color{0, 0, 0, 1}
	Red         = Color{1, 0, 0, 1}
	Green       = Color{0, 1, 0, 1}
	Blue        = Color{0, 0, 1, 1}
	Purple      = Color{1, 0, 1, 1}
	Yellow      = Color{1, 1, 0, 1}
	Cyan        = Color{0, 1, 1, 1}
	Transparent = Color{0, 0, 0, 0}

	DarkGrey = Color{0.2, 0.2, 0.2, 1}
)

// Color holds 32-bit RGBA colors
type Color struct {
	R, G, B, A float32
}

// Color4 creates a color struct from its RGBA components
func Color4(r, g, b, a float32) Color {
	return Color{r, g, b, a}
}

// RGBA returns an 8-bit RGBA image/color
func (c Color) RGBA() color.RGBA {
	return color.RGBA{
		uint8(255.0 * c.R),
		uint8(255.0 * c.G),
		uint8(255.0 * c.B),
		uint8(255.0 * c.A),
	}
}

// Vec3 returns a vec3 containing the RGB components of the color
func (c Color) Vec3() vec3.T {
	return vec3.New(c.R, c.G, c.B)
}

// Vec4 returns a vec4 containing the RGBA components of the color
func (c Color) Vec4() vec4.T {
	return vec4.New(c.R, c.G, c.B, c.A)
}

func (c Color) Byte4() byte4.T {
	return byte4.New(
		byte(255.0*c.R),
		byte(255.0*c.G),
		byte(255.0*c.B),
		byte(255.0*c.A))
}

func (c Color) String() string {
	return fmt.Sprintf("(R:%.2f G:%.2f B:%.2f A:%.2f)", c.R, c.G, c.B, c.A)
}

// WithAlpha returns a new color with a modified alpha value
func (c Color) WithAlpha(a float32) Color {
	c.A = a
	return c
}

func Hex(s string) Color {
	if s[0] != '#' {
		panic("invalid color value")
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		panic("invalid color value")
	}

	c := Color{A: 1}
	switch len(s) {
	case 7:
		c.R = float32(hexToByte(s[1])<<4+hexToByte(s[2])) / 255
		c.G = float32(hexToByte(s[3])<<4+hexToByte(s[4])) / 255
		c.B = float32(hexToByte(s[5])<<4+hexToByte(s[6])) / 255
	case 4:
		c.R = float32(hexToByte(s[1])*17) / 255
		c.G = float32(hexToByte(s[2])*17) / 255
		c.B = float32(hexToByte(s[3])*17) / 255
	default:
		panic("invalid color value")
	}
	return c
}
