package valley

// Type ...
type Type struct{}

// Constraints ...
func (t Type) Constraints(_ ...Constraint) {}

// Field ...
func (t Type) Field(_ interface{}) Field {
	return Field{}
}

// Field ...
type Field struct{}

// Constraints ...
func (f Field) Constraints(_ ...Constraint) Field {
	return f
}

// Elements ...
func (f Field) Elements(_ ...Constraint) Field {
	return f
}
