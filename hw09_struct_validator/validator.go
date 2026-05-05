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
	panic("implement me")
}

func Validate(v interface{}) error {
	err := validate(reflect.ValueOf(v))
	if err != nil {
		if errors.As(err, &ValidationErrors{}) {
			validateErr, _ := err.(ValidationErrors)
			return validateErr
		}
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
	if reflectType.Kind() == reflect.Struct {
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
			fieldName := fieldType.Name
			var err error
			switch fieldValue.Kind() {
			case reflect.String:
				err = ValidateString(fieldName, fieldValue.String(), tag)
			case reflect.Int:
				err = ValidateInteger(fieldName, int(fieldValue.Int()), tag)
			case reflect.Slice:
				sliceSize := fieldValue.Len()
				for i := range sliceSize {
					switch fieldType.Type.Elem().Kind() {
					case reflect.String:
						err = ValidateString(fieldName, fieldValue.Index(i).String(), tag)
					case reflect.Int:
						err = ValidateInteger(fieldName, int(fieldValue.Index(i).Int()), tag)
					}
					if err != nil {
						break
					}
				}
			case reflect.Struct:
				err = validate(fieldValue)
			}
			if err != nil {
				if errors.As(err, &ValidationErrors{}) {
					validateErr, _ := err.(ValidationErrors)
					validationErrors = append(validationErrors, validateErr...)
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

func ValidateString(fieldName string, value string, tagValue string) error {
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

func ValidateInteger(fieldName string, value int, tagValue string) error {
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
