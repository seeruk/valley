package main

import (
	"testing"

	"github.com/seeruk/valley"
)

func BenchmarkRequired(b *testing.B) {
	violations := make([]valley.ConstraintViolation, 1)

	ints := []int{}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if len(ints) == 0 {
			violations[0] = valley.ConstraintViolation{
				Path:    "test",
				Message: "a value is required",
			}
		}
	}

	_ = violations
}

func BenchmarkExample_ValidateHappy(b *testing.B) {
	var example Example
	var violations []valley.ConstraintViolation

	example.Text = "Hello"
	//example.Texts = []string{"Hello", "World!"}
	example.TextMap = map[string]string{"hello": "world"}
	example.Int = 999
	example.Ints = []int{1}
	example.Nested = &NestedExample{Text: "Hello, World!"}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		violations = example.Validate(valley.NewPath())
	}

	if len(violations) > 0 {
		b.Error("expected no violations")
		b.Log(violations)
	}
}

func BenchmarkExample_ValidateUnhappy(b *testing.B) {
	var example Example
	var violations []valley.ConstraintViolation

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		violations = example.Validate(valley.NewPath())
	}

	if len(violations) == 0 {
		b.Error("expected no violations")
	}
}
