// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package db

import (
	"context"

	"go.universe.tf/garden/gogarden/types"
	"go.universe.tf/garden/gogarden/types/plantfamily"
	"go.universe.tf/garden/gogarden/types/tribool"
)

const createLocation = `-- name: CreateLocation :one
insert into locations (
  name,
  qr_id,
  qr_state,
  inserted_at,
  updated_at) values (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) returning id, name, inserted_at, updated_at, qr_id, qr_state
`

type CreateLocationParams struct {
	Name    string        `json:"name"`
	QRID    string        `json:"qr_id"`
	QRState types.QRState `json:"qr_state"`
}

func (q *Queries) CreateLocation(ctx context.Context, arg CreateLocationParams) (Location, error) {
	row := q.db.QueryRowContext(ctx, createLocation, arg.Name, arg.QRID, arg.QRState)
	var i Location
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.InsertedAt,
		&i.UpdatedAt,
		&i.QRID,
		&i.QRState,
	)
	return i, err
}

const createPlant = `-- name: CreatePlant :one
insert into plants (
  name,
  seed_id,
  name_from_seed,
  inserted_at,
  updated_at) values (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) returning id, name, seed_id, inserted_at, updated_at, name_from_seed
`

type CreatePlantParams struct {
	Name         string `json:"name"`
	SeedID       *int64 `json:"seed_id"`
	NameFromSeed int64  `json:"name_from_seed"`
}

func (q *Queries) CreatePlant(ctx context.Context, arg CreatePlantParams) (Plant, error) {
	row := q.db.QueryRowContext(ctx, createPlant, arg.Name, arg.SeedID, arg.NameFromSeed)
	var i Plant
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SeedID,
		&i.InsertedAt,
		&i.UpdatedAt,
		&i.NameFromSeed,
	)
	return i, err
}

const createPlantLocation = `-- name: CreatePlantLocation :one
insert into plant_locations (
  plant_id,
  location_id,
  start) values (?,?,?) returning id, plant_id, location_id, start, "end"
`

type CreatePlantLocationParams struct {
	PlantID    int64          `json:"plant_id"`
	LocationID int64          `json:"location_id"`
	Start      types.TextTime `json:"start"`
}

func (q *Queries) CreatePlantLocation(ctx context.Context, arg CreatePlantLocationParams) (PlantLocation, error) {
	row := q.db.QueryRowContext(ctx, createPlantLocation, arg.PlantID, arg.LocationID, arg.Start)
	var i PlantLocation
	err := row.Scan(
		&i.ID,
		&i.PlantID,
		&i.LocationID,
		&i.Start,
		&i.End,
	)
	return i, err
}

const createSeed = `-- name: CreateSeed :one
insert into seeds (
  name,
  family,
  inserted_at,
  updated_at,
  year,
  edible,
  needs_trellis,
  needs_bird_netting,
  is_keto,
  is_native,
  is_invasive,
  is_cover_crop,
  grows_well_from_seed,
  is_bad_for_cats,
  is_deer_resistant) values (?,?,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,?,?,?,?,?,?,?,?,?,?,?) returning id, name, inserted_at, updated_at, front_image_id, back_image_id, year, edible, needs_trellis, needs_bird_netting, is_keto, is_native, is_invasive, is_cover_crop, grows_well_from_seed, is_bad_for_cats, is_deer_resistant, type, lifespan, family
`

type CreateSeedParams struct {
	Name              string                  `json:"name"`
	Family            plantfamily.PlantFamily `json:"family"`
	Year              *int64                  `json:"year"`
	Edible            tribool.Tribool         `json:"edible"`
	NeedsTrellis      tribool.Tribool         `json:"needs_trellis"`
	NeedsBirdNetting  tribool.Tribool         `json:"needs_bird_netting"`
	IsKeto            tribool.Tribool         `json:"is_keto"`
	IsNative          tribool.Tribool         `json:"is_native"`
	IsInvasive        tribool.Tribool         `json:"is_invasive"`
	IsCoverCrop       tribool.Tribool         `json:"is_cover_crop"`
	GrowsWellFromSeed tribool.Tribool         `json:"grows_well_from_seed"`
	IsBadForCats      tribool.Tribool         `json:"is_bad_for_cats"`
	IsDeerResistant   tribool.Tribool         `json:"is_deer_resistant"`
}

func (q *Queries) CreateSeed(ctx context.Context, arg CreateSeedParams) (Seed, error) {
	row := q.db.QueryRowContext(ctx, createSeed,
		arg.Name,
		arg.Family,
		arg.Year,
		arg.Edible,
		arg.NeedsTrellis,
		arg.NeedsBirdNetting,
		arg.IsKeto,
		arg.IsNative,
		arg.IsInvasive,
		arg.IsCoverCrop,
		arg.GrowsWellFromSeed,
		arg.IsBadForCats,
		arg.IsDeerResistant,
	)
	var i Seed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.InsertedAt,
		&i.UpdatedAt,
		&i.FrontImageID,
		&i.BackImageID,
		&i.Year,
		&i.Edible,
		&i.NeedsTrellis,
		&i.NeedsBirdNetting,
		&i.IsKeto,
		&i.IsNative,
		&i.IsInvasive,
		&i.IsCoverCrop,
		&i.GrowsWellFromSeed,
		&i.IsBadForCats,
		&i.IsDeerResistant,
		&i.Type,
		&i.Lifespan,
		&i.Family,
	)
	return i, err
}

const getLocation = `-- name: GetLocation :one
select id, name, inserted_at, updated_at, qr_id, qr_state from locations where id=?
`

func (q *Queries) GetLocation(ctx context.Context, id int64) (Location, error) {
	row := q.db.QueryRowContext(ctx, getLocation, id)
	var i Location
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.InsertedAt,
		&i.UpdatedAt,
		&i.QRID,
		&i.QRState,
	)
	return i, err
}

const getPlant = `-- name: GetPlant :one
select id, name, seed_id, inserted_at, updated_at, name_from_seed from plants where id=?
`

func (q *Queries) GetPlant(ctx context.Context, id int64) (Plant, error) {
	row := q.db.QueryRowContext(ctx, getPlant, id)
	var i Plant
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SeedID,
		&i.InsertedAt,
		&i.UpdatedAt,
		&i.NameFromSeed,
	)
	return i, err
}

const getPlantForUpdate = `-- name: GetPlantForUpdate :one
select pl.location_id,p.seed_id,p.name
  from plants as p inner join plant_locations as pl on pl.plant_id=p.id
 order by pl.start desc
 limit 1
`

type GetPlantForUpdateRow struct {
	LocationID int64  `json:"location_id"`
	SeedID     *int64 `json:"seed_id"`
	Name       string `json:"name"`
}

func (q *Queries) GetPlantForUpdate(ctx context.Context) (GetPlantForUpdateRow, error) {
	row := q.db.QueryRowContext(ctx, getPlantForUpdate)
	var i GetPlantForUpdateRow
	err := row.Scan(&i.LocationID, &i.SeedID, &i.Name)
	return i, err
}

const getPlantsInLocation = `-- name: GetPlantsInLocation :many
select p.name,pl.start,pl.end from locations as l
                                   inner join plant_locations as pl on l.id=pl.location_id
                                   inner join plants as p on p.id=pl.plant_id
 where l.id=?
 order by p.name collate nocase,
          pl.end desc,
          pl.start desc
`

type GetPlantsInLocationRow struct {
	Name  string         `json:"name"`
	Start types.TextTime `json:"start"`
	End   types.TextTime `json:"end"`
}

func (q *Queries) GetPlantsInLocation(ctx context.Context, id int64) ([]GetPlantsInLocationRow, error) {
	rows, err := q.db.QueryContext(ctx, getPlantsInLocation, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPlantsInLocationRow
	for rows.Next() {
		var i GetPlantsInLocationRow
		if err := rows.Scan(&i.Name, &i.Start, &i.End); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSeed = `-- name: GetSeed :one
select id, name, inserted_at, updated_at, front_image_id, back_image_id, year, edible, needs_trellis, needs_bird_netting, is_keto, is_native, is_invasive, is_cover_crop, grows_well_from_seed, is_bad_for_cats, is_deer_resistant, type, lifespan, family from seeds where id=?
`

func (q *Queries) GetSeed(ctx context.Context, id int64) (Seed, error) {
	row := q.db.QueryRowContext(ctx, getSeed, id)
	var i Seed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.InsertedAt,
		&i.UpdatedAt,
		&i.FrontImageID,
		&i.BackImageID,
		&i.Year,
		&i.Edible,
		&i.NeedsTrellis,
		&i.NeedsBirdNetting,
		&i.IsKeto,
		&i.IsNative,
		&i.IsInvasive,
		&i.IsCoverCrop,
		&i.GrowsWellFromSeed,
		&i.IsBadForCats,
		&i.IsDeerResistant,
		&i.Type,
		&i.Lifespan,
		&i.Family,
	)
	return i, err
}

const listLocations = `-- name: ListLocations :many
select l.id, l.name, l.inserted_at, l.updated_at, l.qr_id, l.qr_state,count(pl.id) as num_plants
  from locations as l left join plant_locations as pl on l.id=pl.location_id
  where pl.end is null
  group by l.id
  order by l.name collate nocase
`

type ListLocationsRow struct {
	ID         int64          `json:"id"`
	Name       string         `json:"name"`
	InsertedAt types.TextTime `json:"inserted_at"`
	UpdatedAt  types.TextTime `json:"updated_at"`
	QRID       string         `json:"qr_id"`
	QRState    types.QRState  `json:"qr_state"`
	NumPlants  int64          `json:"num_plants"`
}

func (q *Queries) ListLocations(ctx context.Context) ([]ListLocationsRow, error) {
	rows, err := q.db.QueryContext(ctx, listLocations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListLocationsRow
	for rows.Next() {
		var i ListLocationsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.InsertedAt,
			&i.UpdatedAt,
			&i.QRID,
			&i.QRState,
			&i.NumPlants,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listLocationsForSelector = `-- name: ListLocationsForSelector :many
select id,name from locations order by name collate nocase
`

type ListLocationsForSelectorRow struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) ListLocationsForSelector(ctx context.Context) ([]ListLocationsForSelectorRow, error) {
	rows, err := q.db.QueryContext(ctx, listLocationsForSelector)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListLocationsForSelectorRow
	for rows.Next() {
		var i ListLocationsForSelectorRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPlants = `-- name: ListPlants :many
select id, name, seed_id, inserted_at, updated_at, name_from_seed from plants order by name collate nocase
`

func (q *Queries) ListPlants(ctx context.Context) ([]Plant, error) {
	rows, err := q.db.QueryContext(ctx, listPlants)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Plant
	for rows.Next() {
		var i Plant
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.SeedID,
			&i.InsertedAt,
			&i.UpdatedAt,
			&i.NameFromSeed,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSeeds = `-- name: ListSeeds :many
select id, name, inserted_at, updated_at, front_image_id, back_image_id, year, edible, needs_trellis, needs_bird_netting, is_keto, is_native, is_invasive, is_cover_crop, grows_well_from_seed, is_bad_for_cats, is_deer_resistant, type, lifespan, family from seeds order by name collate nocase
`

func (q *Queries) ListSeeds(ctx context.Context) ([]Seed, error) {
	rows, err := q.db.QueryContext(ctx, listSeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Seed
	for rows.Next() {
		var i Seed
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.InsertedAt,
			&i.UpdatedAt,
			&i.FrontImageID,
			&i.BackImageID,
			&i.Year,
			&i.Edible,
			&i.NeedsTrellis,
			&i.NeedsBirdNetting,
			&i.IsKeto,
			&i.IsNative,
			&i.IsInvasive,
			&i.IsCoverCrop,
			&i.GrowsWellFromSeed,
			&i.IsBadForCats,
			&i.IsDeerResistant,
			&i.Type,
			&i.Lifespan,
			&i.Family,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSeedsForSelector = `-- name: ListSeedsForSelector :many
select id,name from seeds order by name collate nocase
`

type ListSeedsForSelectorRow struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) ListSeedsForSelector(ctx context.Context) ([]ListSeedsForSelectorRow, error) {
	rows, err := q.db.QueryContext(ctx, listSeedsForSelector)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListSeedsForSelectorRow
	for rows.Next() {
		var i ListSeedsForSelectorRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateLocation = `-- name: UpdateLocation :one
update locations set name=?,qr_id=?,qr_state=?,updated_at=CURRENT_TIMESTAMP where id=? returning id, name, inserted_at, updated_at, qr_id, qr_state
`

type UpdateLocationParams struct {
	Name    string        `json:"name"`
	QRID    string        `json:"qr_id"`
	QRState types.QRState `json:"qr_state"`
	ID      int64         `json:"id"`
}

func (q *Queries) UpdateLocation(ctx context.Context, arg UpdateLocationParams) (Location, error) {
	row := q.db.QueryRowContext(ctx, updateLocation,
		arg.Name,
		arg.QRID,
		arg.QRState,
		arg.ID,
	)
	var i Location
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.InsertedAt,
		&i.UpdatedAt,
		&i.QRID,
		&i.QRState,
	)
	return i, err
}

const updatePlant = `-- name: UpdatePlant :one
update plants set name=?,seed_id=?,name_from_seed=?,updated_at=CURRENT_TIMESTAMP where id=? returning id, name, seed_id, inserted_at, updated_at, name_from_seed
`

type UpdatePlantParams struct {
	Name         string `json:"name"`
	SeedID       *int64 `json:"seed_id"`
	NameFromSeed int64  `json:"name_from_seed"`
	ID           int64  `json:"id"`
}

func (q *Queries) UpdatePlant(ctx context.Context, arg UpdatePlantParams) (Plant, error) {
	row := q.db.QueryRowContext(ctx, updatePlant,
		arg.Name,
		arg.SeedID,
		arg.NameFromSeed,
		arg.ID,
	)
	var i Plant
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SeedID,
		&i.InsertedAt,
		&i.UpdatedAt,
		&i.NameFromSeed,
	)
	return i, err
}

const updateSeed = `-- name: UpdateSeed :one
update seeds set name=?,family=?,updated_at=CURRENT_TIMESTAMP,year=?,edible=?,needs_trellis=?,needs_bird_netting=?,is_keto=?,is_native=?,is_invasive=?,is_cover_crop=?,grows_well_from_seed=?,is_bad_for_cats=?,is_deer_resistant=? where id=? returning id, name, inserted_at, updated_at, front_image_id, back_image_id, year, edible, needs_trellis, needs_bird_netting, is_keto, is_native, is_invasive, is_cover_crop, grows_well_from_seed, is_bad_for_cats, is_deer_resistant, type, lifespan, family
`

type UpdateSeedParams struct {
	Name              string                  `json:"name"`
	Family            plantfamily.PlantFamily `json:"family"`
	Year              *int64                  `json:"year"`
	Edible            tribool.Tribool         `json:"edible"`
	NeedsTrellis      tribool.Tribool         `json:"needs_trellis"`
	NeedsBirdNetting  tribool.Tribool         `json:"needs_bird_netting"`
	IsKeto            tribool.Tribool         `json:"is_keto"`
	IsNative          tribool.Tribool         `json:"is_native"`
	IsInvasive        tribool.Tribool         `json:"is_invasive"`
	IsCoverCrop       tribool.Tribool         `json:"is_cover_crop"`
	GrowsWellFromSeed tribool.Tribool         `json:"grows_well_from_seed"`
	IsBadForCats      tribool.Tribool         `json:"is_bad_for_cats"`
	IsDeerResistant   tribool.Tribool         `json:"is_deer_resistant"`
	ID                int64                   `json:"id"`
}

func (q *Queries) UpdateSeed(ctx context.Context, arg UpdateSeedParams) (Seed, error) {
	row := q.db.QueryRowContext(ctx, updateSeed,
		arg.Name,
		arg.Family,
		arg.Year,
		arg.Edible,
		arg.NeedsTrellis,
		arg.NeedsBirdNetting,
		arg.IsKeto,
		arg.IsNative,
		arg.IsInvasive,
		arg.IsCoverCrop,
		arg.GrowsWellFromSeed,
		arg.IsBadForCats,
		arg.IsDeerResistant,
		arg.ID,
	)
	var i Seed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.InsertedAt,
		&i.UpdatedAt,
		&i.FrontImageID,
		&i.BackImageID,
		&i.Year,
		&i.Edible,
		&i.NeedsTrellis,
		&i.NeedsBirdNetting,
		&i.IsKeto,
		&i.IsNative,
		&i.IsInvasive,
		&i.IsCoverCrop,
		&i.GrowsWellFromSeed,
		&i.IsBadForCats,
		&i.IsDeerResistant,
		&i.Type,
		&i.Lifespan,
		&i.Family,
	)
	return i, err
}
