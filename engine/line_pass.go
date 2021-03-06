package engine

import "github.com/johanhenriksson/goworld/engine/object"

type LineDrawable interface {
	DrawLines(DrawArgs)
}

// LinePass draws line geometry
type LinePass struct {
}

// NewLinePass sets up a line geometry pass.
func NewLinePass() *LinePass {
	return &LinePass{}
}

func (p *LinePass) Resize(width, height int) {}

// DrawPass executes the line pass
func (p *LinePass) Draw(scene *Scene) {
	scene.Camera.Use()

	query := object.NewQuery(func(c object.Component) bool {
		_, ok := c.(LineDrawable)
		return ok
	})
	scene.Collect(&query)

	args := scene.Camera.DrawArgs()
	for _, component := range query.Results {
		drawable := component.(LineDrawable)
		drawable.DrawLines(args.Apply(component.Parent().Transform()))
	}
}
