-- name: ListLocations :many
select * from locations order by name collate nocase;

-- name: UpdateLocation :exec
update locations set name=?,qr_state=? where id=?;

-- name: ListSeeds :many
select * from seeds order by name collate nocase;

-- name: GetSeed :one
select * from seeds where id=?;

-- name: CreateSeed :one
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
  is_deer_resistant) values (?,?,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,?,?,?,?,?,?,?,?,?,?,?) returning *;

-- name: UpdateSeed :one
update seeds set name=?,family=?,updated_at=CURRENT_TIMESTAMP,year=?,edible=?,needs_trellis=?,needs_bird_netting=?,is_keto=?,is_native=?,is_invasive=?,is_cover_crop=?,grows_well_from_seed=?,is_bad_for_cats=?,is_deer_resistant=? where id=? returning *;
