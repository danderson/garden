defmodule Garden.Plants do
  @moduledoc """
  The Plants context.
  """

  import Ecto.Query, warn: false
  alias Garden.Repo

  alias Garden.Plants.Plant

  defp base_query() do
    from p in Plant, preload: [:location, :seed], order_by: [:name]
  end

  def list_plants do
    base_query() |> Repo.all()
  end

  def get_plant!(id) do
    base_query() |> Repo.get!(id)
  end

  def expand_plant(plant), do: Repo.preload(plant, [:location, :seed])

  def new_plant(), do: %Plant{} |> expand_plant()

  def create_plant(attrs \\ %{}) do
    %Plant{}
    |> Plant.changeset(attrs)
    |> Repo.insert()
  end

  def update_plant(%Plant{} = plant, attrs) do
    plant
    |> Plant.changeset(attrs)
    |> Repo.update()
  end

  def delete_plant(%Plant{} = plant) do
    Repo.delete(plant)
  end

  def change_plant(%Plant{} = plant, attrs \\ %{}) do
    Plant.changeset(plant, attrs)
  end
end
