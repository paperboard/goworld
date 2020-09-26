package vec3

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
)

var (
	Zero  = T{0, 0, 0}
	One   = T{1, 1, 1}
	UnitX = T{1, 0, 0}
	UnitY = T{0, 1, 0}
	UnitZ = T{0, 0, 1}
)

// T holds a 3-component vector of 32-bit floats
type T struct {
	X, Y, Z float32
}

func (v T) Slice() [3]float32 {
	return [3]float32{v.X, v.Y, v.Z}
}

// Length returns the length of the vector.
// See also LengthSqr and Normalize.
func (v *T) Length() float32 {
	return math.Sqrt(v.LengthSqr())
}

// LengthSqr returns the squared length of the vector.
// See also Length and Normalize.
func (v *T) LengthSqr() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Absed returns a copy of the vector containing the absolute values.
func (v *T) Abs() T {
	return T{math.Abs(v.X), math.Abs(v.Y), math.Abs(v.Z)}
}

// Normalize normalizes the vector to unit length.
func (v *T) Normalize() {
	sl := v.LengthSqr()
	if sl == 0 || sl == 1 {
		return
	}
	s := 1 / math.Sqrt(sl)
	v.X *= s
	v.Y *= s
	v.Z *= s
}

// Normalized returns a unit length normalized copy of the vector.
func (v T) Normalized() T {
	v.Normalize()
	return v
}

// Scale the vector by a constant (in-place)
func (v *T) Scale(f float32) {
	v.X *= f
	v.Y *= f
	v.Z *= f
}

// Scaled returns a scaled vector
func (v T) Scaled(f float32) T {
	return T{v.X * f, v.Y * f, v.Z * f}
}

// Scaled returns a scaled vector
func (v T) ScaleI(i int) T {
	return v.Scaled(float32(i))
}

// Invert the vector components
func (v *T) Invert() {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
}

// Inverted returns an inverted vector
func (v *T) Inverted() T {
	i := *v
	i.Invert()
	return i
}

// Add each element of the vector with the corresponding element of another vector
func (v T) Add(v2 T) T {
	return T{
		v.X + v2.X,
		v.Y + v2.Y,
		v.Z + v2.Z,
	}
}

// Sub subtracts each element of the vector with the corresponding element of another vector
func (v T) Sub(v2 T) T {
	return T{
		v.X - v2.X,
		v.Y - v2.Y,
		v.Z - v2.Z,
	}
}

// Mul multiplies each element of the vector with the corresponding element of another vector
func (v T) Mul(v2 T) T {
	return T{
		v.X * v2.X,
		v.Y * v2.Y,
		v.Z * v2.Z,
	}
}

// XY returns a 2-component vector with the X, Y components of this vector
func (v T) XY() vec2.T {
	return vec2.T{X: v.X, Y: v.Y}
}

// XZ returns a 2-component vector with the X, Z components of this vector
func (v T) XZ() vec2.T {
	return vec2.T{X: v.X, Y: v.Z}
}

// YZ returns a 2-component vector with the Y, Z components of this vector
func (v T) YZ() vec2.T {
	return vec2.T{X: v.Y, Y: v.Z}
}

// Div divides each element of the vector with the corresponding element of another vector
func (v T) Div(v2 T) T {
	return T{v.X / v2.X, v.Y / v2.Y, v.Z / v2.Z}
}