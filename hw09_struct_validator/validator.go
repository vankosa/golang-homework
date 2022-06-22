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
	ErrLength = errors.New("invalid length")
	ErrRegexp = errors.New("regexp doesn't match")
	ErrNotIn  = errors.New("value isn't in")
	ErrMin    = errors.New("value < min")
	ErrMax    = errors.New("value > max")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	result := ""
	for _, e := range v {
		result = result + e.Error() + "; "
	}
	return result
}

func (v ValidationError) Error() string {
	return v.Field + ":" + v.Err.Error()
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return errors.New("not a struct")
	}

	t := reflect.TypeOf(v)

	finalErrors := ValidationErrors{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		fieldValue := val.Field(i)

		tagValue := field.Tag.Get("validate")

		success, errs := validateField(field.Name, fieldValue, tagValue)
		if !success {
			return errs[0]
		}

		if len(errs) == 0 {
			continue
		}

		for _, e := range errs {
			validateError, ok := e.(ValidationError) //nolint:errorlint
			if !ok {
				return errors.New("unexpected error")
			}

			finalErrors = append(finalErrors, validateError)
		}
	}
	return finalErrors
}

func validateField(fieldName string, v reflect.Value, tagValue string) (bool, []error) {
	result := []error{}
	switch v.Kind() { //nolint:exhaustive
	case reflect.String:
		tags := strings.Split(tagValue, "|")
		for _, tag := range tags {
			success, err := validateString(tag, fieldName, v.String())
			if success {
				if err != nil {
					result = append(result, err)
				}
			} else {
				return false, []error{err}
			}
		}
	case reflect.Int:
		tags := strings.Split(tagValue, "|")
		for _, tag := range tags {
			success, err := validateInt(tag, fieldName, int(v.Int()))
			if success {
				if err != nil {
					result = append(result, err)
				}
			} else {
				return false, []error{err}
			}
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			success, errs := validateField(fieldName, v.Index(i), tagValue)
			if success {
				result = append(result, errs...)
			} else {
				return false, errs
			}
		}
	}

	return true, result
}

var patternRule = regexp.MustCompile(`(\w+):(.+)`)

func validateString(tag string, name, value string) (bool, error) {
	m := patternRule.FindStringSubmatch(tag)
	if m == nil {
		return false, errors.New("invalid tag")
	}
	rule := m[1]
	ruleValue := m[2]

	switch rule {
	case "len":
		length, err := strconv.Atoi(ruleValue)
		if err != nil {
			return false, errors.New("invalid length value")
		}
		if len(value) != length {
			return true, ValidationError{Field: name, Err: fmt.Errorf("%s: %w", tag, ErrLength)}
		}
	case "regexp":
		re := regexp.MustCompile(ruleValue)
		if !re.Match([]byte(value)) {
			return true, ValidationError{Field: name, Err: fmt.Errorf("%s: %w", tag, ErrRegexp)}
		}
	case "in":
		variants := strings.Split(ruleValue, ",")

		exist := false
		for _, variant := range variants {
			if variant == value {
				exist = true
			}
		}

		if !exist {
			return true, ValidationError{Field: name, Err: fmt.Errorf("%s: %w", tag, ErrNotIn)}
		}
	}

	return true, nil
}

func validateInt(tag string, name string, value int) (bool, error) {
	m := patternRule.FindStringSubmatch(tag)
	if m == nil {
		return false, errors.New("invalid tag")
	}
	rule := m[1]
	ruleValue := m[2]

	switch rule {
	case "min":
		i, err := strconv.Atoi(ruleValue)
		if err != nil {
			return false, errors.New("invalid min value")
		}
		if value < i {
			return true, ValidationError{Field: name, Err: fmt.Errorf("%s: %w", tag, ErrMin)}
		}
	case "max":
		i, err := strconv.Atoi(ruleValue)
		if err != nil {
			return false, errors.New("invalid max value")
		}
		if value > i {
			return true, ValidationError{Field: name, Err: fmt.Errorf("%s: %w", tag, ErrMax)}
		}
	case "in":
		variants := strings.Split(ruleValue, ",")

		exist := false
		for _, variant := range variants {
			if variant == strconv.Itoa(value) {
				exist = true
			}
		}

		if !exist {
			return true, ValidationError{Field: name, Err: fmt.Errorf("%s: %w", tag, ErrNotIn)}
		}
	}

	return true, nil
}
