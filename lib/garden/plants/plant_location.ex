defmodule Garden.Plants.PlantLocation do
  use Ecto.Schema
  import Ecto.Changeset

  alias Garden.DateTime
  alias Garden.Plants.Plant
  alias Garden.Locations.Location

  schema "plant_locations" do
    belongs_to :plant, Plant
    belongs_to :location, Location
    field :start, DateTime
    field :end, DateTime
  end

  def new_changeset(attrs \\ %{}) do
    %__MODULE__{}
    |> cast(attrs, [:location_id])
    |> validate_required([:location_id])
    |> cast_assoc(:plant, with: &Plant.new_changeset/2, required: true)
    |> put_change(:start, DateTime.now!())
  end
end
