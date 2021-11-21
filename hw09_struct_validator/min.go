package hw09_struct_validator

type minValidator struct {
	Min int
}

func (mv minValidator) validate(field string, value interface{}) *ValidationError {
	vint, ok := value.(int)
	if !ok {
		return &ValidationError{
			Field: field,
			Err:   ErrInvalidArg,
		}
	}

	if vint < mv.Min {
		return &ValidationError{
			Field: field,
			Err:   ErrValueLessMin,
		}
	}

	return nil
}
