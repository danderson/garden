defmodule Garden.Locations.Location do
  use Ecto.Schema
  import Ecto.Changeset

  schema "locations" do
    field :name, :string
    has_many :images, Garden.Locations.Location.LocationImage

    timestamps(type: :utc_datetime)
  end

  @doc false
  def changeset(location, attrs) do
    location
    |> cast(attrs, [:name])
    |> validate_required([:name])
  end

  defmodule LocationImage do
    use Ecto.Schema

    schema "locations_images" do
      field :image_id, :string
      belongs_to :location, Garden.Locations.Location
    end
  end
end
