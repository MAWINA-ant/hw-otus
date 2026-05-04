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

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	var validationErrors ValidationErrors
	reflectValue := reflect.ValueOf(v)
	reflectType := reflectValue.Type()
	if reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectType = reflectType.Elem()
	}
	if reflectType.Kind() == reflect.Struct {
		fieldCount := reflectType.NumField()
		for i := range fieldCount {
			field := reflectType.Field(i)
			// проверка на публичность поля структуры
			if field.PkgPath != "" {
				continue
			}
			tag := field.Tag.Get("validate")
			// проверка на наличие правил "validate"
			if tag == "" || tag == "-" {
				continue
			}
			fieldName := field.Name
			fieldType := field.Type
			fieldValue := reflectValue.Field(i)
			var err error
			switch fieldType.Kind() {
			case reflect.String:
				err = ValidateString(fieldName, fieldValue.String(), tag)
			case reflect.Int:
				err = ValidateInteger(fieldName, int(fieldValue.Int()), tag)
			case reflect.Slice:
				sliceSize := fieldValue.Len()
				for i := range sliceSize {
					switch field.Type.Elem().Kind() {
					case reflect.String:
						err = ValidateString(fieldName, fieldValue.Index(i).String(), tag)
					case reflect.Int:
						err = ValidateInteger(fieldName, int(fieldValue.Index(i).Int()), tag)
					}
					if err != nil {
						var validateErr ValidationErrors
						if errors.Is(err, validateErr) {
							validateErr, _ = err.(ValidationErrors)
							validationErrors = append(validationErrors, validateErr...)
						}
						return err
					}
				}
			case reflect.Struct:
				err = Validate(fieldValue)
			}
			if err != nil {
				var validateErr ValidationErrors
				if errors.Is(err, validateErr) {
					validateErr, _ = err.(ValidationErrors)
					validationErrors = append(validationErrors, validateErr...)
				}
				return err
			}
		}
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
				validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: fmt.Errorf("\"len\" rule is not fulfilled")})
			}
		case "regexp":
			re := regexp.MustCompile(ruleValue)
			if !re.MatchString(value) {
				validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: fmt.Errorf("\"regexp\" rule is not fulfilled")})
			}
		case "in":
			ruleValueSlice := strings.Split(ruleValue, ",")
			if !slices.Contains(ruleValueSlice, value) {
				validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: fmt.Errorf("\"in\" rule is not fulfilled")})
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
				validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: fmt.Errorf("\"min\" rule is not fulfilled")})
			}
		case "max":
			ruleValueInt, err := strconv.Atoi(ruleValue)
			if err != nil {
				return fmt.Errorf("couldn't parse \"max\" rule value to int")
			}
			if value > ruleValueInt {
				validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: fmt.Errorf("\"max\" rule is not fulfilled")})
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
			validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: fmt.Errorf("\"in\" rule is not fulfilled")})
		}
	}
	if len(validationErrors) == 0 {
		return nil
	}
	return validationErrors
}
