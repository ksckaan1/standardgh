package standardgh

import (
	"context"
	"encoding/json/v2"
	"errors"
	"net/http"
	"strings"

	"github.com/ksckaan1/m2s"
)

func GH[Req any, Res any](handlerFunc func(context.Context, *Req) (Res, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(Req)

		if err := bindURIFromRequest(r, req); err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := parseBody(r, req); err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := bindHeaders(r, req); err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := bindQuery(r, req); err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, status, err := handlerFunc(r.Context(), req)
		if err != nil {
			writeError(w, err.Error(), status)
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

	switch {
	case contentType == "application/json" && isBodyFilled:
		return json.UnmarshalRead(r.Body, out)
	case strings.HasPrefix(contentType, "multipart/form-data") && isBodyFilled:
		if err := r.ParseMultipartForm(1024); err != nil {
			return err
		}
		return m2s.Convert(r.MultipartForm, out)
	default:
		return errors.New("unsupported content type")
	}
}
