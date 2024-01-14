-- name: ListSeeds :many
select * from seeds order by name collate nocase;

-- name: ListSeedsForSelector :many
select id,name from seeds order by name collate nocase;

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

-- name: ListLocations :many
select l.*,count(pl.id) as num_plants
  from locations as l left join plant_locations as pl on l.id=pl.location_id
  where pl.end is null
  group by l.id
  order by l.name collate nocase;

-- name: ListLocationsForSelector :many
select id,name from locations order by name collate nocase;

-- name: GetLocation :one
select * from locations where id=?;

-- name: GetPlantsInLocation :many
select p.name,pl.start,pl.end from locations as l
                                   inner join plant_locations as pl on l.id=pl.location_id
                                   inner join plants as p on p.id=pl.plant_id
 where l.id=?
 order by p.name collate nocase,
          pl.end desc,
          pl.start desc;

-- name: CreateLocation :one
insert into locations (
  name,
  qr_id,
  qr_state,
  inserted_at,
  updated_at) values (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) returning *;

-- name: UpdateLocation :one
update locations set name=?,qr_id=?,qr_state=?,updated_at=CURRENT_TIMESTAMP where id=? returning *;

-- name: ListPlants :many
select * from plants order by name collate nocase;

-- name: GetPlant :one
select * from plants where id=?;

-- name: CreatePlant :one
insert into plants (
  name,
  seed_id,
  name_from_seed,
  inserted_at,
  updated_at) values (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) returning *;

-- name: CreatePlantLocation :one
insert into plant_locations (
  plant_id,
  location_id,
  start) values (?,?,?) returning *;

-- name: GetPlantForUpdate :one
select pl.location_id,p.seed_id,p.name
  from plants as p inner join plant_locations as pl on pl.plant_id=p.id
 order by pl.start desc
 limit 1;

-- name: UpdatePlant :one
update plants set name=?,seed_id=?,name_from_seed=?,updated_at=CURRENT_TIMESTAMP where id=? returning *;
