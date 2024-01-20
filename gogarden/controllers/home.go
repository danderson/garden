package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.universe.tf/garden/gogarden/db"
	"go.universe.tf/garden/gogarden/views"
)

func Home(r *chi.Mux, db *db.DB) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, r, views.Home())
	})
}
