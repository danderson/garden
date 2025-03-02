package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.universe.tf/garden/db"
	"go.universe.tf/garden/forms"
	"go.universe.tf/garden/htu"
	"go.universe.tf/garden/types/plantfamily"
	"go.universe.tf/garden/views"
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
	r.Get("/seeds/search", chiFn(s.searchSeeds))
	r.Get("/seeds/search-complete", chiFn(s.searchSeedsAutocomplete))
	r.Get("/seeds/search-family", chiFn(s.searchFamilyAutocomplete))
}

func (s *seeds) listSeeds(w http.ResponseWriter, r *http.Request) error {
	seeds, err := s.db.SearchSeeds(r.Context(), "%")
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
	tx, err := s.db.ReadTx(r.Context())
	if err != nil {
		return dbGetErrorf("starting transaction: %w", err)
	}
	defer tx.Rollback()
	seed, err := tx.GetSeed(r.Context(), id)
	if err != nil {
		return dbGetErrorf("getting seed: %w", err)
	}
	rawHist, err := tx.GetSeedHistory(r.Context(), id)
	if err != nil {
		return dbGetErrorf("getting seed history: %w", err)
	}

	var hist [][]db.GetSeedHistoryRow
	lastID := int64(-1)
	for _, h := range rawHist {
		if h.ID != lastID {
			hist = append(hist, nil)
			lastID = h.ID
		}
		hist[len(hist)-1] = append(hist[len(hist)-1], h)
	}

	render(w, r, views.Seed(seed, hist))
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

	w.Header().Set("HX-Redirect", fmt.Sprintf("/seeds/%d", seed.ID))
	w.WriteHeader(http.StatusOK)
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

	w.Header().Set("HX-Redirect", fmt.Sprintf("/seeds/%d", seed.ID))
	w.WriteHeader(http.StatusOK)
	return nil
}

func (s *seeds) searchSeeds(w http.ResponseWriter, r *http.Request) error {
	q := strings.Trim(r.FormValue("q"), "%")
	q = fmt.Sprintf("%%%s%%", q)
	seeds, err := s.db.SearchSeeds(r.Context(), q)
	if err != nil {
		return internalErrorf("executing search: %w", err)
	}
	views.SeedList(seeds).Render(r.Context(), w)
	return nil
}

func (s *seeds) searchSeedsAutocomplete(w http.ResponseWriter, r *http.Request) error {
	q := strings.Trim(r.FormValue("q"), "%")
	q = fmt.Sprintf("%%%s%%", q)
	seeds, err := s.db.SearchSeeds(r.Context(), q)
	if err != nil {
		return internalErrorf("executing search: %w", err)
	}
	views.SeedListAutocomplete(seeds).Render(r.Context(), w)
	return nil
}

func (s *seeds) searchFamilyAutocomplete(w http.ResponseWriter, r *http.Request) error {
	q := strings.ToLower(r.FormValue("q"))
	var opts []forms.SelectOption
	for _, opt := range plantfamily.Unknown.SelectOptions() {
		if strings.Contains(strings.ToLower(opt.Label), q) {
			opts = append(opts, opt)
		}
	}
	views.ComboInputOptions(opts).Render(r.Context(), w)
	return nil
}
