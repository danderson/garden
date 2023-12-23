defmodule Garden.Repo.Migrations.SeedFamily do
  use Ecto.Migration

  def change do
    alter table(:seeds) do
      add :family, :string
    end
  end
end
