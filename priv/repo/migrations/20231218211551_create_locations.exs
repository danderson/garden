defmodule Garden.Repo.Migrations.CreateLocations do
  use Ecto.Migration

  def change do
    create table(:locations) do
      add :name, :string

      timestamps(type: :utc_datetime)
    end
  end
end
