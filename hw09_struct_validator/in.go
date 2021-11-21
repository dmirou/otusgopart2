package hw09_struct_validator //nolint:golint,stylecheck

import "fmt"

type inValidator struct {
	Items []string
}

func (iv inValidator) validate(field string, value interface{}) *ValidationError {
	str := fmt.Sprintf("%v", value)

	for _, item := range iv.Items {
		if str == item {
			return nil
		}
	}

	return &ValidationError{
		Field: field,
		Err:   ErrValueNotIn,
	}
}
