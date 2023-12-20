defmodule Garden.Repo.Migrations.SeedTwoImages do
  use Ecto.Migration

  def change do
    rename table(:seeds), :image_id, to: :front_image_id

    alter table(:seeds) do
      add :back_image_id, :string
    end
  end
end
