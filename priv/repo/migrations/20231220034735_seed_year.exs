defmodule Garden.Repo.Migrations.SeedYear do
  use Ecto.Migration

  def change do
    alter table(:seeds) do
      add :year, :integer
    end
  end
end
