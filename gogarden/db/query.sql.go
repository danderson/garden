// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package db

import (
	"context"

	"go.universe.tf/garden/gogarden/types"
)

const listLocations = `-- name: ListLocations :many
select id, name, inserted_at, updated_at, qr_id, qr_state from locations order by name
`

func (q *Queries) ListLocations(ctx context.Context) ([]Location, error) {
	rows, err := q.db.QueryContext(ctx, listLocations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Location
	for rows.Next() {
		var i Location
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.InsertedAt,
			&i.UpdatedAt,
			&i.QrID,
			&i.QRState,
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
select id, name, inserted_at, updated_at, front_image_id, back_image_id, year, edible, needs_trellis, needs_bird_netting, is_keto, is_native, is_invasive, is_cover_crop, grows_well_from_seed, is_bad_for_cats, is_deer_resistant, type, lifespan, family from seeds order by name
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

const updateLocation = `-- name: UpdateLocation :exec
update locations set name=?,qr_state=? where id=?
`

type UpdateLocationParams struct {
	Name    *string       `json:"name"`
	QRState types.QRState `json:"qr_state"`
	ID      int64         `json:"id"`
}

func (q *Queries) UpdateLocation(ctx context.Context, arg UpdateLocationParams) error {
	_, err := q.db.ExecContext(ctx, updateLocation, arg.Name, arg.QRState, arg.ID)
	return err
}
