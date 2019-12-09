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

// Write appends the given string to the end of the internal buffer.
func (r *Path) Write(in string) int {
	r.buf = append(r.buf, in...)
	return len(in)
}

// TruncateRight cuts n bytes off of the end of the buffer. The backing array for the buffer does
// not shrink, meaning we can re-use that memory if we need to.
func (r *Path) TruncateRight(n int) {
	r.buf = r.buf[:len(r.buf)-n]
}

// String renders this path as a string, to be sent to the frontend.
func (r *Path) String() string {
	return string(r.buf)
}
