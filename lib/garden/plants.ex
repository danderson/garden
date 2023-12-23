defmodule Garden.Plants do
  @moduledoc """
  The Plants context.
  """

  import Ecto.Query, warn: false
  alias Garden.Repo

  alias Garden.DateTime
  alias Garden.Plants.Plant
  alias Garden.Plants.PlantLocation
  alias Ecto.Multi
  alias Ecto.Changeset

  defp query(kw) do
    from(p in Plant)
    |> query_locations(kw[:locations])
    |> query_seed(kw[:seed])
  end

  def list(kw \\ []) do
    query(kw) |> Repo.all() |> Enum.map(&set_current_location/1)
  end

  def get!(id, kw \\ []), do: query(kw) |> Repo.get!(id) |> set_current_location

  defp query_locations(q, nil), do: q
  defp query_locations(q, :all), do: from(q, preload: [locations: :location])

  defp query_locations(q, :current) do
    pl = from(pl in PlantLocation, where: is_nil(pl.end), preload: [:location])
    from(p in q, preload: [locations: ^pl])
  end

  defp query_seed(q, nil), do: q
  defp query_seed(q, true), do: from(q, preload: [:seed])

  defp set_current_location(%Plant{locations: location} = plant) when is_list(location) do
    case Enum.find(plant.locations, fn pl -> pl.end == nil end) do
      %PlantLocation{} = pl -> Map.put(plant, :current_location, pl.location)
      _ -> plant
    end
  end

  defp set_current_location(%Plant{} = plant), do: plant

  defdelegate new_changeset(attrs), to: PlantLocation

  def new(attrs \\ %{}) do
    case new_changeset(attrs) |> Repo.insert() do
      {:ok, pl} -> {:ok, pl.plant}
      _ = err -> err
    end
  end

  defdelegate edit_changeset(plant, attrs \\ %{}), to: Plant

  def edit(%Plant{} = plant, attrs \\ %{}) do
    edit_changeset(plant, attrs) |> Repo.update()
  end

  def move(plant_id, location_id) do
    {:ok, _} =
      Multi.new()
      |> Multi.one(
        :current,
        from(pl in PlantLocation,
          where: pl.plant_id == ^plant_id,
          where: pl.location_id != ^location_id,
          where: is_nil(pl.end)
        )
      )
      |> Multi.update(:remove_current, fn %{current: current} ->
        Changeset.change(current, end: DateTime.now!())
      end)
      |> Multi.insert(:add_new, %PlantLocation{
        plant_id: plant_id,
        location_id: location_id,
        start: DateTime.now!()
      })
      |> Repo.transaction()

    {:ok, plant_id}
  end
end
