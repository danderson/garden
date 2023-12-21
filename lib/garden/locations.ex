defmodule Garden.Locations do
  @moduledoc """
  The Locations context.
  """

  import Ecto.Query, warn: false
  alias Garden.Repo

  alias Garden.Locations.Location

  defp base_query() do
    from l in Location, order_by: [:name], preload: [:plants]
  end

  def list_locations do
    base_query() |> Repo.all()
  end

  def get_location!(id) do
    base_query() |> Repo.get!(id)
  end

  def expand_location(%Location{} = location) do
    Repo.preload(location, [:plants, plants: :seed])
  end

  def new_location(), do: expand_location(%Location{})

  def create_location(attrs \\ %{}) do
    %Location{}
    |> Location.changeset(attrs)
    |> Repo.insert()
  end

  def update_location(%Location{} = location, attrs) do
    location
    |> Location.changeset(attrs)
    |> Repo.update()
  end

  def delete_location(%Location{} = location) do
    Repo.delete(location)
  end

  def change_location(%Location{} = location, attrs \\ %{}) do
    Location.changeset(location, attrs)
  end
end
