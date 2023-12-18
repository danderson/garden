defmodule Garden.Repo.Migrations.SeedPhoto do
  use Ecto.Migration

  def change do
    alter table(:seeds) do
      add :photo, :binary
    end
  end
end
