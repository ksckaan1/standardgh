package standardgh

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

func bindURIFromRequest(r *http.Request, out any) error {
	outVal := reflectValue(out)
	if !outVal.IsValid() || outVal.Kind() != reflect.Struct {
		return nil
	}

	t := outVal.Type()
	for i := range outVal.NumField() {
		field := t.Field(i)
		fieldValue := outVal.Field(i)

		tag := field.Tag.Get("uri")
		if tag == "" {
			continue
		}

		if err := bindURIField(r, field, fieldValue, tag); err != nil {
			return err
		}
	}

	return nil
}

func bindURIFromParams(params map[string]string, out any) error {
	outVal := reflectValue(out)
	if !outVal.IsValid() || outVal.Kind() != reflect.Struct {
		return nil
	}

	t := outVal.Type()
	for i := range outVal.NumField() {
		field := t.Field(i)
		fieldValue := outVal.Field(i)

		tag := field.Tag.Get("uri")
		if tag == "" {
			continue
		}

		name, opts := parseTag(tag)
		if name == "" {
			name = strings.ToLower(field.Name)
		}

		value, ok := params[name]
		if !ok || value == "" {
			if dv, dOk := opts["default"]; dOk {
				value = dv
			} else if _, rOk := opts["required"]; rOk {
				return &EmptyFieldError{Key: name, Source: "uri"}
			}
		}

		if value != "" && fieldValue.CanSet() {
			if err := setField(fieldValue, value); err != nil {
				return fmt.Errorf("%s uri: %w", name, err)
			}
		}
	}

	return nil
}
