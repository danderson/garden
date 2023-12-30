defmodule Garden.Repo.Migrations.PlantNameFromSeed do
  use Ecto.Migration

  def change do
    alter table(:plants) do
      add :name_from_seed, :boolean
    end

    execute "UPDATE plants SET name_from_seed=0"
  end
end
