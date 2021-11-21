package hw09_struct_validator

type maxValidator struct {
	Max int
}

func (mv maxValidator) validate(field string, value interface{}) *ValidationError {
	vint, ok := value.(int)
	if !ok {
		return &ValidationError{
			Field: field,
			Err:   ErrInvalidArg,
		}
	}

	if vint > mv.Max {
		return &ValidationError{
			Field: field,
			Err:   ErrValueGreaterMax,
		}
	}

	return nil
}
