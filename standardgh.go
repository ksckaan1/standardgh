package standardgh

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/ksckaan1/m2s"
)

func GH[Req any, Res any](handlerFunc func(context.Context, *Req) (Res, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(Req)
		reqVal := reflectValue(req)

		if err := bindAll(r, reqVal); err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := parseBody(r, req); err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := validateStruct(req); err != nil {
			writeError(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		res, status, err := handlerFunc(r.Context(), req)
		if err != nil {
			writeError(w, err.Error(), status)
			return
		}

		if err := encodeResponseHeaders(w, res); err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := encodeResponseCookies(w, res); err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(status)

		if status == http.StatusNoContent {
			return
		}

		if err = json.MarshalWrite(w, res); err != nil {
			writeError(w, err.Error(), status)
			return
		}
	}
}

func bindAll(r *http.Request, reqVal reflect.Value) error {
	if !reqVal.IsValid() || reqVal.Kind() != reflect.Struct {
		return nil
	}

	t := reqVal.Type()

	cookies := parseCookies(r)

	for i := range reqVal.NumField() {
		field := t.Field(i)
		fieldValue := reqVal.Field(i)

		if !field.IsExported() || !fieldValue.CanSet() {
			continue
		}

		if tag := field.Tag.Get("uri"); tag != "" {
			if err := bindURIField(r, field, fieldValue, tag); err != nil {
				return err
			}
		}

		if tag := field.Tag.Get("header"); tag != "" {
			if err := bindHeaderField(r, field, fieldValue, tag); err != nil {
				return err
			}
		}

		if tag := field.Tag.Get("cookie"); tag != "" {
			if err := bindCookieField(field, fieldValue, tag, cookies); err != nil {
				return err
			}
		}
	}

	return bindQuery(r, reqVal.Addr().Interface())
}

func bindURIField(r *http.Request, field reflect.StructField, fieldValue reflect.Value, tag string) error {
	name, opts := parseTag(tag)
	if name == "" {
		name = strings.ToLower(field.Name)
	}

	value := r.PathValue(name)

	if value == "" {
		if dv, ok := opts["default"]; ok {
			value = dv
		} else if _, ok := opts["required"]; ok {
			return &EmptyFieldError{Key: name, Source: "uri"}
		}
	}

	if value != "" {
		return setField(fieldValue, value)
	}
	return nil
}

func bindHeaderField(r *http.Request, field reflect.StructField, fieldValue reflect.Value, tag string) error {
	name, opts := parseTag(tag)
	if name == "" {
		return nil
	}

	value := getHeader(r, name)

	if value == "" {
		if dv, ok := opts["default"]; ok {
			value = dv
		} else if _, ok := opts["required"]; ok {
			return &EmptyFieldError{Key: name, Source: "header"}
		}
	}

	if value != "" {
		return setField(fieldValue, value)
	}
	return nil
}

func bindCookieField(field reflect.StructField, fieldValue reflect.Value, tag string, cookies map[string]string) error {
	name, opts := parseTag(tag)
	if name == "" {
		name = strings.ToLower(field.Name)
	}

	value, ok := cookies[name]
	if !ok || value == "" {
		if dv, ok := opts["default"]; ok {
			value = dv
		} else if _, ok := opts["required"]; ok {
			return &EmptyFieldError{Key: name, Source: "cookie"}
		}
	}

	if value != "" {
		if err := setField(fieldValue, value); err != nil {
			return fmt.Errorf("%s cookie: %w", name, err)
		}
	}
	return nil
}

func writeError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.MarshalWrite(w, &struct {
		Error string `json:"error"`
	}{Error: msg})
}

func parseBody(r *http.Request, out any) error {
	contentType := r.Header.Get("Content-Type")
	isBodyFilled := r.Body != nil && r.ContentLength > 0

	if !isBodyFilled {
		return nil
	}

	switch {
	case contentType == "application/json":
		return json.UnmarshalRead(r.Body, out)
	case strings.HasPrefix(contentType, "multipart/form-data"):
		if err := r.ParseMultipartForm(1024); err != nil {
			return err
		}
		return m2s.Convert(r.MultipartForm, out)
	default:
		return errors.New("unsupported content type")
	}
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
