package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"go.universe.tf/garden/db"
	"go.universe.tf/garden/forms"
	"go.universe.tf/garden/htu"
	"go.universe.tf/garden/types"
	"go.universe.tf/garden/views"
)

type plants struct {
	db *db.DB
}

func Plants(r *chi.Mux, db *db.DB) {
	s := plants{db}
	r.Get("/plants", chiFn(s.listPlants))
	r.Get("/plants/new", chiFn(s.newPlant))
	r.Post("/plants/new", chiFn(s.newPlant))
	r.Get("/plants/{id}", chiFn(s.showPlant))
	r.Get("/plants/{id}/edit", chiFn(s.editPlant))
	r.Post("/plants/{id}/edit", chiFn(s.editPlant))
	r.Get("/plants/search", chiFn(s.searchPlant))
	r.Get("/plants/{id}/uproot", chiFn(s.uprootPlant))
	r.Post("/plants/{id}/uproot", chiFn(s.uprootPlant))
}

func (s *plants) listPlants(w http.ResponseWriter, r *http.Request) error {
	plants, err := s.db.SearchPlants(r.Context(), "%")
	if err != nil {
		return internalErrorf("listing plants: %w", err)
	}
	render(w, r, views.Plants(plants))
	return nil
}

func (s *plants) showPlant(w http.ResponseWriter, r *http.Request) error {
	id, err := htu.Int64Param(r, "id")
	if err != nil {
		return badRequest(err)
	}
	return s.showPlantByID(w, r, id)
}

func (s *plants) showPlantByID(w http.ResponseWriter, r *http.Request, id int64) error {
	plant, err := s.db.GetPlant(r.Context(), id)
	if err != nil {
		return dbGetErrorf("getting plant: %w", err)
	}
	locs, err := s.db.GetPlantLocations(r.Context(), id)
	if err != nil {
		return dbGetErrorf("getting plant locations: %w", err)
	}

	render(w, r, views.Plant(plant, locs))
	return nil
}

func (s *plants) selectors(ctx context.Context) (seeds []forms.SelectOption, locations []forms.SelectOption, err error) {
	seedData, err := s.db.ListSeedsForSelector(ctx)
	if err != nil {
		return nil, nil, err
	}
	seeds = make([]forms.SelectOption, 0, len(seedData)+1)
	seeds = append(seeds, forms.SelectOption{
		Value: "",
		Label: "(none)",
	})
	for _, s := range seedData {
		seeds = append(seeds, forms.SelectOption{
			Value: fmt.Sprint(s.ID),
			Label: s.Name,
		})
	}

	locationData, err := s.db.ListLocationsForSelector(ctx)
	if err != nil {
		return nil, nil, err
	}
	locations = make([]forms.SelectOption, 0, len(locationData)+1)
	locations = append(locations, forms.SelectOption{
		Value: "",
		Label: "",
	})
	for _, s := range locationData {
		locations = append(locations, forms.SelectOption{
			Value: fmt.Sprint(s.ID),
			Label: s.Name,
		})
	}

	return seeds, locations, nil
}

func (s *plants) newPlantFormSelectors(ctx context.Context, form *forms.Form) error {
	seeds, locations, err := s.selectors(ctx)
	if err != nil {
		return err
	}
	form.SetSelectOptions("SeedID", seeds)
	form.SetSelectOptions("LocationID", locations)
	return nil
}

func (s *plants) newPlant(w http.ResponseWriter, r *http.Request) error {
	type createParams struct {
		Name       string
		LocationID int64
		Date       types.TextDate
		SeedID     *int64
	}

	if r.Method == "GET" {
		_, form, err := forms.FromRequest(&createParams{
			Date: types.TextDate{Time: time.Now()},
		}, r)
		if err != nil {
			return internalErrorf("parsing form: %w", err)
		}
		if err := s.newPlantFormSelectors(r.Context(), form); err != nil {
			return internalErrorf("adding form selectors: %w", err)
		}
		render(w, r, views.NewPlant(form))
		return nil
	}

	np, form, err := forms.FromRequest(&createParams{}, r)
	if err != nil {
		return internalErrorf("parsing form: %w", err)
	}
	if np.LocationID == 0 {
		form.AddError("LocationID", "required")
	}
	if np.Date.IsZero() {
		form.AddError("Date", "required")
	}
	if np.Date.After(time.Now()) {
		form.AddError("Date", "can't be in the future")
	}
	if np.SeedID == nil && np.Name == "" {
		form.AddFormError("One of seed or name is required")
	}
	if form.HasErrors() {
		if err := s.newPlantFormSelectors(r.Context(), form); err != nil {
			return internalErrorf("adding form selectors: %w", err)
		}
		render(w, r, views.NewPlant(form))
		return nil
	}

	tx, err := s.db.Tx(r.Context())
	if err != nil {
		return internalErrorf("starting transaction: %w", err)
	}
	defer tx.Rollback()

	nameFromSeed := int64(0)
	if np.Name == "" {
		seed, err := tx.GetSeed(r.Context(), *np.SeedID)
		if err != nil {
			return dbGetErrorf("getting seed: %w", err)
		}
		np.Name = seed.Name
		nameFromSeed = 1
	}

	p, err := s.db.CreatePlant(r.Context(), db.CreatePlantParams{
		Name:         np.Name,
		SeedID:       np.SeedID,
		NameFromSeed: nameFromSeed,
	})
	if err != nil {
		return internalErrorf("creating plant: %w", err)
	}
	_, err = s.db.CreatePlantLocation(r.Context(), db.CreatePlantLocationParams{
		PlantID:    p.ID,
		LocationID: np.LocationID,
		Start:      types.TextTime{Time: np.Date.Time},
	})
	if err != nil {
		return internalErrorf("creating plant location: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return internalErrorf("commit: %w", err)
	}

	w.Header().Set("HX-Replace-Url", fmt.Sprintf("/plants/%d", p.ID))
	return s.showPlantByID(w, r, p.ID)
}

func (s *plants) editPlantFormSelectors(ctx context.Context, form *forms.Form) error {
	seeds, locations, err := s.selectors(ctx)
	if err != nil {
		return err
	}
	form.SetSelectOptions("SeedID", seeds)
	form.SetSelectOptions("LocationID", locations)
	return nil
}

func (s *plants) editPlant(w http.ResponseWriter, r *http.Request) error {
	type editForm struct {
		SeedID     *int64
		LocationID int64
		Name       string
		Date       types.TextDate
	}

	id, err := htu.Int64Param(r, "id")
	if err != nil {
		return badRequest(err)
	}

	if r.Method == "GET" {
		plant, err := s.db.GetPlantForUpdate(r.Context(), id)
		if err != nil {
			return dbGetErrorf("getting plant: %w", err)
		}
		form := forms.FromStruct(&editForm{
			SeedID:     plant.SeedID,
			LocationID: plant.LocationID,
			Name:       plant.Name,
			Date:       types.TextDate{time.Now()},
		})
		if err := s.editPlantFormSelectors(r.Context(), form); err != nil {
			return internalErrorf("adding form selectors: %w", err)
		}
		render(w, r, views.EditPlant(id, form, false))
		return nil
	}

	tx, err := s.db.Tx(r.Context())
	if err != nil {
		return internalErrorf("starting transaction: %w", err)
	}
	defer tx.Rollback()

	curLocation, err := tx.GetPlantCurrentLocation(r.Context(), id)
	if err != nil {
		return internalErrorf("getting plant current location: %w", err)
	}

	params, form, err := forms.FromRequest(&editForm{}, r)
	if err != nil {
		return internalErrorf("parsing form: %w", err)
	}
	if params.LocationID == 0 {
		form.AddError("LocationID", "required")
	}
	locationChanged := params.LocationID != curLocation.LocationID
	if locationChanged {
		if params.Date.IsZero() {
			form.AddError("Date", "required")
		} else if params.Date.After(time.Now()) {
			form.AddError("Date", "can't be in the future")
		} else if params.Date.Before(curLocation.Start.Time) {
			form.AddError("Date", "can't be before last transplant")
		}
	}
	if params.SeedID == nil && params.Name == "" {
		form.AddFormError("One of seed or name is required")
	}
	if form.HasErrors() {
		if err := s.editPlantFormSelectors(r.Context(), form); err != nil {
			return internalErrorf("adding form selectors: %w", err)
		}
		render(w, r, views.EditPlant(id, form, locationChanged))
		return nil
	}

	up := db.UpdatePlantParams{
		ID:     id,
		Name:   params.Name,
		SeedID: params.SeedID,
	}

	if up.Name == "" {
		up.Name, err = tx.GetSeedName(r.Context(), *up.SeedID)
		if err != nil {
			return internalErrorf("getting seed name: %w", err)
		}
		up.NameFromSeed = 1
	} else {
		up.NameFromSeed = 0
	}

	plant, err := tx.UpdatePlant(r.Context(), up)
	if err != nil {
		return internalErrorf("updating plant: %w", err)
	}

	if locationChanged {
		transplant := types.TextTime{Time: params.Date.Time}
		err := tx.UprootPlant(r.Context(), db.UprootPlantParams{
			PlantID: id,
			End:     transplant,
		})
		if err != nil {
			return internalErrorf("removing plant from current location: %w", err)
		}

		_, err = tx.CreatePlantLocation(r.Context(), db.CreatePlantLocationParams{
			PlantID:    id,
			LocationID: params.LocationID,
			Start:      transplant,
		})
		if err != nil {
			return internalErrorf("setting new plant location: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return internalErrorf("committing transaction: %w", err)
	}

	w.Header().Set("HX-Replace-Url", fmt.Sprintf("/plants/%d", plant.ID))
	return s.showPlantByID(w, r, plant.ID)
}

func (s *plants) searchPlant(w http.ResponseWriter, r *http.Request) error {
	q := strings.Trim(r.FormValue("q"), "%")
	q = fmt.Sprintf("%%%s%%", q)
	plants, err := s.db.SearchPlants(r.Context(), q)
	if err != nil {
		return internalErrorf("executing search: %w", err)
	}
	views.PlantList(plants).Render(r.Context(), w)
	return nil
}

func (s *plants) uprootPlant(w http.ResponseWriter, r *http.Request) error {
	type uprootForm struct {
		End types.TextDate
	}

	id, err := htu.Int64Param(r, "id")
	if err != nil {
		return badRequest(err)
	}

	if r.Method == "GET" {
		form := forms.FromStruct(&uprootForm{types.TextDate{Time: time.Now()}})
		views.UprootPlantForm(id, form).Render(r.Context(), w)
		return nil
	}

	f, form, err := forms.FromRequest(&uprootForm{}, r)
	if err != nil {
		return internalErrorf("parsing form: %w", err)
	}
	if f.End.IsZero() {
		form.AddError("End", "required")
	}
	if f.End.After(time.Now()) {
		form.AddError("End", "can't be in the future")
	}
	if form.HasErrors() {
		views.UprootPlantForm(id, form).Render(r.Context(), w)
		return nil
	}

	tx, err := s.db.Tx(r.Context())
	if err != nil {
		return internalErrorf("start transaction: %w", err)
	}
	defer tx.Rollback()

	loc, err := tx.GetPlantCurrentLocation(r.Context(), id)
	if err != nil {
		return internalErrorf("getting plant location: %w", err)
	}

	err = tx.UprootPlant(r.Context(), db.UprootPlantParams{
		PlantID: id,
		End:     types.TextTime{Time: f.End.Time},
	})
	if err != nil {
		return internalErrorf("uprooting plant: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return internalErrorf("commit: %w", err)
	}

	w.Header().Set("Hx-Location", fmt.Sprintf("/locations/%d", loc.LocationID))
	w.WriteHeader(http.StatusOK)
	return nil
}
