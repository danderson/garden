create table schema_migrations (
  version integer primary key,
  inserted_at text
);

create table seeds (
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
  family text
);

create table locations (
  id integer primary key autoincrement,
  name text not null,
  inserted_at text not null,
  updated_at text not null,
  qr_id text not null,
  qr_state text
);

create table locations_images (
  id integer primary key autoincrement,
  image_id text,
  location_id integer constraint locations_images_location_id_fkey references locations(id),
  inserted_at text not null,
  updated_at text not null
);

create table plants (
  id integer primary key autoincrement,
  name text,
  seed_id integer constraint plants_seed_id_fkey references seeds(id) on delete restrict,
  inserted_at text not null,
  updated_at text not null,
  name_from_seed integer
);

create index plants_seed_id_index on plants (
  seed_id
);

create table plant_locations (
  id integer primary key autoincrement,
  plant_id integer constraint plant_locations_plant_id_fkey references plants(id) on delete restrict,
  location_id integer constraint plant_locations_location_id_fkey references locations(id) on delete restrict,
  start text not null,
  end text null);

create index plant_locations_plant_id_index on plant_locations (plant_id);

create index plant_locations_location_id_index on plant_locations (location_id);
