package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"regexp"
)

type regexpValidator struct {
	Pattern string
}

func (rv regexpValidator) validate(field string, value interface{}) *ValidationError {
	str := fmt.Sprintf("%s", value)

	re := regexp.MustCompile(rv.Pattern)

	if !re.MatchString(str) {
		return &ValidationError{
			Field: field,
			Err:   ErrPatternNotMatch,
		}
	}

	return nil
}
