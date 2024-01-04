package htu

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
)

func ErrHandler(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err == nil {
			w.WriteHeader(200)
			return
		}
		var hErr httpErr
		if errors.As(err, &hErr) {
			w.WriteHeader(hErr.code)
			io.WriteString(w, hErr.userMsg)
			slog.Error("Handler error", "err", hErr.err, "usermsg", hErr.userMsg, "code", hErr.code)
		} else {
			w.WriteHeader(500)
			slog.Error("Handler error", "err", hErr.err, "code", 500)
		}
	})
}

func BadRequest(userMsg string, internalErr error) error {
	return HTTPErr(400, userMsg, internalErr)
}

func HTTPErr(code int, userMsg string, internalErr error) error {
	return httpErr{code, userMsg, internalErr}
}

type httpErr struct {
	code    int
	userMsg string
	err     error
}

func (e httpErr) Error() string {
	return fmt.Sprintf("%s (%d)", e.userMsg, e.code)
}

func (e httpErr) Unwrap() error {
	return e.err
}

func JSONBody(r *http.Request, out any) error {
	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		return fmt.Errorf("wrong content type %q", ct)
	}
	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		return err
	}
	return nil
}

func Int64Param(r *http.Request, name string) (int64, error) {
	v := chi.URLParam(r, name)
	ret, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func RespondJSON(w http.ResponseWriter, obj any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(obj)
}
