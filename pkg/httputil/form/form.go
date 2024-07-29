package form

import (
	"errors"
	"mime/multipart"
	"reflect"
	"time"
)

type Value map[string]Item

type Options struct {
	Field        string
	TimeFormat   string
	TimeLocation *time.Location
}

type Item interface {
	SetValue(value reflect.Value, opt *Options) error
}

type StringSlice []string

func (s StringSlice) SetValue(value reflect.Value, opt *Options) error {
	typ := value.Type()
	switch typ.Kind() {
	case reflect.Array:
		return setStringToArray(s, value, opt)
	case reflect.Slice:
		return setStringToSlice(s, value, opt)
	default:
		if len(s) > 0 {
			return setStringToValue(s[len(s)-1], value, opt)
		}
		return nil
	}
}

type FileSlice []*multipart.FileHeader

func (s FileSlice) SetValue(value reflect.Value, opt *Options) error {
	typ := value.Type()
	switch typ.Kind() {
	case reflect.Array:
		return setFileToArray(s, value, opt)
	case reflect.Slice:
		return setFileToSlice(s, value, opt)
	default:
		if len(s) > 0 {
			return setFileToValue(s[len(s)-1], value, opt)
		}
		return nil
	}
}

func NewOptionFromStruct(field reflect.StructField, defOpt *Options) (*Options, error) {
	opt := NewDefaultOption()
	opt.Field = defOpt.Field
	opt.TimeFormat = defOpt.TimeFormat
	opt.TimeLocation = defOpt.TimeLocation
	name := field.Tag.Get("time_location")
	if name != "" {
		loc, err := time.LoadLocation(name)
		if err != nil {
			return nil, err
		}
		opt.TimeLocation = loc
	}
	fmt := field.Tag.Get("time_format")
	if name != "" {
		opt.TimeFormat = fmt
	}
	return opt, nil
}

func NewDefaultOption() *Options {
	return &Options{
		Field:        "form",
		TimeLocation: time.Local,
		TimeFormat:   time.RFC3339,
	}
}

func getName(name string, _ reflect.Value, field reflect.StructField, opt *Options) string {
	if v := field.Tag.Get(opt.Field); v != "" {
		return v
	}
	return name
}

func mapping(input map[string]Item, val reflect.Value, opt *Options) error {
	typ := val.Type()
	if typ.Kind() != reflect.Pointer {
		return errors.New("output must be pointer")
	}

	val = val.Elem()
	typ = val.Type()
	if val.Kind() != reflect.Struct {
		return errors.New("output must be pointer of struct")
	}

	n := val.NumField()
	for i := 0; i < n; i++ {
		field := typ.Field(i)
		fieldOpt, err := NewOptionFromStruct(field, opt)
		if err != nil {
			return err
		}
		item := val.Field(i)
		name := getName(field.Name, item, field, fieldOpt)
		if v, ok := input[name]; ok {
			if err := v.SetValue(item, fieldOpt); err != nil {
				return err
			}
		}
	}
	return nil
}

func MappingForm(input map[string]Item, out interface{}) error {
	return MappingFormOptions(input, out, NewDefaultOption())
}

func MappingFormOptions(input map[string]Item, out interface{}, opt *Options) error {
	return mapping(input, reflect.ValueOf(out), opt)
}

func NewForm(value map[string][]string) map[string]Item {
	val := map[string]Item{}
	for k, v := range value {
		val[k] = StringSlice(v)
	}
	return val
}

func NewMultipartForm(form *multipart.Form) map[string]Item {
	val := map[string]Item{}
	for k, v := range form.Value {
		val[k] = StringSlice(v)
	}
	for k, v := range form.File {
		val[k] = FileSlice(v)
	}
	return val
}
