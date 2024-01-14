package forms

import (
	"encoding"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// Form decorates a struct with validation errors, and provides
// helpers to parse and validate HTML form data.
type Form struct {
	Fields map[string]Field
	Errors []string
}

// Field is a struct field and its validation errors, if any.
type Field struct {
	ID      string
	Value   any
	Errors  []string
	Options []string // for <select>
}

// FromStruct returns a Form initialized with st's data, and no
// validation errors.
func New[T any]() *Form {
	var st T
	return FromStruct(&st)
}

func FromStruct[T any](st *T) *Form {
	ret := &Form{
		Fields: map[string]Field{},
	}
	err := eachStructField(st, func(fi reflect.StructField, fv reflect.Value) error {
		ret.Fields[fi.Name] = Field{
			ID:      fi.Name,
			Value:   fv.Interface(),
			Options: selectOptions(fv),
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return ret
}

// FromForm returns a Form for st, with form values from r patched in.
func FromRequest[T any](st *T, r *http.Request) (*T, *Form, error) {
	if err := r.ParseForm(); err != nil {
		return nil, nil, err
	}
	ret := &Form{
		Fields: map[string]Field{},
	}
	err := eachStructField(st, func(fi reflect.StructField, fv reflect.Value) error {
		name := fi.Name
		if !r.Form.Has(name) {
			name = strings.ToLower(name)
		}
		if !r.Form.Has(name) {
			return nil
		}
		val := r.Form.Get(name)
		var errs []string
		if err := castValue(val, fv); err != nil {
			errs = []string{err.Error()}
		}
		ret.Fields[fi.Name] = Field{
			ID:     fi.Name,
			Value:  fv.Interface(),
			Errors: errs,
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return st, ret, nil
}

// HasErrors reports whether the form has any validation errors.
func (f *Form) HasErrors() bool {
	if len(f.Errors) > 0 {
		return true
	}
	for _, f := range f.Fields {
		if len(f.Errors) > 0 {
			return true
		}
	}
	return false
}

// AddError adds a validation error to field.
func (f *Form) AddError(field string, msg string, args ...any) {
	fd, ok := f.Fields[field]
	if !ok {
		panic(fmt.Sprintf("added error on unknown form field %q", field))
	}
	fd.Errors = append(fd.Errors, fmt.Sprintf(msg, args...))
	f.Fields[field] = fd
}

func (f *Form) AddFormError(msg string, args ...any) {
	f.Errors = append(f.Errors, fmt.Sprintf(msg, args...))
}

// eachStructField calls fn for every field in st, which must be a
// pointer to a struct. If fn returns an error, eachStructField
// returns immediately with that error.
func eachStructField[T any](st *T, fn func(fi reflect.StructField, fv reflect.Value) error) error {
	v := reflect.ValueOf(st).Elem()
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
	if um, ok := dest.Addr().Interface().(encoding.TextUnmarshaler); ok {
		return um.UnmarshalText([]byte(raw))
	}

	// Otherwise, handle the basic Go types.
	switch dest.Kind() {
	case reflect.Pointer:
		if raw == "" {
			dest.Set(reflect.Zero(dest.Type()))
			return nil
		}
		destp := reflect.New(dest.Elem().Elem().Type())
		if err := castValue(raw, destp); err != nil {
			return err
		}
		dest.Set(destp.Addr())
		return nil
	case reflect.Int64:
		i, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return err
		}
		dest.SetInt(i)
	case reflect.String:
		dest.Set(reflect.ValueOf(raw))
		return nil
	default:
	}
	return fmt.Errorf("unhandled form kind %v", dest.Kind())
}

// A Selecter can provide a list of available options for a <select>
// HTML form input.
type Selecter interface {
	SelectOptions() []string
}

func selectOptions(v reflect.Value) []string {
	log.Printf("selectOption %s", v.Type().Name())
	if s, ok := v.Interface().(Selecter); ok {
		return s.SelectOptions()
	}
	return nil
}
