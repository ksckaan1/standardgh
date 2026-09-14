package standardgh

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type routePattern struct {
	segments []routeSegment
}

type routeSegment struct {
	isParam bool
	param   string
	literal string
}

func parseRoutePattern(pattern string) *routePattern {
	segments := strings.Split(pattern, "/")
	rp := &routePattern{}

	for _, seg := range segments {
		if seg == "" {
			continue
		}
		if strings.HasPrefix(seg, ":") {
			rp.segments = append(rp.segments, routeSegment{
				isParam: true,
				param:   seg[1:],
			})
		} else if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
			rp.segments = append(rp.segments, routeSegment{
				isParam: true,
				param:   seg[1 : len(seg)-1],
			})
		} else {
			rp.segments = append(rp.segments, routeSegment{
				isParam: false,
				literal: seg,
			})
		}
	}

	return rp
}

func (rp *routePattern) extractParams(path string) map[string]string {
	parts := strings.Split(path, "/")
	params := make(map[string]string)

	j := 0
	for _, part := range parts {
		if part == "" {
			continue
		}
		if j >= len(rp.segments) {
			break
		}
		seg := rp.segments[j]
		if seg.isParam {
			params[seg.param] = part
		}
		j++
	}

	return params
}

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

		name, opts := parseTag(tag)
		if name == "" {
			name = strings.ToLower(field.Name)
		}

		value := r.PathValue(name)

		if value == "" {
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
