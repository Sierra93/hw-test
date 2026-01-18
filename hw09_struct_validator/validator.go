package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidInputType = errors.New("input parameter type must be a struct")
	ErrInvalidInteger   = errors.New("invalid integer")
	ErrInvalidRegexp    = errors.New("invalig regexp")
	ErrInvalidFieldType = errors.New("invalid field type")
	ErrInvalidParam     = errors.New("invalid validate parameter")

	ErrValidationMin    = errors.New("value less than a required minimum")
	ErrValidationMax    = errors.New("value more than a required maximum")
	ErrValidationIn     = errors.New("value not found in list of possible values")
	ErrValidationLen    = errors.New("invalid value length")
	ErrValidationRegexp = errors.New("value does not match regexp")
)

type ValidationError struct {
	Field string
	Value any
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	buf := strings.Builder{}
	for _, err := range v {
		buf.WriteString(fmt.Sprintf("%s = %v: %s\n", err.Field, err.Value, err.Err))
	}
	return buf.String()
}

func Validate(v interface{}) error {
	validationErrors := make(ValidationErrors, 0)

	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return ErrInvalidInputType
	}
	typeOfValue := value.Type()

	for i := 0; i < typeOfValue.NumField(); i++ {
		field := typeOfValue.Field(i)
		tags, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}

		params := strings.Split(tags, "|")
		for _, param := range params {
			tag := strings.Split(param, ":")

			var err error
			switch tag[0] {
			case "min":
				err = validateMin(tag[1], field, value.Field(i))
			case "max":
				err = validateMax(tag[1], field, value.Field(i))
			case "len":
				err = validateLen(tag[1], field, value.Field(i))
			case "in":
				err = validateIn(tag[1], field, value.Field(i))
			case "regexp":
				err = validateRegexp(tag[1], field, value.Field(i))
			default:
				return ErrInvalidParam
			}

			var verr ValidationErrors
			if errors.As(err, &verr) {
				validationErrors = append(validationErrors, verr...)
			} else {
				return err
			}
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateMin(paramValue string, field reflect.StructField, value reflect.Value) error {
	validationErrors := make(ValidationErrors, 0)

	minVal, err := strconv.Atoi(paramValue)
	if err != nil {
		return ErrInvalidInteger
	}

	switch field.Type.Kind() { //nolint: exhaustive
	case reflect.Int:
		if err := checkMin(int(value.Int()), minVal); err != nil {
			validationErrors = append(validationErrors,
				ValidationError{
					Field: field.Name,
					Value: value.Int(),
					Err:   err,
				})
		}

	case reflect.Slice:
		switch fieldValues := value.Interface().(type) {
		case []int:
			for _, fv := range fieldValues {
				if err := checkMin(fv, minVal); err != nil {
					validationErrors = append(validationErrors,
						ValidationError{
							Field: field.Name,
							Value: fv,
							Err:   err,
						})
				}
			}
		default:
			return ErrInvalidFieldType
		}
	default:
		return ErrInvalidFieldType
	}

	return validationErrors
}

func checkMin(fieldValue, minVal int) error {
	if fieldValue < minVal {
		return ErrValidationMin
	}
	return nil
}

func validateMax(paramValue string, field reflect.StructField, value reflect.Value) error {
	validationErrors := make(ValidationErrors, 0)

	maxVal, err := strconv.Atoi(paramValue)
	if err != nil {
		return ErrInvalidInteger
	}

	switch field.Type.Kind() { //nolint: exhaustive
	case reflect.Int:
		if err := checkMax(int(value.Int()), maxVal); err != nil {
			validationErrors = append(validationErrors,
				ValidationError{
					Field: field.Name,
					Value: value.Int(),
					Err:   err,
				})
		}

	case reflect.Slice:
		switch fieldValues := value.Interface().(type) {
		case []int:
			for _, fv := range fieldValues {
				if err := checkMax(fv, maxVal); err != nil {
					validationErrors = append(validationErrors,
						ValidationError{
							Field: field.Name,
							Value: fv,
							Err:   err,
						})
				}
			}
		default:
			return ErrInvalidFieldType
		}
	default:
		return ErrInvalidFieldType
	}

	return validationErrors
}

func checkMax(fieldValue, maxVal int) error {
	if fieldValue > maxVal {
		return ErrValidationMax
	}
	return nil
}

func validateLen(paramValue string, field reflect.StructField, value reflect.Value) error {
	validationErrors := make(ValidationErrors, 0)

	length, err := strconv.Atoi(paramValue)
	if err != nil {
		return ErrInvalidInteger
	}

	switch field.Type.Kind() { //nolint: exhaustive
	case reflect.String:
		if err := checkLen(value.String(), length); err != nil {
			validationErrors = append(validationErrors,
				ValidationError{
					Field: field.Name,
					Value: value.String(),
					Err:   err,
				})
		}

	case reflect.Slice:
		switch fieldValues := value.Interface().(type) {
		case []string:
			for _, fv := range fieldValues {
				if err := checkLen(fv, length); err != nil {
					validationErrors = append(validationErrors,
						ValidationError{
							Field: field.Name,
							Value: fv,
							Err:   err,
						})
				}
			}
		default:
			return ErrInvalidFieldType
		}
	default:
		return ErrInvalidFieldType
	}

	return validationErrors
}

func checkLen(fieldValue string, length int) error {
	if len(fieldValue) != length {
		return ErrValidationLen
	}
	return nil
}

func validateIn(paramValue string, field reflect.StructField, value reflect.Value) error {
	validationErrors := make(ValidationErrors, 0)
	possibleValuesStr := strings.Split(paramValue, ",")

	switch field.Type.Kind() { //nolint: exhaustive
	case reflect.Int:
		possibleValuesInt, err := parseValues(possibleValuesStr)
		if err != nil {
			return err
		}

		if err := checkIn(int(value.Int()), possibleValuesInt); err != nil {
			validationErrors = append(validationErrors,
				ValidationError{
					Field: field.Name,
					Value: value.Int(),
					Err:   err,
				})
		}

	case reflect.String:
		if err := checkIn(value.String(), possibleValuesStr); err != nil {
			validationErrors = append(validationErrors,
				ValidationError{
					Field: field.Name,
					Value: value.String(),
					Err:   err,
				})
		}

	case reflect.Slice:
		switch fieldValues := value.Interface().(type) {
		case []int:
			possibleValuesInt, err := parseValues(possibleValuesStr)
			if err != nil {
				return err
			}

			for _, fv := range fieldValues {
				if err := checkIn(fv, possibleValuesInt); err != nil {
					validationErrors = append(validationErrors,
						ValidationError{
							Field: field.Name,
							Value: fv,
							Err:   err,
						})
				}
			}
		case []string:
			for _, fv := range fieldValues {
				if err := checkIn(fv, possibleValuesStr); err != nil {
					validationErrors = append(validationErrors,
						ValidationError{
							Field: field.Name,
							Value: fv,
							Err:   err,
						})
				}
			}
		default:
			return ErrInvalidFieldType
		}
	default:
		return ErrInvalidFieldType
	}

	return validationErrors
}

func parseValues(values []string) ([]int, error) {
	valuesInt := make([]int, 0, len(values))
	for _, value := range values {
		value, err := strconv.Atoi(value)
		if err != nil {
			return nil, ErrInvalidInteger
		}
		valuesInt = append(valuesInt, value)
	}
	return valuesInt, nil
}

func checkIn[T comparable](fieldValue T, paramValues []T) error {
	for _, pv := range paramValues {
		if fieldValue == pv {
			return nil
		}
	}
	return ErrValidationIn
}

func validateRegexp(paramValue string, field reflect.StructField, value reflect.Value) error {
	validationErrors := make(ValidationErrors, 0)

	re, err := regexp.Compile(paramValue)
	if err != nil {
		return ErrInvalidRegexp
	}

	switch field.Type.Kind() { //nolint: exhaustive
	case reflect.String:
		if err := checkRegexp(value.String(), re); err != nil {
			validationErrors = append(validationErrors,
				ValidationError{
					Field: field.Name,
					Value: value.String(),
					Err:   err,
				})
		}

	case reflect.Slice:
		if fieldValues, ok := value.Interface().([]string); ok {
			for _, fv := range fieldValues {
				if err := checkRegexp(fv, re); err != nil {
					validationErrors = append(validationErrors,
						ValidationError{
							Field: field.Name,
							Value: fv,
							Err:   err,
						})
				}
			}
		}
	default:
		return ErrInvalidFieldType
	}
	return validationErrors
}

func checkRegexp(fieldValue string, re *regexp.Regexp) error {
	if !re.MatchString(fieldValue) {
		return ErrValidationRegexp
	}
	return nil
}
