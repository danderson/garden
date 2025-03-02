CREATE TABLE schema_migrations (
  version integer primary key,
  inserted_at text
);

CREATE TABLE seeds (
  id integer primary key autoincrement,
  name text not null,
  inserted_at text not null,
  updated_at text not null,
  front_image_id text,
  back_image_id text,
  year integer,
  edible integer,
  needs_trellis integer,
  needs_bird_netting integer,
  is_keto integer,
  is_native integer,
  is_invasive integer,
  is_cover_crop integer,
  grows_well_from_seed integer,
  is_bad_for_cats integer,
  is_deer_resistant integer,
  type text,
  lifespan text,
  family text,
  latin_name text not null default "",
  needs_stratification integer,
  sun_type text check (sun_type in ('Full', 'Partial', 'Shade', 'Unknown')) not null default 'Unknown',
  soil_type text check (soil_type in ('Dry', 'Wet', 'Both', 'Unknown')) not null default 'Unknown'
);

CREATE TABLE locations (
  id integer primary key autoincrement,
  name text not null,
  inserted_at text not null,
  updated_at text not null,
  qr_id text not null,
  qr_state text
);

CREATE TABLE locations_images (
  id integer primary key autoincrement,
  image_id text,
  location_id integer constraint locations_images_location_id_fkey references locations(id),
  inserted_at text not null,
  updated_at text not null
);

CREATE TABLE plants (
  id integer primary key autoincrement,
  name text not null,
  seed_id integer constraint plants_seed_id_fkey references seeds(id) on delete restrict,
  inserted_at text not null,
  updated_at text not null,
  name_from_seed integer not null
);

CREATE INDEX plants_seed_id_index on plants (
  seed_id
);

CREATE TABLE plant_locations (
  id integer primary key autoincrement,
  plant_id integer not null constraint plant_locations_plant_id_fkey references plants(id) on delete restrict,
  location_id integer not null constraint plant_locations_location_id_fkey references locations(id) on delete restrict,
  start text not null,
  end text null
);

CREATE INDEX plant_locations_plant_id_index on plant_locations (plant_id);

CREATE INDEX plant_locations_location_id_index on plant_locations (location_id);
