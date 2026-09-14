package standardgh

import (
	"encoding"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const maxSliceSize = 16000

var errUnmatchedBrackets = errors.New("unmatched brackets in query parameter")

func bindQuery(r *http.Request, out any) error {
	outVal := reflectValue(out)
	if !outVal.IsValid() {
		return nil
	}

	params, err := parseQueryParams(r)
	if err != nil {
		return err
	}

	if outVal.Kind() == reflect.Map {
		return bindMap(outVal, params)
	}

	if outVal.Kind() != reflect.Struct {
		return nil
	}

	fields := collectFields(outVal.Type(), "")
	multiErr := &MultiError{Errors: make(map[string]error)}

	for key, values := range params {
		finfo, ok := fields[key]
		if !ok {
			continue
		}

		vals := values
		if finfo.field.Type.Kind() == reflect.Slice {
			vals = splitCommaValues(values)
		}

		if err := setByPath(outVal, finfo.path, vals); err != nil {
			if _, ok := errors.AsType[*ConversionError](err); ok {
				multiErr.Errors[key] = err
			} else {
				return err
			}
		}
	}

	for _, finfo := range fields {
		if err := checkRequired(outVal, finfo, params); err != nil {
			return err
		}
		if err := applyDefault(outVal, finfo); err != nil {
			return err
		}
	}

	if len(multiErr.Errors) > 0 {
		return multiErr
	}

	return nil
}

func checkRequired(root reflect.Value, finfo *fieldInfo, params map[string][]string) error {
	if !finfo.isRequired {
		return nil
	}

	val := root.FieldByIndex(finfo.path)
	if val.IsValid() && !isZeroValue(val) {
		return nil
	}

	if _, has := params[finfo.alias]; has {
		return nil
	}

	if _, hasDefault := finfo.options["default"]; hasDefault {
		return nil
	}

	return &EmptyFieldError{Key: finfo.alias, Source: "query"}
}

func applyDefault(root reflect.Value, finfo *fieldInfo) error {
	dv, ok := finfo.options["default"]
	if !ok {
		return nil
	}

	val := root.FieldByIndex(finfo.path)
	if val.IsValid() && !isZeroValue(val) {
		return nil
	}

	return setDefaultByPath(root, finfo.path, dv, finfo.field.Type)
}

// ─── Field Collection ────────────────────────────────────────────────────────

type fieldInfo struct {
	alias      string
	field      reflect.StructField
	path       []int
	options    map[string]string
	isRequired bool
}

func collectFields(t reflect.Type, prefix string) map[string]*fieldInfo {
	fields := make(map[string]*fieldInfo)

	for f := range t.Fields() {
		f := f

		if !f.IsExported() {
			continue
		}

		if f.Anonymous {
			collectEmbeddedFields(f, prefix, fields)
			continue
		}

		tag := f.Tag.Get("query")
		alias, opts := parseTag(tag)
		if alias == "" {
			alias = strings.ToLower(f.Name)
		}
		alias = prefix + alias

		if f.Type.Kind() == reflect.Struct && !implementsTextUnmarshaler(f.Type) {
			collectNestedFields(f, alias, fields)
			continue
		}

		_, isRequired := opts["required"]
		fields[alias] = &fieldInfo{
			alias:      alias,
			field:      f,
			path:       f.Index,
			options:    opts,
			isRequired: isRequired,
		}
	}

	return fields
}

func collectEmbeddedFields(f reflect.StructField, prefix string, fields map[string]*fieldInfo) {
	embeddedType := f.Type
	if embeddedType.Kind() == reflect.Pointer {
		embeddedType = embeddedType.Elem()
	}

	if embeddedType.Kind() != reflect.Struct {
		return
	}

	for k, v := range collectFields(embeddedType, prefix) {
		v.path = prependPath(f.Index, v.path)
		fields[k] = v
	}
}

func collectNestedFields(f reflect.StructField, prefix string, fields map[string]*fieldInfo) {
	nestedType := f.Type
	if f.Type.Kind() == reflect.Pointer {
		nestedType = f.Type.Elem()
	}

	if nestedType.Kind() != reflect.Struct {
		return
	}

	for k, v := range collectFields(nestedType, prefix+".") {
		v.path = prependPath(f.Index, v.path)
		fields[k] = v
	}
}

func prependPath(base, rest []int) []int {
	path := make([]int, len(base)+len(rest))
	copy(path, base)
	copy(path[len(base):], rest)
	return path
}

// ─── Map Binding ─────────────────────────────────────────────────────────────

func bindMap(field reflect.Value, params map[string][]string) error {
	mapType := field.Type()
	if mapType.Key().Kind() != reflect.String {
		return fmt.Errorf("map key must be string, got %s", mapType.Key())
	}

	if field.IsNil() {
		field.Set(reflect.MakeMap(mapType))
	}

	valType := mapType.Elem()

	for key, values := range params {
		if valType.Kind() == reflect.Slice {
			slice := reflect.MakeSlice(valType, len(values), len(values))
			for i, v := range values {
				if err := setField(slice.Index(i), v); err != nil {
					return err
				}
			}
			field.SetMapIndex(reflect.ValueOf(key), slice)
		} else if len(values) > 0 {
			val := reflect.New(valType).Elem()
			if err := setField(val, values[len(values)-1]); err != nil {
				return err
			}
			field.SetMapIndex(reflect.ValueOf(key), val)
		}
	}

	return nil
}

// ─── Query Parsing ───────────────────────────────────────────────────────────

func parseQueryParams(r *http.Request) (map[string][]string, error) {
	params := make(map[string][]string)
	query := r.URL.RawQuery

	if query == "" {
		return params, nil
	}

	for query != "" {
		var key string
		key, query, _ = strings.Cut(query, "&")
		if key == "" {
			continue
		}

		k, v, _ := strings.Cut(key, "=")

		k, _ = url.QueryUnescape(k)

		var err error
		k, err = parseBracketNotation(k)
		if err != nil {
			return nil, err
		}
		k = strings.ToLower(k)

		v, _ = url.QueryUnescape(v)

		params[k] = append(params[k], v)
	}

	return params, nil
}

func parseBracketNotation(key string) (string, error) {
	if !strings.Contains(key, "[") {
		return key, nil
	}

	var result strings.Builder
	depth := 0

	for i := 0; i < len(key); i++ {
		ch := key[i]
		switch ch {
		case '[':
			if i+1 < len(key) && key[i+1] == ']' {
				i++
				continue
			}
			result.WriteByte('.')
			depth++
		case ']':
			depth--
			if depth < 0 {
				return "", errUnmatchedBrackets
			}
		default:
			result.WriteByte(ch)
		}
	}

	if depth != 0 {
		return "", errUnmatchedBrackets
	}

	return result.String(), nil
}

func splitCommaValues(values []string) []string {
	var result []string
	for _, v := range values {
		if strings.Contains(v, ",") {
			result = append(result, strings.Split(v, ",")...)
		} else {
			result = append(result, v)
		}
	}
	return result
}

// ─── TextUnmarshaler Check ──────────────────────────────────────────────────

func implementsTextUnmarshaler(t reflect.Type) bool {
	return t.Implements(reflect.TypeFor[encoding.TextUnmarshaler]())
}
