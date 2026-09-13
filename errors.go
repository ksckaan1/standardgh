package standardgh

import (
	"fmt"
	"reflect"
)

type ConversionError struct {
	Key   string
	Type  reflect.Type
	Index int
	Err   error
}

func (e *ConversionError) Error() string {
	if e.Index >= 0 {
		return fmt.Sprintf("error converting value for %q at index %d: %v", e.Key, e.Index, e.Err)
	}
	return fmt.Sprintf("error converting value for %q: %v", e.Key, e.Err)
}

func (e *ConversionError) Unwrap() error { return e.Err }

type EmptyFieldError struct {
	Key    string
	Source string
}

func (e *EmptyFieldError) Error() string {
	return fmt.Sprintf("%s %s is required", e.Key, e.Source)
}

type MultiError struct {
	Errors map[string]error
}

func (e *MultiError) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}
	for k, v := range e.Errors {
		return fmt.Sprintf("%s: %v", k, v)
	}
	return ""
}

func (e *MultiError) Unwrap() error {
	for _, v := range e.Errors {
		return v
	}
	return nil
}
