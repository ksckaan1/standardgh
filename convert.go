package standardgh

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func setField(field reflect.Value, value string) error {
	if unmarshaler, ok := field.Interface().(encoding.TextUnmarshaler); ok {
		return unmarshaler.UnmarshalText([]byte(value))
	}

	switch field.Kind() {
	case reflect.Pointer:
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		return setField(field.Elem(), value)
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setInt(field, value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setUint(field, value)
	case reflect.Float32, reflect.Float64:
		return setFloat(field, value)
	case reflect.Bool:
		return setBool(field, value)
	case reflect.Slice:
		return setSlice(field, value)
	default:
		return fmt.Errorf("unsupported type: %s", field.Type())
	}
	return nil
}

func setByPath(root reflect.Value, path []int, values []string) error {
	field := root.FieldByIndex(path)

	if field.Kind() == reflect.Pointer {
		return setPtrField(field, values)
	}

	if field.Kind() == reflect.Slice {
		return setSliceField(field, values)
	}

	if len(values) == 0 {
		return nil
	}

	return setField(field, values[len(values)-1])
}

func setPtrField(field reflect.Value, values []string) error {
	if len(values) == 0 || (len(values) == 1 && values[0] == "") {
		return nil
	}

	if field.IsNil() {
		field.Set(reflect.New(field.Type().Elem()))
	}

	return setField(field.Elem(), values[0])
}

func setSliceField(field reflect.Value, values []string) error {
	elemType := field.Type().Elem()
	if elemType.Kind() == reflect.Pointer {
		elemType = elemType.Elem()
	}

	if elemType.Kind() == reflect.Struct && !implementsTextUnmarshaler(elemType) {
		return setIndexedSliceOfStructs(field, values, elemType)
	}

	if len(values) > maxSliceSize {
		return &ConversionError{
			Key:   "slice",
			Type:  field.Type(),
			Index: -1,
			Err:   fmt.Errorf("slice size %d exceeds max %d", len(values), maxSliceSize),
		}
	}

	slice := reflect.MakeSlice(field.Type(), len(values), len(values))
	for i, val := range values {
		elem := slice.Index(i)
		if elem.Kind() == reflect.Pointer {
			if val == "" {
				continue
			}
			elem.Set(reflect.New(elemType))
			elem = elem.Elem()
		}
		if err := setField(elem, val); err != nil {
			return &ConversionError{Key: "value", Type: field.Type(), Index: i, Err: err}
		}
	}

	field.Set(slice)
	return nil
}

func setIndexedSliceOfStructs(field reflect.Value, values []string, elemType reflect.Type) error {
	maxIdx := 0
	for _, val := range values {
		parts := strings.SplitN(val, ".", 2)
		if len(parts) != 2 {
			continue
		}
		idx := parseIndex(parts[0])
		if idx >= 0 && idx+1 > maxIdx {
			maxIdx = idx + 1
		}
	}

	if maxIdx > maxSliceSize {
		return &ConversionError{
			Key:   "slice",
			Type:  field.Type(),
			Index: -1,
			Err:   fmt.Errorf("slice size %d exceeds max %d", maxIdx, maxSliceSize),
		}
	}

	slice := reflect.MakeSlice(field.Type(), maxIdx, maxIdx)
	for i := range maxIdx {
		slice.Index(i).Set(reflect.New(elemType).Elem())
	}

	for _, val := range values {
		parts := strings.SplitN(val, ".", 2)
		if len(parts) != 2 {
			continue
		}

		idx := parseIndex(parts[0])
		if idx < 0 || idx >= maxIdx {
			continue
		}

		subParts := strings.SplitN(parts[1], ".", 2)
		elem := slice.Index(idx)

		if len(subParts) == 2 {
			nested := elem.FieldByName(subParts[0])
			if nested.IsValid() && nested.Kind() == reflect.Struct {
				setIndexedSliceOfStructs(nested, []string{subParts[1]}, nested.Type())
			}
		} else {
			setField(elem.FieldByName(subParts[0]), "")
		}
	}

	field.Set(slice)
	return nil
}

func parseIndex(s string) int {
	idx := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			idx = idx*10 + int(c-'0')
		} else {
			return -1
		}
	}
	return idx
}

func setDefaultByPath(root reflect.Value, path []int, defaultValue string, fieldType reflect.Type) error {
	field := root.FieldByIndex(path)

	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		return setField(field.Elem(), defaultValue)
	}

	if field.Kind() == reflect.Slice {
		return setSlice(field, defaultValue)
	}

	return setField(field, defaultValue)
}

// ─── Primitive Setters ───────────────────────────────────────────────────────

func setInt(field reflect.Value, value string) error {
	n, err := strconv.ParseInt(value, 10, field.Type().Bits())
	if err != nil {
		return err
	}
	field.SetInt(n)
	return nil
}

func setUint(field reflect.Value, value string) error {
	n, err := strconv.ParseUint(value, 10, field.Type().Bits())
	if err != nil {
		return err
	}
	field.SetUint(n)
	return nil
}

func setFloat(field reflect.Value, value string) error {
	n, err := strconv.ParseFloat(value, field.Type().Bits())
	if err != nil {
		return err
	}
	field.SetFloat(n)
	return nil
}

func setBool(field reflect.Value, value string) error {
	if value == "on" {
		field.SetBool(true)
		return nil
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	field.SetBool(b)
	return nil
}

func setSlice(field reflect.Value, value string) error {
	if value == "" {
		return nil
	}

	parts := strings.Split(value, "|")
	slice := reflect.MakeSlice(field.Type(), len(parts), len(parts))

	for i, part := range parts {
		if err := setField(slice.Index(i), strings.TrimSpace(part)); err != nil {
			return err
		}
	}

	field.Set(slice)
	return nil
}
