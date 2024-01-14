package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.universe.tf/garden/gogarden/db"
	"go.universe.tf/garden/gogarden/forms"
	"go.universe.tf/garden/gogarden/htu"
	"go.universe.tf/garden/gogarden/types"
	"go.universe.tf/garden/gogarden/views"
)

type plants struct {
	db *db.DB
}

func Plants(r *chi.Mux, db *db.DB) {
	s := plants{db}
	r.Get("/plants", s.listPlants)
	r.Get("/plants/new", s.newPlant)
	r.Post("/plants/new", s.newPlant)
	r.Get("/plants/{id}", s.showPlant)
	r.Get("/plants/{id}/edit", s.editPlant)
	r.Post("/plants/{id}/edit", s.editPlant)
}

func (s *plants) listPlants(w http.ResponseWriter, r *http.Request) {
	plants, err := s.db.ListPlants(r.Context())
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	render(w, r, views.Plants(plants))
}

func (s *plants) showPlant(w http.ResponseWriter, r *http.Request) {
	id, err := htu.Int64Param(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	plant, err := s.db.GetPlant(r.Context(), id)
	if err != nil {
		http.Error(w, "plant not found", http.StatusNotFound)
		return
	}

	render(w, r, views.Plant(plant))
}

func (s *plants) selectors(ctx context.Context) (seeds []forms.SelectOption, locations []forms.SelectOption, err error) {
	seedData, err := s.db.ListSeedsForSelector(ctx)
	if err != nil {
		return nil, nil, err
	}
	seeds = make([]forms.SelectOption, 0, len(seedData))
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
	locations = make([]forms.SelectOption, 0, len(locationData))
	for _, s := range locationData {
		locations = append(locations, forms.SelectOption{
			Value: fmt.Sprint(s.ID),
			Label: s.Name,
		})
	}

	return seeds, locations, nil
}

func (s *plants) formSelectors(ctx context.Context, form *forms.Form) error {
	seeds, locations, err := s.selectors(ctx)
	if err != nil {
		return err
	}
	form.SetSelectOptions("SeedID", seeds)
	form.SetSelectOptions("LocationID", locations)
	return nil
}

func (s *plants) newPlant(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		form := forms.New[db.GetPlantForUpdateRow]()
		if err := s.formSelectors(r.Context(), form); err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
		render(w, r, views.NewPlant(form))
		return
	}

	np, form, err := forms.FromRequest(&db.GetPlantForUpdateRow{}, r)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	if np.LocationID == 0 {
		form.AddError("LocationID", "required")
	}
	if np.SeedID == nil && np.Name == "" {
		form.AddFormError("One of seed or name is required")
	}
	if form.HasErrors() {
		if err := s.formSelectors(r.Context(), form); err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		render(w, r, views.NewPlant(form))
		return
	}

	tx, err := s.db.Tx(r.Context())
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	nameFromSeed := int64(0)
	if np.Name == "" {
		seed, err := tx.GetSeed(r.Context(), *np.SeedID)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
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
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	_, err = s.db.CreatePlantLocation(r.Context(), db.CreatePlantLocationParams{
		PlantID:    p.ID,
		LocationID: np.LocationID,
		Start:      types.TextTime{Time: time.Now()},
	})
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Replace-Url", fmt.Sprintf("/plants/%d", p.ID))
	render(w, r, views.Plant(p))
}

func (s *plants) editPlant(w http.ResponseWriter, r *http.Request) {
	// id, err := htu.Int64Param(r, "id")
	// if err != nil {
	// 	http.Error(w, "invalid id", http.StatusBadRequest)
	// 	return
	// }

	// if r.Method == "GET" {
	// 	plant, err := s.db.GetPlant(r.Context(), id)
	// 	if err != nil {
	// 		http.Error(w, "seed not found", http.StatusNotFound)
	// 		return
	// 	}
	// 	form := forms.FromStruct(&plant)
	// 	if err := s.formSelectors(r.Context(), form); err != nil {
	// 		http.Error(w, "database error", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	render(w, r, views.EditPlant(plant.ID, form))
	// 	return
	// }

	// ulp, form, err := forms.FromRequest(&db.UpdatePlantParams{ID: id}, r)
	// if err != nil {
	// 	http.Error(w, "invalid input", http.StatusBadRequest)
	// 	return
	// }
	// if ulp.Name == "" {
	// 	form.AddError("Name", "required")
	// }
	// if form.HasErrors() {
	// 	if err := s.formSelectors(r.Context(), form); err != nil {
	// 		http.Error(w, "database error", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	render(w, r, views.EditPlant(id, form))
	// 	return
	// }

	// plant, err := s.db.UpdatePlant(r.Context(), *ulp)
	// if err != nil {
	// 	form.AddFormError("Internal error, please try again")
	// 	if err := s.formSelectors(r.Context(), form); err != nil {
	// 		http.Error(w, "database error", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	render(w, r, views.NewPlant(form))
	// 	return
	// }

	// w.Header().Set("HX-Replace-Url", fmt.Sprintf("/plants/%d", plant.ID))
	// render(w, r, views.Plant(plant))
}
