package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

type ValidationError struct {
	Field string
	Err   error
}

func CreateValidationError(field, rule string) ValidationError {
	err := fmt.Errorf("\"%s\" rule is not fulfilled", rule)
	return ValidationError{Field: field, Err: err}
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errorStrign := ""
	if len(v) > 0 {
		errorStringSlice := make([]string, len(v))
		for i, e := range v {
			errorStringSlice[i] = fmt.Sprintf("field is %s and error is %s", e.Field, e.Err)
		}
		errorStrign = strings.Join(errorStringSlice, "; ")
	}
	panic(errorStrign)
}

func Validate(v interface{}) error {
	err := validate(reflect.ValueOf(v))
	if err != nil {
		return err
	}
	return nil
}

func validate(reflectValue reflect.Value) error {
	var validationErrors ValidationErrors
	reflectType := reflectValue.Type()
	if reflectValue.Kind() == reflect.Pointer {
		reflectValue = reflectValue.Elem()
		reflectType = reflectType.Elem()
	}
	if reflectType.Kind() == reflect.Struct { //nolint:nestif
		fieldCount := reflectType.NumField()
		for i := range fieldCount {
			fieldType := reflectType.Field(i)
			fieldValue := reflectValue.Field(i)

			// проверка на публичность поля структуры
			if fieldType.PkgPath != "" {
				continue
			}
			tag := fieldType.Tag.Get("validate")
			// проверка на наличие правил "validate"
			if tag == "" || tag == "-" {
				continue
			}
			err := validateField(fieldType, fieldValue, tag)
			if err != nil {
				if errors.As(err, &ValidationErrors{}) {
					var validateErr ValidationErrors
					if errors.As(err, &validateErr) {
						validationErrors = append(validationErrors, validateErr...)
					}
					continue
				}
				return err
			}
		}
	}
	if len(validationErrors) == 0 {
		return nil
	}
	return validationErrors
}

func validateField(fieldType reflect.StructField, fieldValue reflect.Value, tagValue string) error {
	fieldName := fieldType.Name
	var err error
	switch fieldValue.Kind() { //nolint:exhaustive
	case reflect.String:
		err = validateString(fieldName, fieldValue.String(), tagValue)
	case reflect.Int:
		err = validateInteger(fieldName, int(fieldValue.Int()), tagValue)
	case reflect.Slice:
		sliceSize := fieldValue.Len()
		for i := range sliceSize {
			switch fieldType.Type.Elem().Kind() { //nolint:exhaustive
			case reflect.String:
				err = validateString(fieldName, fieldValue.Index(i).String(), tagValue)
			case reflect.Int:
				err = validateInteger(fieldName, int(fieldValue.Index(i).Int()), tagValue)
			}
			if err != nil {
				break
			}
		}
	case reflect.Struct:
		err = validate(fieldValue)
	}
	return err
}

func validateString(fieldName string, value string, tagValue string) error {
	var validationErrors ValidationErrors
	rules := strings.Split(tagValue, "|")
	for _, rule := range rules {
		ruleSplit := strings.Split(rule, ":")
		ruleName := ruleSplit[0]
		if len(ruleSplit) != 2 {
			return fmt.Errorf("the rule \"%s\" without value", ruleName)
		}
		ruleValue := ruleSplit[1]
		switch ruleName {
		case "len":
			ruleValueInt, err := strconv.Atoi(ruleValue)
			if err != nil {
				return fmt.Errorf("couldn't parse \"len\" rule value to int")
			}
			runeCount := utf8.RuneCountInString(value)
			if runeCount != ruleValueInt {
				validationErrors = append(validationErrors, CreateValidationError(fieldName, ruleName))
			}
		case "regexp":
			re := regexp.MustCompile(ruleValue)
			if !re.MatchString(value) {
				validationErrors = append(validationErrors, CreateValidationError(fieldName, ruleName))
			}
		case "in":
			ruleValueSlice := strings.Split(ruleValue, ",")
			if !slices.Contains(ruleValueSlice, value) {
				validationErrors = append(validationErrors, CreateValidationError(fieldName, ruleName))
			}
		}
	}
	if len(validationErrors) == 0 {
		return nil
	}
	return validationErrors
}

func validateInteger(fieldName string, value int, tagValue string) error {
	var validationErrors ValidationErrors
	rules := strings.Split(tagValue, "|")
	for _, rule := range rules {
		ruleSplit := strings.Split(rule, ":")
		ruleName := ruleSplit[0]
		if len(ruleSplit) != 2 {
			return fmt.Errorf("the rule \"%s\" without value", ruleName)
		}
		ruleValue := ruleSplit[1]
		switch ruleName {
		case "min":
			ruleValueInt, err := strconv.Atoi(ruleValue)
			if err != nil {
				return fmt.Errorf("couldn't parse \"min\" rule value to int")
			}
			if value < ruleValueInt {
				validationErrors = append(validationErrors, CreateValidationError(fieldName, ruleName))
			}
		case "max":
			ruleValueInt, err := strconv.Atoi(ruleValue)
			if err != nil {
				return fmt.Errorf("couldn't parse \"max\" rule value to int")
			}
			if value > ruleValueInt {
				validationErrors = append(validationErrors, CreateValidationError(fieldName, ruleName))
			}
		case "in":
			ruleValueSlice := strings.Split(ruleValue, ",")
			for _, v := range ruleValueSlice {
				ruleValueInt, err := strconv.Atoi(v)
				if err != nil {
					return fmt.Errorf("couldn't parse \"in\" rule value to int")
				}
				if value == ruleValueInt {
					break
				}
			}
			validationErrors = append(validationErrors, CreateValidationError(fieldName, ruleName))
		}
	}
	if len(validationErrors) == 0 {
		return nil
	}
	return validationErrors
}
