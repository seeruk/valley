// Code generated by valley. DO NOT EDIT.
package main

import fmt "fmt"
import valley "github.com/seeruk/valley"
import math "math"
import reflect "reflect"
import regexp "regexp"
import strconv "strconv"
import strings "strings"
import time "time"

// Reference imports to suppress errors if they aren't otherwise used
var _ = fmt.Sprintf
var _ = strconv.Itoa

// Variables generated by constraints:
var github_com_seeruk_valley_validation_constraints_RegexpString_Example_25 = regexp.MustCompile("^Hello")

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

		if len(nonEmpty) > 1 {

			violations = append(violations, valley.ConstraintViolation{
				Path:     path.String(),
				PathKind: "struct",
				Message:  "fields are mutually exclusive",
				Details: map[string]interface{}{
					"fields": nonEmpty,
				},
			})

		}
	}

	{
		// MutuallyInclusive uses it's own block to lock down nonEmpty's scope.
		var nonEmpty []string

		if !(e.Int == 0) {
			nonEmpty = append(nonEmpty, "Int")
		}

		if !(e.Int2 == nil) {
			nonEmpty = append(nonEmpty, "Int2")
		}

		if !(len(e.Ints) == 0) {
			nonEmpty = append(nonEmpty, "Ints")
		}

		if len(nonEmpty) > 0 && len(nonEmpty) != 3 {

			violations = append(violations, valley.ConstraintViolation{
				Path:     path.String(),
				PathKind: "struct",
				Message:  "fields are mutually inclusive",
				Details: map[string]interface{}{
					"fields": nonEmpty,
				},
			})

		}
	}

	if e.Adults < 1 {
		size := path.Write("Adults")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "minimum value not met",
			Details: map[string]interface{}{
				"minimum": 1,
			},
		})
		path.TruncateRight(size)
	}

	if e.Adults > 9 {
		size := path.Write("Adults")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "maximum value exceeded",
			Details: map[string]interface{}{
				"maximum": 9,
			},
		})
		path.TruncateRight(size)
	}

	if e.Bool == false {
		size := path.Write("Bool")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "values must not be equal",
			Details: map[string]interface{}{
				"equal_to": false,
			},
		})
		path.TruncateRight(size)
	}

	if !reflect.DeepEqual(e.Bool, true) {
		size := path.Write("Bool")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "values must be deeply equal",
			Details: map[string]interface{}{
				"deeply_equal_to": true,
			},
		})
		path.TruncateRight(size)
	}

	if len(e.Chan) > 12 {
		size := path.Write("Chan")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "maximum length exceeded",
			Details: map[string]interface{}{
				"maximum": 12,
			},
		})
		path.TruncateRight(size)
	}

	if e.Children < 0 {
		size := path.Write("Children")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "minimum value not met",
			Details: map[string]interface{}{
				"minimum": 0,
			},
		})
		path.TruncateRight(size)
	}

	if e.Children != e.Adults+2 {
		size := path.Write("Children")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "values must be equal",
			Details: map[string]interface{}{
				"equal_to": e.Adults + 2,
			},
		})
		path.TruncateRight(size)
	}

	if e.Children > int(math.Max(float64(8-(e.Adults-1)), 0)) {
		size := path.Write("Children")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "maximum value exceeded",
			Details: map[string]interface{}{
				"maximum": int(math.Max(float64(8-(e.Adults-1)), 0)),
			},
		})
		path.TruncateRight(size)
	}

	if e.Float != math.Pi {
		size := path.Write("Float")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "values must be equal",
			Details: map[string]interface{}{
				"equal_to": math.Pi,
			},
		})
		path.TruncateRight(size)
	}

	if e.Int == 0 {
		size := path.Write("Int")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "a value is required",
		})
		path.TruncateRight(size)
	}

	if e.Int2 == nil {
		size := path.Write("Int2")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "a value is required",
		})
		path.TruncateRight(size)
	}

	if e.Int2 == nil {
		size := path.Write("Int2")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "value must not be nil",
		})
		path.TruncateRight(size)
	}

	if e.Int2 != nil && *e.Int2 < 0 {
		size := path.Write("Int2")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "minimum value not met",
			Details: map[string]interface{}{
				"minimum": 0,
			},
		})
		path.TruncateRight(size)
	}

	if len(e.Ints) == 0 {
		size := path.Write("Ints")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "a value is required",
		})
		path.TruncateRight(size)
	}

	if len(e.Ints) > 3 {
		size := path.Write("Ints")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "maximum length exceeded",
			Details: map[string]interface{}{
				"maximum": 3,
			},
		})
		path.TruncateRight(size)
	}

	for i, element := range e.Ints {

		if element == 0 {
			size := path.Write("Ints.[" + strconv.Itoa(i) + "]")
			violations = append(violations, valley.ConstraintViolation{
				Path:     path.String(),
				PathKind: "element",
				Message:  "a value is required",
			})
			path.TruncateRight(size)
		}

		if element < 0 {
			size := path.Write("Ints.[" + strconv.Itoa(i) + "]")
			violations = append(violations, valley.ConstraintViolation{
				Path:     path.String(),
				PathKind: "element",
				Message:  "minimum value not met",
				Details: map[string]interface{}{
					"minimum": 0,
				},
			})
			path.TruncateRight(size)
		}

	}

	if e.Nested == nil {
		size := path.Write("Nested")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "a value is required",
		})
		path.TruncateRight(size)
	}

	if e.Nested != nil {
		size := path.Write("Nested")
		violations = append(violations, e.Nested.Validate(path)...)
		path.TruncateRight(size)
	}

	for key := range e.NestedMap {
		size := path.Write("NestedMap.[" + fmt.Sprintf("%v", key) + "]")
		violations = append(violations, key.Validate(path)...)
		path.TruncateRight(size)

	}

	for i, element := range e.Nesteds {
		if element != nil {
			size := path.Write("Nesteds.[" + strconv.Itoa(i) + "]")
			violations = append(violations, element.Validate(path)...)
			path.TruncateRight(size)
		}

	}

	if len(e.Text) == 0 {
		size := path.Write("Text")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "a value is required",
		})
		path.TruncateRight(size)
	}

	if !github_com_seeruk_valley_validation_constraints_RegexpString_Example_25.MatchString(e.Text) {
		size := path.Write("Text")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "value must match regular expression",
			Details: map[string]interface{}{
				"regexp": github_com_seeruk_valley_validation_constraints_RegexpString_Example_25.String(),
			},
		})
		path.TruncateRight(size)
	}

	if strings.HasPrefix(e.Text, "custom") && len(e.Text) == 32 {
		size := path.Write("Text")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "\"value must be a valid custom ID\"",
		})
		path.TruncateRight(size)
	}

	if len(e.Text) > 12 {
		size := path.Write("Text")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "maximum length exceeded",
			Details: map[string]interface{}{
				"maximum": 12,
			},
		})
		path.TruncateRight(size)
	}

	if len(e.Text) != 5 {
		size := path.Write("Text")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "exact length not met",
			Details: map[string]interface{}{
				"exactly": 5,
			},
		})
		path.TruncateRight(size)
	}

	if len(e.TextMap) == 0 {
		size := path.Write("TextMap")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "a value is required",
		})
		path.TruncateRight(size)
	}

	for i, element := range e.TextMap {

		if len(element) == 0 {
			size := path.Write("TextMap.[" + fmt.Sprintf("%v", i) + "]")
			violations = append(violations, valley.ConstraintViolation{
				Path:     path.String(),
				PathKind: "element",
				Message:  "a value is required",
			})
			path.TruncateRight(size)
		}

	}

	for key := range e.TextMap {

		if len(key) < 10 {
			size := path.Write("TextMap.[" + fmt.Sprintf("%v", key) + "]")
			violations = append(violations, valley.ConstraintViolation{
				Path:     path.String(),
				PathKind: "key",
				Message:  "minimum length not met",
				Details: map[string]interface{}{
					"minimum": 10,
				},
			})
			path.TruncateRight(size)
		}

	}

	if !e.Time.Before(timeYosemite) {
		size := path.Write("Time")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "value must be before time",
			Details: map[string]interface{}{
				"time": timeYosemite.Format(time.RFC3339),
			},
		})
		path.TruncateRight(size)
	}

	if len(e.Times) < 1 {
		size := path.Write("Times")
		violations = append(violations, valley.ConstraintViolation{
			Path:     path.String(),
			PathKind: "field",
			Message:  "minimum length not met",
			Details: map[string]interface{}{
				"minimum": 1,
			},
		})
		path.TruncateRight(size)
	}

	for i, element := range e.Times {

		if !element.Before(timeYosemite) {
			size := path.Write("Times.[" + strconv.Itoa(i) + "]")
			violations = append(violations, valley.ConstraintViolation{
				Path:     path.String(),
				PathKind: "element",
				Message:  "value must be before time",
				Details: map[string]interface{}{
					"time": timeYosemite.Format(time.RFC3339),
				},
			})
			path.TruncateRight(size)
		}

	}

	path.TruncateRight(1)

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
			Path:     path.String(),
			PathKind: "field",
			Message:  "a value is required",
		})
		path.TruncateRight(size)
	}

	path.TruncateRight(1)

	return violations
}
