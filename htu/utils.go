package htu

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ToStd(h Handler) http.Handler {
	return errHandler{h}
}

type Handler interface {
	ServeHTTPErr(http.ResponseWriter, *http.Request) error
}

type errHandler struct {
	inner Handler
}

func (e errHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := e.inner.ServeHTTPErr(w, r)
	if err == nil {
		return
	}

	var (
		code    = http.StatusInternalServerError
		userMsg string
		logErr  = err
	)

	var hErr HTTPError
	if errors.As(err, &hErr) {
		if hErr.Code != 0 {
			code = hErr.Code
		}
		if hErr.UserErr != nil {
			userMsg = hErr.UserErr.Error()
		}
		logErr = hErr.InternalErr
	}

	if userMsg == "" {
		userMsg = http.StatusText(code)
	}

	w.WriteHeader(code)
	io.WriteString(w, userMsg)
	log.Printf("%s %s: %v", r.Method, r.URL.Path, logErr)
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	errHandler{h}.ServeHTTP(w, r)
}

func (h HandlerFunc) ServeHTTPErr(w http.ResponseWriter, r *http.Request) error {
	return h(w, r)
}

type HTTPError struct {
	Code        int
	UserErr     error
	InternalErr error
}

func (e HTTPError) Error() string {
	if e.UserErr != nil {
		return e.UserErr.Error()
	}
	return fmt.Sprintf("HTTP %d", e.Code)
}

func Int64Param(r *http.Request, name string) (int64, error) {
	v := chi.URLParam(r, name)
	ret, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, err
	}
	return ret, nil
}
