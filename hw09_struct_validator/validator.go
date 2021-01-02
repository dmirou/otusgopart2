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

type validator interface {
	validate(field string, value interface{}) *ValidationError
}

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

// validateField validate field value by the rules.
func validateField(field reflect.StructField, value reflect.Value, rules []string) error {
	var vr validator
	var err error

	for _, rule := range rules {
		parts := strings.Split(rule, ":")

		if vr, err = createValidator(parts[0], parts[1]); err != nil {
			return err
		}

		if field.Type.Kind() == reflect.Slice {
			for i := 0; i < value.Len(); i++ {
				if err := vr.validate(field.Name, value.Index(i).Interface()); err != nil {
					return err
				}
			}
			continue
		}

		if err := vr.validate(field.Name, value.Interface()); err != nil {
			return err
		}
	}

	return nil
}

// createValidator creates validator by name with specified params.
func createValidator(name, params string) (validator, error) {
	switch name {
	case "len":
		length, err := strconv.Atoi(params)
		if err != nil {
			return nil, err
		}

		return lengthValidator{
			Length: length,
		}, nil
	case "min":
		min, err := strconv.Atoi(params)
		if err != nil {
			return nil, err
		}

		return minValidator{
			Min: min,
		}, nil
	case "max":
		max, err := strconv.Atoi(params)
		if err != nil {
			return nil, err
		}

		return maxValidator{
			Max: max,
		}, nil
	case "in":
		return inValidator{
			Items: strings.Split(params, ","),
		}, nil
	case "regexp":
		return regexpValidator{
			Pattern: params,
		}, nil
	default:
	}

	return nil, ErrInvalidRule
}
