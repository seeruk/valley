package valley

import "strings"

// Path is used to represent the current position in a structure, to output a useful field value to
// identify where a ConstraintViolation occurred.
type Path struct {
	items []string
}

// NewPath returns a new Path instance.
func NewPath() *Path {
	return &Path{}
}

// Push adds an item to the end of the path.
func (r *Path) Push(item string) {
	r.items = append(r.items, item)
}

// Pop removes an item from the end of the path, and returns it.
func (r *Path) Pop() string {
	var p string
	p, r.items = r.items[len(r.items)-1], r.items[:len(r.items)-1]
	return p
}

// Render renders this path as a string, to be sent to the frontend.
func (r *Path) Render() string {
	return strings.Join(r.items, ".")
}
