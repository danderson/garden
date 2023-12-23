defmodule Garden.Repo.Migrations.PlantMultipleLocations do
  use Ecto.Migration

  def change do
    create table(:plant_locations) do
      add :plant_id, references(:plants, on_delete: :restrict)
      add :location_id, references(:locations, on_delete: :restrict)
      add :start, :date, null: false
      add :end, :date, null: true
    end

    create index(:plant_locations, [:plant_id])
    create index(:plant_locations, [:location_id])

    drop index(:plants, [:location_id])

    alter table(:plants) do
      remove :location_id
    end
  end
end
