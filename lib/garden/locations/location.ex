defmodule Garden.Locations.Location do
  use Ecto.Schema
  import Ecto.Changeset

  schema "locations" do
    field :name, :string
    field :qr_id, :string
    field :qr_state, Ecto.Enum, values: [:none, :wanted, :applied], default: :wanted

    has_many :images, Garden.Locations.Location.LocationImage
    has_many :plants, Garden.Plants.PlantLocation

    timestamps(type: :utc_datetime)
  end

  def upsert_changeset(location, attrs) do
    location
    |> cast(attrs, [:name, :qr_state])
    |> validate_required([:name, :qr_state])
  end

  defmodule LocationImage do
    use Ecto.Schema

    schema "locations_images" do
      field :image_id, :string
      belongs_to :location, Garden.Locations.Location
    end
  end
end
