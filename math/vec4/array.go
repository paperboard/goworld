package vec4

import "unsafe"

// Array holds an array of 4-component vectors
type Array []T

func (a Array) Elements() int {
	return len(a)
}

func (a Array) Size() int {
	return 16
}

func (a Array) Pointer() unsafe.Pointer {
	return unsafe.Pointer(&a[0])
}
