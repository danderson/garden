package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.universe.tf/garden/gogarden/db"
	"go.universe.tf/garden/gogarden/forms"
	"go.universe.tf/garden/gogarden/htu"
	"go.universe.tf/garden/gogarden/views"
)

type locations struct {
	db *db.DB
}

func Locations(r *chi.Mux, db *db.DB) {
	s := locations{db}
	r.Get("/locations", chiFn(s.listLocations))
	r.Get("/locations/new", chiFn(s.newLocation))
	r.Post("/locations/new", chiFn(s.newLocation))
	r.Get("/locations/{id}", chiFn(s.showLocation))
	r.Get("/locations/{id}/edit", chiFn(s.editLocation))
	r.Post("/locations/{id}/edit", chiFn(s.editLocation))
}

func (s *locations) listLocations(w http.ResponseWriter, r *http.Request) error {
	locations, err := s.db.ListLocations(r.Context())
	if err != nil {
		return internalErrorf("listing locations: %w", err)
	}
	render(w, r, views.Locations(locations))
	return nil
}

func (s *locations) plantsInLocation(ctx context.Context, id int64) (current, old []db.GetPlantsInLocationRow, err error) {
	plants, err := s.db.GetPlantsInLocation(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	for _, p := range plants {
		if p.End.IsZero() {
			current = append(current, p)
		} else {
			old = append(old, p)
		}
	}
	return current, old, nil
}

func (s *locations) showLocation(w http.ResponseWriter, r *http.Request) error {
	id, err := htu.Int64Param(r, "id")
	if err != nil {
		return badRequest(err)
	}
	location, err := s.db.GetLocation(r.Context(), id)
	if err != nil {
		return dbGetErrorf("getting location: %w", err)
	}
	current, old, err := s.plantsInLocation(r.Context(), id)
	if err != nil {
		return dbGetErrorf("getting plants in location: %w", err)
	}

	render(w, r, views.Location(location, current, old))
	return nil
}

func (s *locations) newLocation(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		form := forms.New[db.CreateLocationParams]()
		render(w, r, views.NewLocation(form))
		return nil
	}

	csp, form, err := forms.FromRequest(&db.CreateLocationParams{}, r)
	if err != nil {
		return internalErrorf("parsing form: %w", err)
	}
	if csp.Name == "" {
		form.AddError("Name", "required")
	}
	if form.HasErrors() {
		render(w, r, views.NewLocation(form))
		return nil
	}

	location, err := s.db.CreateLocation(r.Context(), *csp)
	if err != nil {
		return internalErrorf("creating location: %w", err)
	}
	current, old, err := s.plantsInLocation(r.Context(), location.ID)
	if err != nil {
		return internalErrorf("getting plants in location: %w", err)
	}

	w.Header().Set("HX-Replace-Url", fmt.Sprintf("/locations/%d", location.ID))
	render(w, r, views.Location(location, current, old))
	return nil
}

func (s *locations) editLocation(w http.ResponseWriter, r *http.Request) error {
	id, err := htu.Int64Param(r, "id")
	if err != nil {
		return badRequest(err)
	}

	if r.Method == "GET" {
		location, err := s.db.GetLocation(r.Context(), id)
		if err != nil {
			return dbGetErrorf("getting location: %w", err)
		}
		form := forms.FromStruct(&location)
		render(w, r, views.EditLocation(location.ID, form))
		return nil
	}

	ulp, form, err := forms.FromRequest(&db.UpdateLocationParams{ID: id}, r)
	if err != nil {
		return internalErrorf("parsing form: %w", err)
	}
	if ulp.Name == "" {
		form.AddError("Name", "required")
	}
	if form.HasErrors() {
		render(w, r, views.EditLocation(id, form))
		return nil
	}

	location, err := s.db.UpdateLocation(r.Context(), *ulp)
	if err != nil {
		return internalErrorf("updating location: %w", err)
	}
	current, old, err := s.plantsInLocation(r.Context(), location.ID)
	if err != nil {
		return internalErrorf("getting plants in location: %w", err)
	}

	w.Header().Set("HX-Replace-Url", fmt.Sprintf("/locations/%d", location.ID))
	render(w, r, views.Location(location, current, old))
	return nil
}
