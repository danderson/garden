defmodule Garden.Repo.Migrations.LocationQrCode do
  use Ecto.Migration

  def change do
    alter table(:locations) do
      add :qr_id, :string, null: false
      add :qr_state, :string
    end
  end
end
