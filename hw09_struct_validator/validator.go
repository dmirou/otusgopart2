package hw09_struct_validator //nolint:golint,stylecheck
import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidArg      = errors.New("invalid argument")
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

func Validate(v interface{}) error {
	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)

	if rt.Kind() != reflect.Struct {
		return ErrInvalidArg
	}

	var errs ValidationErrors

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i)

		tag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}

		rules := strings.Split(tag, "|")

		err := validateField(field, value, rules)
		if err != nil {
			var e *ValidationError

			if errors.As(err, &e) {
				errs = append(errs, *e)
				continue
			}

			return err
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}

// validateField validate field value by the rules
func validateField(field reflect.StructField, value reflect.Value, rules []string) error {
	for _, rule := range rules {
		parts := strings.Split(rule, ":")

		switch parts[0] {
		case "len":
			length, err := strconv.Atoi(parts[1])
			if err != nil {
				return err
			}

			if field.Type.Kind() == reflect.Slice {
				for i := 0; i < value.Len(); i++ {
					if err := validateLength(field.Name, value.Index(i).Interface(), length); err != nil {
						return err
					}
				}
				continue
			}

			if err := validateLength(field.Name, value.Interface(), length); err != nil {
				return err
			}
		case "min":
			min, err := strconv.Atoi(parts[1])
			if err != nil {
				return err
			}

			if field.Type.Kind() == reflect.Slice {
				for i := 0; i < value.Len(); i++ {
					if err := validateMin(field.Name, value.Index(i).Interface(), min); err != nil {
						return err
					}
				}
				continue
			}

			if err := validateMin(field.Name, value.Interface(), min); err != nil {
				return err
			}
		case "max":
			max, err := strconv.Atoi(parts[1])
			if err != nil {
				return err
			}

			if field.Type.Kind() == reflect.Slice {
				for i := 0; i < value.Len(); i++ {
					if err := validateMax(field.Name, value.Index(i).Interface(), max); err != nil {
						return err
					}
				}
				continue
			}

			if err := validateMax(field.Name, value.Interface(), max); err != nil {
				return err
			}
		case "in":
			items := strings.Split(parts[1], ",")

			if field.Type.Kind() == reflect.Slice {
				for i := 0; i < value.Len(); i++ {
					if err := validateIn(field.Name, value.Index(i).Interface(), items); err != nil {
						return err
					}
				}
				continue
			}

			if err := validateIn(field.Name, value.Interface(), items); err != nil {
				return err
			}
		case "regexp":
			pattern := parts[1]

			if field.Type.Kind() == reflect.Slice {
				for i := 0; i < value.Len(); i++ {
					if err := validateByPattern(field.Name, value.Index(i).Interface(), pattern); err != nil {
						return err
					}
				}
				continue
			}

			if err := validateByPattern(field.Name, value.Interface(), pattern); err != nil {
				return err
			}
		default:
		}
	}

	return nil
}

func validateLength(field string, value interface{}, length int) *ValidationError {
	str := fmt.Sprintf("%s", value)

	if len(str) != length {
		return &ValidationError{
			Field: field,
			Err:   ErrInvalidLength,
		}
	}

	return nil
}

func validateMin(field string, value interface{}, min int) *ValidationError {
	vint, ok := value.(int)
	if !ok {
		return &ValidationError{
			Field: field,
			Err:   ErrInvalidArg,
		}
	}

	if vint < min {
		return &ValidationError{
			Field: field,
			Err:   ErrValueLessMin,
		}
	}

	return nil
}

func validateMax(field string, value interface{}, max int) *ValidationError {
	vint, ok := value.(int)
	if !ok {
		return &ValidationError{
			Field: field,
			Err:   ErrInvalidArg,
		}
	}

	if vint > max {
		return &ValidationError{
			Field: field,
			Err:   ErrValueGreaterMax,
		}
	}

	return nil
}

func validateIn(field string, value interface{}, items []string) *ValidationError {
	str := fmt.Sprintf("%v", value)

	for _, item := range items {
		if str == item {
			return nil
		}
	}

	return &ValidationError{
		Field: field,
		Err:   ErrValueNotIn,
	}
}

func validateByPattern(field string, value interface{}, pattern string) *ValidationError {
	str := fmt.Sprintf("%s", value)

	re := regexp.MustCompile(pattern)

	if !re.MatchString(str) {
		return &ValidationError{
			Field: field,
			Err:   ErrPatternNotMatch,
		}
	}

	return nil
}
