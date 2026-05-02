package hw09structvalidator

import "reflect"

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
	rv := reflect.ValueOf(v)
	rt := rv.Type()
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		rt = rt.Elem()
	}
	if rt.Kind() == reflect.Struct {

	}
	return nil
}
