package valley

// InitialPathSize sets the default size of a new Path's internal buffer.
var InitialPathSize = 32

// Path is used to represent the current position in a structure, to output a useful field value to
// identify where a ConstraintViolation occurred.
type Path struct {
	buf []byte
}

// NewPath returns a new Path instance.
func NewPath() *Path {
	return &Path{
		buf: make([]byte, 0, InitialPathSize),
	}
}

// Write ...
func (r *Path) Write(in string) int {
	r.buf = append(r.buf, in...)
	return len(in)
}

// TruncateRight ...
func (r *Path) TruncateRight(amount int) {
	r.buf = r.buf[:len(r.buf)-amount]
}

// Render renders this path as a string, to be sent to the frontend.
func (r *Path) Render() string {
	return Btos(r.buf)
}
