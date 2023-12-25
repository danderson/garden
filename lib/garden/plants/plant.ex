defmodule Garden.Plants.Plant do
  use Ecto.Schema
  import Ecto.Changeset

  alias Garden.Plants.PlantLocation
  alias Garden.Seeds.Seed

  schema "plants" do
    field :name, :string
    belongs_to :seed, Seed
    has_many :locations, PlantLocation, preload_order: [desc: :start, desc: :end]

    field :current_location, :any, virtual: true

    timestamps(type: :utc_datetime)
  end

  def new_changeset(plant, attrs \\ %{}) do
    plant
    |> cast(attrs, [:name, :seed_id])
    |> validate_required([:name])
  end

  def edit_changeset(plant, attrs \\ %{}) do
    plant
    |> cast(attrs, [:name, :seed_id])
    |> validate_required([:name])
  end
end
