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
	r.Get("/locations", s.listLocations)
	r.Get("/locations/new", s.newLocation)
	r.Post("/locations/new", s.newLocation)
	r.Get("/locations/{id}", s.showLocation)
	r.Get("/locations/{id}/edit", s.editLocation)
	r.Post("/locations/{id}/edit", s.editLocation)
}

func (s *locations) listLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := s.db.ListLocations(r.Context())
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	render(w, r, views.Locations(locations))
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

func (s *locations) showLocation(w http.ResponseWriter, r *http.Request) {
	id, err := htu.Int64Param(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	location, err := s.db.GetLocation(r.Context(), id)
	if err != nil {
		http.Error(w, "location not found", http.StatusNotFound)
		return
	}
	current, old, err := s.plantsInLocation(r.Context(), id)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	render(w, r, views.Location(location, current, old))
}

func (s *locations) newLocation(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		form := forms.New[db.CreateLocationParams]()
		render(w, r, views.NewLocation(form))
		return
	}

	csp, form, err := forms.FromRequest(&db.CreateLocationParams{}, r)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	if csp.Name == "" {
		form.AddError("Name", "required")
	}
	if form.HasErrors() {
		render(w, r, views.NewLocation(form))
		return
	}

	location, err := s.db.CreateLocation(r.Context(), *csp)
	if err != nil {
		form.AddFormError("Internal error, please try again")
		render(w, r, views.NewLocation(form))
		return
	}
	current, old, err := s.plantsInLocation(r.Context(), location.ID)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Replace-Url", fmt.Sprintf("/locations/%d", location.ID))
	render(w, r, views.Location(location, current, old))
}

func (s *locations) editLocation(w http.ResponseWriter, r *http.Request) {
	id, err := htu.Int64Param(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		location, err := s.db.GetLocation(r.Context(), id)
		if err != nil {
			http.Error(w, "seed not found", http.StatusNotFound)
			return
		}
		form := forms.FromStruct(&location)
		render(w, r, views.EditLocation(location.ID, form))
		return
	}

	ulp, form, err := forms.FromRequest(&db.UpdateLocationParams{ID: id}, r)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	if ulp.Name == "" {
		form.AddError("Name", "required")
	}
	if form.HasErrors() {
		render(w, r, views.EditLocation(id, form))
		return
	}

	location, err := s.db.UpdateLocation(r.Context(), *ulp)
	if err != nil {
		form.AddFormError("Internal error, please try again")
		render(w, r, views.NewLocation(form))
		return
	}
	current, old, err := s.plantsInLocation(r.Context(), location.ID)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Replace-Url", fmt.Sprintf("/locations/%d", location.ID))
	render(w, r, views.Location(location, current, old))
}
