package valleydemo

import "github.com/seeruk/valley/valley"

func (e Example) Validate() []valley.ConstraintViolation {
	var violations []valley.ConstraintViolation

	if len(e.Text) == 0 {
		violations = append(violations, valley.ConstraintViolation{
			// TODO: Calculate Field value properly.
			// Field: path.Render(),
			Field: "Text",
			Message: "a value is required",
		})
	}


	if e.Texts == nil || len(e.Texts) == 0 {
		violations = append(violations, valley.ConstraintViolation{
			// TODO: Calculate Field value properly.
			// Field: path.Render(),
			Field: "Texts",
			Message: "a value is required",
		})
	}

	return violations
}
