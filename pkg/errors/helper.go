package errors

import (
	"errors"
	"fmt"
)

var New = errors.New
var Is = errors.Is
var As = errors.As
var Join = errors.Join

func Errorf(format string, params ...interface{}) error {
	return New(fmt.Sprintf(format, params...))
}

func Wrap(err error, msg string) error {
	return Join(err, errors.New(msg))
}

func Wrapf(err error, format string, params ...interface{}) error {
	return Wrap(err, fmt.Sprintf(format, params...))
}

func Unwrap(err error) []error {
	if v, ok := err.(customError); ok {
		return []error{v.Unwrap()}
	}
	if v, ok := err.(interface{ Unwrap() error }); ok {
		return []error{v.Unwrap(), err}
	}
	if v, ok := err.(interface{ Unwrap() []error }); ok {
		return v.Unwrap()
	}
	return []error{err}
}
