// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"go.universe.tf/garden/gogarden/types"
)

type Location struct {
	ID         int64          `json:"id"`
	Name       *string        `json:"name"`
	InsertedAt types.TextTime `json:"inserted_at"`
	UpdatedAt  types.TextTime `json:"updated_at"`
	QrID       string         `json:"qr_id"`
	QRState    types.QRState  `json:"qr_state"`
}

type LocationsImage struct {
	ID         int64   `json:"id"`
	ImageID    *string `json:"image_id"`
	LocationID *int64  `json:"location_id"`
	InsertedAt string  `json:"inserted_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type Plant struct {
	ID           int64   `json:"id"`
	Name         *string `json:"name"`
	SeedID       *int64  `json:"seed_id"`
	InsertedAt   string  `json:"inserted_at"`
	UpdatedAt    string  `json:"updated_at"`
	NameFromSeed *int64  `json:"name_from_seed"`
}

type PlantLocation struct {
	ID         int64       `json:"id"`
	PlantID    *int64      `json:"plant_id"`
	LocationID *int64      `json:"location_id"`
	Start      string      `json:"start"`
	End        interface{} `json:"end"`
}

type SchemaMigration struct {
	Version    int64   `json:"version"`
	InsertedAt *string `json:"inserted_at"`
}

type Seed struct {
	ID                int64   `json:"id"`
	Name              *string `json:"name"`
	InsertedAt        string  `json:"inserted_at"`
	UpdatedAt         string  `json:"updated_at"`
	FrontImageID      *string `json:"front_image_id"`
	BackImageID       *string `json:"back_image_id"`
	Year              *int64  `json:"year"`
	Edible            *int64  `json:"edible"`
	NeedsTrellis      *int64  `json:"needs_trellis"`
	NeedsBirdNetting  *int64  `json:"needs_bird_netting"`
	IsKeto            *int64  `json:"is_keto"`
	IsNative          *int64  `json:"is_native"`
	IsInvasive        *int64  `json:"is_invasive"`
	IsCoverCrop       *int64  `json:"is_cover_crop"`
	GrowsWellFromSeed *int64  `json:"grows_well_from_seed"`
	IsBadForCats      *int64  `json:"is_bad_for_cats"`
	IsDeerResistant   *int64  `json:"is_deer_resistant"`
	Type              *string `json:"type"`
	Lifespan          *string `json:"lifespan"`
	Family            *string `json:"family"`
}
