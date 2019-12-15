package main

import (
	"math"
	"testing"
	"time"

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

	example.Bool = true
	example.Text = "Hello"
	//example.Texts = []string{"Hello", "World!"}
	example.TextMap = map[string]string{"hello longer key": "world"}
	example.Int = 999
	example.Int2 = &example.Int
	example.Ints = []int{1}
	example.Float = math.Pi
	example.Nested = &NestedExample{Text: "Hello, World!"}
	example.Adults = 2
	example.Children = 4
	example.Times = []time.Time{
		time.Date(1800, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		violations = example.Validate(valley.NewPath())
	}

	if len(violations) > 0 {
		b.Error("expected no violations")
		b.Logf("%+v\n", violations)
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
		b.Error("expected violations")
	}
}
