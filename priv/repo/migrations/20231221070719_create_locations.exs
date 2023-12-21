defmodule Garden.Repo.Migrations.CreateLocations do
  use Ecto.Migration

  def change do
    create table(:locations) do
      add :name, :string

      timestamps(type: :utc_datetime)
    end

    create table(:locations_images) do
      add :image_id, :string
      add :location_id, references(:locations)

      timestamps(type: :utc_datetime)
    end
  end
end
