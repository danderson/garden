defmodule Garden.Repo.Migrations.CreateSeeds do
  use Ecto.Migration

  def change do
    create table(:seeds) do
      add :name, :string

      timestamps(type: :utc_datetime)
    end
  end
end
