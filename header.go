package standardgh

import (
	"net/http"
	"reflect"
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

		if err := bindHeaderField(r, field, fieldValue, tag); err != nil {
			return err
		}
	}

	return nil
}