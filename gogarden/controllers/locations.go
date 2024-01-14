package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.universe.tf/garden/gogarden/db"
	"go.universe.tf/garden/gogarden/views"
)

type locations struct {
	db *db.DB
}

func Locations(r *chi.Mux, db *db.DB) {
	s := locations{db}
	r.Get("/locations", s.listLocations)
}

func (s *locations) listLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := s.db.ListLocations(r.Context())
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	render(w, r, views.Locations(locations))
}
