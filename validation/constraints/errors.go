package constraints

import "errors"

var (
	// ErrTypeWarning is an error returned when a constraint might have been used on an unsupported
	// type. Some constraints may choose to be permissive and continue anyway. This error will only
	// result in a warning being printed. It's up to the constraint if it halts constraint
	// generation (i.e. if that constraint will produce no code, but execution will continue...)
	ErrTypeWarning = errors.New("type used may not produce valid code (is it a custom type?)")
)
