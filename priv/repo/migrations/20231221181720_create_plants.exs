defmodule Garden.Repo.Migrations.CreatePlants do
  use Ecto.Migration

  def change do
    create table(:plants) do
      add :name, :string
      add :seed_id, references(:seeds, on_delete: :restrict)
      add :location_id, references(:locations, on_delete: :restrict)

      timestamps(type: :utc_datetime)
    end

    create index(:plants, [:seed_id])
    create index(:plants, [:location_id])
  end
end
