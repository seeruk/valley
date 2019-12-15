package valley

// Type is the "fake" interface used to configure Valley for a Go type. It's methods are known by
// the config building process.
type Type struct{}

// Constraints accepts some constraints to generate code for.
func (t Type) Constraints(_ ...Constraint) Type {
	return t
}

// Field accepts a field to generate constraints for.
func (t Type) Field(_ interface{}) Field {
	return Field{}
}

// When accepts a predicate which will be used to wrap generated code to conditionally apply
// constraints.
func (t Type) When(_ bool) Type {
	return t
}

// Field represents the options for adding constraints to fields on a type.
type Field struct{}

// Constraints accepts some constraints to generate code for, for a specific field.
func (f Field) Constraints(_ ...Constraint) Field {
	return f
}

// Elements accepts some constraints to generate code for, on the elements of a specific field.
func (f Field) Elements(_ ...Constraint) Field {
	return f
}

// Keys accepts some constraints to generate code for, on the keys of a specific field.
func (f Field) Keys(_ ...Constraint) Field {
	return f
}
