// Code generated by valley. DO NOT EDIT.
package main

import "strconv"
import "github.com/seeruk/valley/valley"

// Validate validates this Example.
// This method was generated by Valley.
func (e Example) Validate(path *valley.Path) []valley.ConstraintViolation {
	var violations []valley.ConstraintViolation

	path.Write(".")

	{
		// MutuallyExclusive uses it's own block to lock down nonEmpty's scope.
		var nonEmpty []string

		if !(len(e.Text) == 0) {
			nonEmpty = append(nonEmpty, "Text")
		}

		if !(len(e.Texts) == 0) {
			nonEmpty = append(nonEmpty, "Texts")
		}

		if !(len(e.TextMap) == 0) {
			nonEmpty = append(nonEmpty, "TextMap")
		}

		if len(nonEmpty) > 1 {

			violations = append(violations, valley.ConstraintViolation{
				Field:   path.String(),
				Message: "fields are mutually exclusive",
				Details: map[string]interface{}{
					"fields": nonEmpty,
				},
			})

		}
	}

	if e.Int == 0 {
		size := path.Write("Int")
		violations = append(violations, valley.ConstraintViolation{
			Field:   path.String(),
			Message: "a value is required",
		})
		path.TruncateRight(size)
	}

	if len(e.Ints) == 0 {
		size := path.Write("Ints")
		violations = append(violations, valley.ConstraintViolation{
			Field:   path.String(),
			Message: "a value is required",
		})
		path.TruncateRight(size)
	}

	for i, element := range e.Ints {

		if element == 0 {
			size := path.Write("Ints.[" + strconv.Itoa(i) + "]")
			violations = append(violations, valley.ConstraintViolation{
				Field:   path.String(),
				Message: "a value is required",
			})
			path.TruncateRight(size)
		}

	}

	if e.Nested == nil {
		size := path.Write("Nested")
		violations = append(violations, valley.ConstraintViolation{
			Field:   path.String(),
			Message: "a value is required",
		})
		path.TruncateRight(size)
	}

	if e.Nested != nil {
		size := path.Write("Nested")
		violations = append(violations, e.Nested.Validate(path)...)
		path.TruncateRight(size)
	}

	if len(e.Text) == 0 {
		size := path.Write("Text")
		violations = append(violations, valley.ConstraintViolation{
			Field:   path.String(),
			Message: "a value is required",
		})
		path.TruncateRight(size)
	}

	return violations
}

// Validate validates this NestedExample.
// This method was generated by Valley.
func (n NestedExample) Validate(path *valley.Path) []valley.ConstraintViolation {
	var violations []valley.ConstraintViolation

	path.Write(".")

	if len(n.Text) == 0 {
		size := path.Write("Text")
		violations = append(violations, valley.ConstraintViolation{
			Field:   path.String(),
			Message: "a value is required",
		})
		path.TruncateRight(size)
	}

	return violations
}
