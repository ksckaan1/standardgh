package standardgh

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

func bindHeaders(r *http.Request, out any) error {
	outVal := reflectValue(out)
	if !outVal.IsValid() || outVal.Kind() != reflect.Struct {
		return nil
	}

	t := outVal.Type()
	for i := range outVal.NumField() {
		field := t.Field(i)
		fieldValue := outVal.Field(i)

		tag := field.Tag.Get("header")
		if tag == "" {
			continue
		}

		name, opts := parseTag(tag)
		value := getHeader(r, name)

		if value == "" {
			if dv, ok := opts["default"]; ok {
				value = dv
			} else if _, ok := opts["required"]; ok {
				return &EmptyFieldError{Key: name, Source: "header"}
			}
		}

		if value != "" && fieldValue.CanSet() {
			if err := setField(fieldValue, value); err != nil {
				return fmt.Errorf("%s header: %w", name, err)
			}
		}
	}

	return nil
}

func getHeader(r *http.Request, name string) string {
	lower := strings.ToLower(name)
	for key, vals := range r.Header {
		if strings.ToLower(key) == lower && len(vals) > 0 {
			return vals[0]
		}
	}
	return ""
}
