defmodule Garden.Plants.Plant do
  use Ecto.Schema
  import Ecto.Changeset

  schema "plants" do
    field :name, :string

    timestamps(type: :utc_datetime)
  end

  @doc false
  def changeset(plant, attrs) do
    plant
    |> cast(attrs, [:name])
    |> validate_required([:name])
  end
end
