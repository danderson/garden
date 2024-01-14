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
	r.Get("/seeds", chiFn(s.listSeeds))
	r.Get("/seeds/{id}", chiFn(s.showSeed))
	r.Get("/seeds/new", chiFn(s.newSeed))
	r.Post("/seeds/new", chiFn(s.newSeed))
	r.Get("/seeds/{id}/edit", chiFn(s.editSeed))
	r.Post("/seeds/{id}/edit", chiFn(s.editSeed))
}

func (s *seeds) listSeeds(w http.ResponseWriter, r *http.Request) error {
	seeds, err := s.db.ListSeeds(r.Context())
	if err != nil {
		return internalErrorf("listing seeds: %w", err)
	}
	render(w, r, views.Seeds(seeds))
	return nil
}

func (s *seeds) showSeed(w http.ResponseWriter, r *http.Request) error {
	id, err := htu.Int64Param(r, "id")
	if err != nil {
		return badRequest(err)
	}
	seed, err := s.db.GetSeed(r.Context(), id)
	if err != nil {
		return dbGetErrorf("getting seed: %w", err)
	}
	render(w, r, views.Seed(seed))
	return nil
}

func (s *seeds) newSeed(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		form := forms.New[db.CreateSeedParams]()
		render(w, r, views.NewSeed(form))
		return nil
	}

	csp, form, err := forms.FromRequest(&db.CreateSeedParams{}, r)
	if err != nil {
		return internalErrorf("parsing form: %w", err)
	}
	if csp.Name == "" {
		form.AddError("Name", "required")
	}
	if form.HasErrors() {
		render(w, r, views.NewSeed(form))
		return nil
	}

	seed, err := s.db.CreateSeed(r.Context(), *csp)
	if err != nil {
		return internalErrorf("creating seed: %w", err)
	}

	w.Header().Set("HX-Replace-Url", fmt.Sprintf("/seeds/%d", seed.ID))
	render(w, r, views.Seed(seed))
	return nil
}

func (s *seeds) editSeed(w http.ResponseWriter, r *http.Request) error {
	id, err := htu.Int64Param(r, "id")
	if err != nil {
		return badRequest(err)
	}
	if r.Method == "GET" {
		seed, err := s.db.GetSeed(r.Context(), id)
		if err != nil {
			return dbGetErrorf("getting seed: %w", err)
		}
		form := forms.FromStruct(&seed)
		render(w, r, views.EditSeed(seed.ID, form))
		return nil
	}

	csp, form, err := forms.FromRequest(&db.UpdateSeedParams{ID: id}, r)
	if err != nil {
		return internalErrorf("parsing form: %w", err)
	}
	if csp.Name == "" {
		form.AddError("Name", "required")
	}
	if form.HasErrors() {
		render(w, r, views.EditSeed(id, form))
		return nil
	}

	seed, err := s.db.UpdateSeed(r.Context(), *csp)
	if err != nil {
		return internalErrorf("updating seed: %w", err)
	}

	w.Header().Set("Hx-Replace-Url", fmt.Sprintf("/seeds/%d", seed.ID))
	render(w, r, views.Seed(seed))
	return nil
}
