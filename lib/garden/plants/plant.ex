defmodule Garden.Plants.Plant do
  use Ecto.Schema
  import Ecto.Changeset

  alias Garden.Plants.PlantLocation
  alias Garden.Seeds.Seed

  schema "plants" do
    field :name, :string
    field :name_from_seed, :boolean
    belongs_to :seed, Seed
    has_many :locations, PlantLocation, preload_order: [desc: :start, desc: :end]

    field :current_location, :any, virtual: true

    timestamps(type: :utc_datetime)
  end

  def new_changeset(plant, attrs \\ %{}) do
    plant
    |> cast(attrs, [:name, :seed_id])
    |> validate_name_or_seed
  end

  def edit_changeset(plant, attrs \\ %{}) do
    plant
    |> cast(attrs, [:name, :seed_id])
    |> validate_name_or_seed
  end

  defp validate_name_or_seed(change) do
    if fetch_field!(change, :seed_id) == nil do
      validate_required(change, [:name], message: "must provide name or link to a seed")
    else
      change
    end
  end
end
