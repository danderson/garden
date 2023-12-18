defmodule Garden.Repo.Migrations.CreatePlants do
  use Ecto.Migration

  def change do
    create table(:plants) do
      add :name, :string

      timestamps(type: :utc_datetime)
    end
  end
end
