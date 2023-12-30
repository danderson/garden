defmodule Garden.Locations do
  @moduledoc """
  The Locations context.
  """

  import Ecto.Query, warn: false
  alias Garden.Repo

  alias Garden.Locations.Location
  alias Garden.Plants.PlantLocation

  defp query(kw) do
    from(l in Location, order_by: [asc: :name])
    |> query_plants(kw[:plants])
    |> query_seeds(kw[:seeds])
  end

  def list(kw \\ []) do
    query(kw) |> Repo.all()
  end

  def get!(id, kw \\ []) do
    query(kw) |> Repo.get!(id)
  end

  def get_from_qr(qr_id, kw \\ []) do
    query(kw) |> where([l], l.qr_id == ^qr_id) |> Repo.one()
  end

  defp query_plants(q, nil), do: q
  defp query_plants(q, :all), do: from(q, preload: [plants: :plant])

  defp query_plants(q, :current) do
    pl = from(pl in PlantLocation, where: is_nil(pl.end), preload: [:plant])
    from(q, preload: [plants: ^pl])
  end

  defp query_seeds(q, nil), do: q
  defp query_seeds(q, true), do: from(q, preload: [plants: [plant: :seed]])

  defdelegate upsert_changeset(loc, attrs \\ %{}), to: Location

  def new(attrs \\ %{}) do
    upsert_changeset(%Location{}, attrs) |> Repo.insert()
  end

  def new!(attrs \\ %{}) do
    {:ok, location} = new(attrs)
    location
  end

  def edit(location, attrs \\ %{}) do
    upsert_changeset(location, attrs) |> Repo.update()
  end
end
