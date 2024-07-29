package form

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func setStringToValue(val string, value reflect.Value, opts *Options) error {
	switch value.Kind() {
	case reflect.Int:
		return setIntField(val, 0, value)
	case reflect.Int8:
		return setIntField(val, 8, value)
	case reflect.Int16:
		return setIntField(val, 16, value)
	case reflect.Int32:
		return setIntField(val, 32, value)
	case reflect.Int64:
		switch value.Interface().(type) {
		case time.Duration:
			return setTimeDuration(val, value)
		}
		return setIntField(val, 64, value)
	case reflect.Uint:
		return setUintField(val, 0, value)
	case reflect.Uint8:
		return setUintField(val, 8, value)
	case reflect.Uint16:
		return setUintField(val, 16, value)
	case reflect.Uint32:
		return setUintField(val, 32, value)
	case reflect.Uint64:
		return setUintField(val, 64, value)
	case reflect.Bool:
		return setBoolField(val, value)
	case reflect.Float32:
		return setFloatField(val, 32, value)
	case reflect.Float64:
		return setFloatField(val, 64, value)
	case reflect.String:
		value.SetString(val)
	case reflect.Struct:
		switch value.Interface().(type) {
		case time.Time:
			return setTimeField(val, value, opts)
		}
		return json.Unmarshal([]byte(val), value.Addr().Interface())
	case reflect.Map:
		return json.Unmarshal([]byte(val), value.Addr().Interface())
	default:
		return fmt.Errorf("unsupported type %s", value.Type().Name())
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setTimeField(val string, value reflect.Value, opts *Options) error {

	timeFormat := time.RFC3339
	timeLocation := time.Local

	if opts != nil {
		timeFormat = opts.TimeFormat
		timeLocation = opts.TimeLocation
	}

	switch tf := strings.ToLower(timeFormat); tf {
	case "unix", "unixnano":
		tv, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}

		d := time.Duration(1)
		if tf == "unixnano" {
			d = time.Second
		}

		t := time.Unix(tv/int64(d), tv%int64(d))
		value.Set(reflect.ValueOf(t))
		return nil
	}

	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	t, err := time.ParseInLocation(timeFormat, val, timeLocation)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}

func setStringToArray(vals []string, value reflect.Value, opts *Options) error {
	for i, s := range vals {
		err := setStringToValue(s, value.Index(i), opts)
		if err != nil {
			return err
		}
	}
	return nil
}

func setStringToSlice(vals []string, value reflect.Value, opts *Options) error {
	slice := reflect.MakeSlice(value.Type(), len(vals), len(vals))
	err := setStringToArray(vals, slice, opts)
	if err != nil {
		return err
	}
	value.Set(slice)
	return nil
}

func setTimeDuration(val string, value reflect.Value) error {
	d, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(d))
	return nil
}

func setFileToArray(vals []*multipart.FileHeader, value reflect.Value, opts *Options) error {
	for i, s := range vals {
		err := setFileToValue(s, value.Index(i), opts)
		if err != nil {
			return err
		}
	}
	return nil
}

func setFileToSlice(vals []*multipart.FileHeader, value reflect.Value, opts *Options) error {
	slice := reflect.MakeSlice(value.Type(), len(vals), len(vals))
	err := setFileToArray(vals, slice, opts)
	if err != nil {
		return err
	}
	value.Set(slice)
	return nil
}

func setFileToValue(val *multipart.FileHeader, value reflect.Value, _ *Options) error {
	_, ok := value.Interface().(*multipart.FileHeader)
	if !ok {
		return fmt.Errorf("want type *multipart.FileHeader got %s", value.Type().String())
	}
	value.Set(reflect.ValueOf(val))
	return nil
}
