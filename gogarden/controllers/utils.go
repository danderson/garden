package controllers

import (
	"net/http"

	"github.com/a-h/templ"
	"go.universe.tf/garden/gogarden/views"
)

func render(w http.ResponseWriter, r *http.Request, page templ.Component) {
	views.Root(page).Render(r.Context(), w)
}
