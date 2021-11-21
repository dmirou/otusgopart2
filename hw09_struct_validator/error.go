package hw09_struct_validator //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidArg      = errors.New("invalid argument")
	ErrInvalidRule     = errors.New("invalid rule")
	ErrInvalidLength   = errors.New("value have invalid length")
	ErrValueLessMin    = errors.New("value should be greater or equal to min")
	ErrValueGreaterMax = errors.New("value should be less or equal to max")
	ErrValueNotIn      = errors.New("value is not allowed")
	ErrPatternNotMatch = errors.New("value pattern is invalid")
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("field: %s, error: %v", v.Field, v.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	b := strings.Builder{}

	for _, err := range v {
		b.WriteString(fmt.Sprintf("(%s)", err))
	}

	return b.String()
}
