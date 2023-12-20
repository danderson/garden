defmodule Garden.Repo.Migrations.SeedImage do
  use Ecto.Migration

  def change do
    alter table(:seeds) do
      add :image_id, :string
    end
  end
end
