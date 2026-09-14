package standardgh

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
	"uuid"
)

func GHforSSE[Req any, Data any](
	retry time.Duration,
	handlerFunc func(context.Context, *Req, func(name string, data Data) error) error,
) http.HandlerFunc {
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

		if err := bindCookies(r, req); err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := validateStruct(req); err != nil {
			writeError(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			writeError(w, "streaming not supported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")
		w.Header().Set("Retry-After", strconv.Itoa(int(retry.Milliseconds())/1000))

		w.WriteHeader(http.StatusOK)
		flusher.Flush()

		ctx := r.Context()

		var mu sync.Mutex

		send := func(name string, data Data) error {
			mu.Lock()
			defer mu.Unlock()

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			eventID := uuid.NewV7()

			buf, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("marshal error: %w", err)
			}

			if name != "" {
				if _, err := fmt.Fprintf(w, "event: %s\n", name); err != nil {
					return err
				}
			}

			if _, err := fmt.Fprintf(w, "id: %s\n", eventID.String()); err != nil {
				return err
			}

			if _, err := fmt.Fprintf(w, "data: %s\n\n", buf); err != nil {
				return err
			}

			flusher.Flush()
			return nil
		}

		_ = handlerFunc(ctx, req, send)
	}
}
