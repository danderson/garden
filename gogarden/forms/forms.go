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
	Value   string
	Errors  []string
	Options []SelectOption // for <select>
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
		val, err := toFormValue(fv)
		if err != nil {
			panic(err)
		}
		ret.Fields[fi.Name] = Field{
			ID:      fi.Name,
			Value:   val,
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
		val := ""
		var errs []string
		if r.Form.Has(name) {
			val = r.Form.Get(name)
			if err := fromFormValue(val, fv); err != nil {
				errs = []string{err.Error()}
			}
		}
		ret.Fields[fi.Name] = Field{
			ID:      fi.Name,
			Value:   val,
			Errors:  errs,
			Options: selectOptions(fv),
		}
		log.Print(ret.Fields[fi.Name])
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

func (f *Form) SetSelectOptions(field string, options []SelectOption) {
	fd, ok := f.Fields[field]
	if !ok {
		panic(fmt.Sprintf("added options on unknown form field %q", field))
	}
	fd.Options = options
	f.Fields[field] = fd
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

func fromFormValue(raw string, dest reflect.Value) error {
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
		destp := reflect.New(dest.Type().Elem())
		if err := fromFormValue(raw, destp.Elem()); err != nil {
			return err
		}
		dest.Set(destp)
		return nil
	case reflect.Int64:
		if raw == "" {
			dest.SetInt(0)
			return nil
		}
		i, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return err
		}
		dest.SetInt(i)
		return nil
	case reflect.String:
		dest.Set(reflect.ValueOf(raw))
		return nil
	default:
	}
	return fmt.Errorf("unhandled form kind %v", dest.Kind())
}

func toFormValue(src reflect.Value) (string, error) {
	if m, ok := src.Addr().Interface().(encoding.TextMarshaler); ok {
		ret, err := m.MarshalText()
		if err != nil {
			return "", err
		}
		return string(ret), nil
	}

	switch src.Kind() {
	case reflect.Pointer:
		if src.IsZero() {
			return "", nil
		}
		return toFormValue(src.Elem())
	case reflect.Int64:
		i := src.Interface().(int64)
		return strconv.FormatInt(i, 10), nil
	case reflect.String:
		return src.Interface().(string), nil
	default:
	}
	return "", fmt.Errorf("unhandled form value %v (%v)", src, src.Type())
}

type SelectOption struct {
	Value string
	Label string
}

// A Selecter can provide a list of available options for a <select>
// HTML form input.
type Selecter interface {
	SelectOptions() []SelectOption
}

func selectOptions(v reflect.Value) []SelectOption {
	if s, ok := v.Interface().(Selecter); ok {
		return s.SelectOptions()
	}
	return nil
}
