package object

import "github.com/johanhenriksson/goworld/math/vec3"

// Builder API for game objects
type Builder struct {
	name     string
	position vec3.T
	rotation vec3.T
	scale    vec3.T
	active   bool
	children []Component
}

// Build instantiates a new object builder.
func Build(name string) *Builder {
	return &Builder{
		name:     name,
		position: vec3.Zero,
		rotation: vec3.Zero,
		scale:    vec3.One,
		active:   true,
	}
}

// Attach a component to the object.
func (b *Builder) Attach(c Component) *Builder {
	b.children = append(b.children, c)
	return b
}

// Position sets the intial position of the object.
func (b *Builder) Position(p vec3.T) *Builder {
	b.position = p
	return b
}

// Rotation sets the intial rotation of the object.
func (b *Builder) Rotation(r vec3.T) *Builder {
	b.rotation = r
	return b
}

// Scale sets the intial scale of the object.
func (b *Builder) Scale(s vec3.T) *Builder {
	b.scale = s
	return b
}

// Active sets the objects active flag.
func (b *Builder) Active(active bool) *Builder {
	b.active = active
	return b
}

// Create instantiates a new object with the current builder settings.
func (b *Builder) Create(parent *T) *T {
	obj := New(b.name, b.children...)
	obj.SetPosition(b.position)
	obj.SetRotation(b.rotation)
	obj.SetScale(b.scale)
	obj.SetActive(b.active)
	if parent != nil {
		parent.Attach(obj)
	}
	return obj
}
