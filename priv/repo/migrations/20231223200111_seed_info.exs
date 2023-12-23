defmodule Garden.Repo.Migrations.SeedInfo do
  use Ecto.Migration

  def change do
    alter table(:seeds) do
      add :edible, :boolean, null: true
      add :needs_trellis, :boolean, null: true
      add :needs_bird_netting, :boolean, null: true
      add :is_keto, :boolean, null: true
      add :is_native, :boolean, null: true
      add :is_invasive, :boolean, null: true
      add :is_cover_crop, :boolean, null: true
      add :grows_well_from_seed, :boolean, null: true
      add :is_bad_for_cats, :boolean, null: true
      add :is_deer_resistant, :boolean, null: true
      add :type, :string
      add :lifespan, :string
    end
  end
end
