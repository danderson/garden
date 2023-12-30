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
  alias Garden.Seeds.Seed

  defp query(kw) do
    from(p in Plant, order_by: [asc: :name])
    |> query_locations(kw[:locations])
    |> query_seed(kw[:seed])
  end

  def list(kw \\ []) do
    query(kw) |> Repo.all() |> Enum.map(&set_current_location/1)
  end

  def list_for_seeds(seed_ids, kw \\ []) do
    query(kw)
    |> where([p], p.seed_id in ^seed_ids)
    |> Repo.all()
    |> Enum.map(&set_current_location/1)
  end

  def get!(id, kw \\ []), do: query(kw) |> Repo.get!(id) |> set_current_location

  defp query_locations(q, nil), do: q
  defp query_locations(q, :all), do: from(q, preload: [locations: :location])

  defp query_locations(q, :current) do
    from(
      p in q,
      join: pl in assoc(p, :locations),
      where: is_nil(pl.end),
      preload: [locations: {pl, [:location]}]
    )
  end

  defp query_seed(q, nil), do: q
  defp query_seed(q, true), do: from(q, preload: [:seed])

  def set_current_location(%Plant{locations: location} = plant) when is_list(location) do
    case Enum.find(plant.locations, fn pl -> pl.end == nil end) do
      %PlantLocation{} = pl -> Map.put(plant, :current_location, pl.location)
      _ -> plant
    end
  end

  def set_current_location(%Plant{} = plant), do: plant

  defdelegate new_changeset(attrs), to: PlantLocation

  def new(attrs \\ %{}) do
    change = new_changeset(attrs)

    if change.valid? do
      new_from_change(
        change,
        change
        |> Changeset.get_assoc(:plant)
        |> Changeset.fetch_field!(:name)
      )
    else
      {:error, change}
    end
  end

  defp new_from_change(change, nil) do
    seed_id = change |> Changeset.get_assoc(:plant) |> Changeset.fetch_field!(:seed_id)

    {:ok, %{create: plantloc}} =
      Multi.new()
      |> Multi.one(:seed, from(s in Seed, where: s.id == ^seed_id))
      |> Multi.run(:adjust, fn _repo, %{seed: seed} ->
        change =
          Changeset.get_assoc(change, :plant)
          |> Changeset.change(%{
            name_from_seed: true,
            name: seed.name
          })
          |> then(&Changeset.put_assoc(change, :plant, &1))

        {:ok, change}
      end)
      |> Multi.insert(:create, fn %{adjust: adjust} ->
        adjust
      end)
      |> Repo.transaction()

    {:ok, plantloc.plant}
  end

  defp new_from_change(change, _), do: Repo.insert(change)

  def new!(attrs \\ %{}) do
    {:ok, plant} = new(attrs)
    plant
  end

  defdelegate edit_changeset(plant, attrs \\ %{}), to: Plant

  def edit(%Plant{} = plant, attrs \\ %{}) do
    change = edit_changeset(plant, attrs)

    if change.valid? do
      edit_from_change(change)
    else
      {:error, change}
    end
  end

  defp edit_from_change(change) do
    cond do
      Changeset.changed?(change, :name) ->
        new_name = Changeset.fetch_change!(change, :name)

        if new_name == nil or new_name == "" do
          set_name_from_seed(change)
        else
          Changeset.change(change, %{name_from_seed: false})
          |> Repo.update()
        end

      Changeset.changed?(change, :seed_id) and Changeset.fetch_field!(change, :name_from_seed) ->
        set_name_from_seed(change)

      true ->
        Repo.update(change)
    end
  end

  defp set_name_from_seed(change) do
    {:ok, %{update: plant}} =
      Multi.new()
      |> Multi.one(
        :seed,
        from(s in Seed, where: s.id == ^Changeset.fetch_field!(change, :seed_id))
      )
      |> Multi.run(:update_name, fn _repo, %{seed: seed} ->
        {:ok, Changeset.change(change, %{name: seed.name, name_from_seed: true})}
      end)
      |> Multi.update(:update, fn %{update_name: change} -> change end)
      |> Repo.transaction()

    {:ok, plant}
  end

  def update_names_from_seed(seed_id, new_name) do
    from(p in Plant, where: p.seed_id == ^seed_id and p.name_from_seed == true)
    |> Repo.update_all(set: [name: new_name])
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
