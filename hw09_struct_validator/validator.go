package hw09_struct_validator //nolint:golint,stylecheck
import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type validator interface {
	validate(field string, value interface{}) *ValidationError
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
