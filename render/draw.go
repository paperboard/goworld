package render

import (
	"github.com/johanhenriksson/goworld/math/mat4"
)

type DrawPass int32

const (
	GeometryPass DrawPass = iota
	LightPass
	LinePass
)

/** UI Component render interface */
type Drawable interface {
	Draw(DrawArgs)

	ZIndex() float32

	/* Render tree */
	Parent() Drawable
	SetParent(Drawable)
	Children() []Drawable
}

/** Passed to Drawables on render */
type DrawArgs struct {
	VP         mat4.T
	MVP        mat4.T
	Projection mat4.T
	View       mat4.T
	Transform  mat4.T
	Pass       DrawPass
}
