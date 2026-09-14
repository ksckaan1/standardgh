package standardgh

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

func bindCookies(r *http.Request, out any) error {
	outVal := reflectValue(out)
	if !outVal.IsValid() {
		return nil
	}

	cookies := parseCookies(r)

	if outVal.Kind() == reflect.Map {
		return bindCookieMap(outVal, cookies)
	}

	if outVal.Kind() != reflect.Struct {
		return nil
	}

	t := outVal.Type()
	for i := range outVal.NumField() {
		field := t.Field(i)
		fieldValue := outVal.Field(i)

		tag := field.Tag.Get("cookie")
		if tag == "" {
			continue
		}

		name, opts := parseTag(tag)
		if name == "" {
			name = strings.ToLower(field.Name)
		}

		value, ok := cookies[name]
		if !ok || value == "" {
			if dv, dOk := opts["default"]; dOk {
				value = dv
			} else if _, rOk := opts["required"]; rOk {
				return &EmptyFieldError{Key: name, Source: "cookie"}
			}
		}

		if value != "" && fieldValue.CanSet() {
			if err := setField(fieldValue, value); err != nil {
				return fmt.Errorf("%s cookie: %w", name, err)
			}
		}
	}

	return nil
}

func parseCookies(r *http.Request) map[string]string {
	cookies := make(map[string]string)

	for _, cookie := range r.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}

	return cookies
}

func bindCookieMap(field reflect.Value, cookies map[string]string) error {
	mapType := field.Type()
	if mapType.Key().Kind() != reflect.String {
		return fmt.Errorf("map key must be string, got %s", mapType.Key())
	}

	if field.IsNil() {
		field.Set(reflect.MakeMap(mapType))
	}

	valType := mapType.Elem()

	for key, value := range cookies {
		if valType.Kind() == reflect.Slice {
			slice := reflect.MakeSlice(valType, 1, 1)
			if err := setField(slice.Index(0), value); err != nil {
				return err
			}
			field.SetMapIndex(reflect.ValueOf(key), slice)
		} else {
			val := reflect.New(valType).Elem()
			if err := setField(val, value); err != nil {
				return err
			}
			field.SetMapIndex(reflect.ValueOf(key), val)
		}
	}

	return nil
}
