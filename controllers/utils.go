package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"go.universe.tf/garden/htu"
	"go.universe.tf/garden/views"
)

func chiFn(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return htu.HandlerFunc(fn).ServeHTTP
}

func render(w http.ResponseWriter, r *http.Request, page templ.Component) {
	views.Root(page).Render(r.Context(), w)
}

func badRequest(err error) error {
	return htu.HTTPError{
		Code:        http.StatusBadRequest,
		InternalErr: err,
	}
}

func badRequestf(msg string, args ...any) error {
	return badRequest(fmt.Errorf(msg, args...))
}

func notFound(err error) error {
	return htu.HTTPError{
		Code:        http.StatusNotFound,
		InternalErr: err,
	}
}

func notFoundf(msg string, args ...any) error {
	return notFound(fmt.Errorf(msg, args...))
}

func internalError(err error) error {
	return htu.HTTPError{
		Code:        http.StatusInternalServerError,
		InternalErr: err,
	}
}

func internalErrorf(msg string, args ...any) error {
	return internalError(fmt.Errorf(msg, args...))
}

func dbGetError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return notFound(err)
	} else {
		return internalError(err)
	}
}

func dbGetErrorf(msg string, args ...any) error {
	return dbGetError(fmt.Errorf(msg, args...))
}
