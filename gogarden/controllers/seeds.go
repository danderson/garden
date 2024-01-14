package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.universe.tf/garden/gogarden/db"
	"go.universe.tf/garden/gogarden/forms"
	"go.universe.tf/garden/gogarden/htu"
	"go.universe.tf/garden/gogarden/views"
)

type seeds struct {
	db *db.DB
}

func Seeds(r *chi.Mux, db *db.DB) {
	s := &seeds{db}
	r.Get("/seeds", s.listSeeds)
	r.Get("/seeds/{id}", s.showSeed)
	r.Get("/seeds/new", s.newSeed)
	r.Post("/seeds/new", s.newSeed)
	r.Get("/seeds/{id}/edit", s.editSeed)
	r.Post("/seeds/{id}/edit", s.editSeed)
}

func (s *seeds) listSeeds(w http.ResponseWriter, r *http.Request) {
	seeds, err := s.db.ListSeeds(r.Context())
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	render(w, r, views.Seeds(seeds))
}

func (s *seeds) showSeed(w http.ResponseWriter, r *http.Request) {
	id, err := htu.Int64Param(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	seed, err := s.db.GetSeed(r.Context(), id)
	if err != nil {
		http.Error(w, "seed not found", http.StatusNotFound)
		return
	}
	render(w, r, views.Seed(seed))
}

func (s *seeds) newSeed(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		form := forms.New[db.CreateSeedParams]()
		render(w, r, views.NewSeed(form))
		return
	}

	csp, form, err := forms.FromRequest(&db.CreateSeedParams{}, r)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
	}
	if csp.Name == "" {
		form.AddError("Name", "required")
	}
	if form.HasErrors() {
		render(w, r, views.NewSeed(form))
		return
	}

	seed, err := s.db.CreateSeed(r.Context(), *csp)
	if err != nil {
		form.AddFormError("Internal error, please try again")
		render(w, r, views.NewSeed(form))
		return
	}

	w.Header().Set("HX-Replace-Url", fmt.Sprintf("/seeds/%d", seed.ID))
	render(w, r, views.Seed(seed))
}

func (s *seeds) editSeed(w http.ResponseWriter, r *http.Request) {
	id, err := htu.Int64Param(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		seed, err := s.db.GetSeed(r.Context(), id)
		if err != nil {
			http.Error(w, "seed not found", http.StatusNotFound)
			return
		}
		form := forms.FromStruct(&seed)
		render(w, r, views.EditSeed(seed.ID, form))
		return
	}

	csp, form, err := forms.FromRequest(&db.UpdateSeedParams{ID: id}, r)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
	}
	if csp.Name == "" {
		form.AddError("Name", "required")
	}
	if form.HasErrors() {
		render(w, r, views.EditSeed(id, form))
		return
	}

	seed, err := s.db.UpdateSeed(r.Context(), *csp)
	if err != nil {
		form.AddFormError("Internal error, please try again")
		render(w, r, views.EditSeed(id, form))
		return
	}

	w.Header().Set("Hx-Replace-Url", fmt.Sprintf("/seeds/%d", seed.ID))
	render(w, r, views.Seed(seed))
}
