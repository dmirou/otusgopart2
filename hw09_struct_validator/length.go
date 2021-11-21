package hw09_struct_validator

import "fmt"

type lengthValidator struct {
	Length int
}

func (lv lengthValidator) validate(field string, value interface{}) *ValidationError {
	str := fmt.Sprintf("%s", value)

	if len(str) != lv.Length {
		return &ValidationError{
			Field: field,
			Err:   ErrInvalidLength,
		}
	}

	return nil
}
