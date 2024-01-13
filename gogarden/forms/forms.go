package forms

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// Form decorates an arbitrary struct with validation errors, and
// provides helpers to parse and validate HTML form data.
type Form[T any] struct {
	Data   T
	Fields map[string]Field
	Errors []string
}

// Field is a struct field and its validation errors, if any.
type Field struct {
	Value  any
	Errors []string
}

// FromStruct returns a Form initialized with st's data, and no
// validation errors.
func FromStruct[T any](st T) *Form[T] {
	ret := &Form[T]{
		Data:   st,
		Fields: map[string]Field{},
	}
	err := eachStructField(st, func(fi reflect.StructField, fv reflect.Value) error {
		ret.Fields[fi.Name] = Field{
			Value: fv.Interface(),
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return ret
}

// FromForm returns a Form for st, with form values from r patched in.
func FromForm[T any](st T, r *http.Request) (*Form[T], error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	ret := &Form[T]{
		Data:   st,
		Fields: map[string]Field{},
	}
	err := eachStructField(ret.Data, func(fi reflect.StructField, fv reflect.Value) error {
		name := fi.Name
		if !r.Form.Has(name) {
			name = strings.ToLower(name)
		}
		if !r.Form.Has(name) {
			return nil
		}
		val := r.Form.Get(name)
		err := castValue(val, fv)
		ret.Fields[fi.Name] = Field{
			Value:  fv.Interface(),
			Errors: []string{err.Error()},
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return ret, nil
}

// eachStructField calls fn for every field in st, which must be a
// struct or a pointer to a struct. If fn returns an error,
// eachStructField returns immediately with that error.
func eachStructField[T any](st T, fn func(fi reflect.StructField, fv reflect.Value) error) error {
	v := reflect.ValueOf(st)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return errors.New("input is not a struct")
	}
	for i := 0; i < t.NumField(); i++ {
		fi := t.Field(i)
		fv := v.Field(i)
		if err := fn(fi, fv); err != nil {
			return err
		}
	}
	return nil
}

func castValue(raw string, dest reflect.Value) error {
	switch dest.Kind() {
	// TODO: cast stuff
	default:
		return fmt.Errorf("unhandled form kind %v", dest.Kind())
	}
}
