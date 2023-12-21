defmodule Garden.Plants.Plant do
  use Ecto.Schema
  import Ecto.Changeset

  schema "plants" do
    field :name, :string
    belongs_to :seed, Garden.Seeds.Seed
    belongs_to :location, Garden.Locations.Location

    timestamps(type: :utc_datetime)
  end

  @doc false
  def changeset(plant, attrs) do
    plant
    |> cast(attrs, [:name, :seed_id, :location_id])
    |> validate_required([:name])
    |> validate_required([:location_id], message: "plants have to be planted somewhere")
    |> foreign_key_constraint(:seed_id)
    |> foreign_key_constraint(:location_id)
  end
end
