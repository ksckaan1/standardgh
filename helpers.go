package standardgh

import (
	"reflect"
	"strings"
)

func reflectValue(out any) reflect.Value {
	v := reflect.ValueOf(out)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return reflect.Value{}
		}
		return v.Elem()
	}
	return v
}

func isZeroValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Array:
		return v.Len() == 0
	case reflect.Struct:
		return v.IsZero()
	default:
		return false
	}
}

func parseTag(tag string) (string, map[string]string) {
	parts := strings.Split(tag, ",")
	name := parts[0]
	opts := make(map[string]string)

	for _, part := range parts[1:] {
		if key, value, ok := strings.Cut(part, ":"); ok {
			opts[key] = value
		} else {
			opts[part] = ""
		}
	}

	return name, opts
}
