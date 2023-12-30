defmodule Garden.Seeds do
  @moduledoc """
  The Seeds context.
  """

  import Ecto.Query, warn: false
  alias Garden.Repo

  alias Garden.Seeds.Seed
  alias Garden.Plants

  defp query(kw) do
    from(s in Seed, order_by: [asc: :name])
    |> query_plants(kw[:plants], kw[:locations])
  end

  def list(kw \\ []) do
    query(kw) |> Repo.all()
  end

  def get!(id, kw \\ []) do
    query(kw) |> Repo.get!(id)
  end

  defp query_plants(q, nil, _), do: q

  defp query_plants(q, true, locs) do
    from(q, preload: [plants: ^(&Plants.list_for_seeds(&1, locations: locs))])
  end

  defdelegate upsert_changeset(seed, attrs \\ %{}, private_attrs \\ %{}), to: Seed

  def new(attrs \\ %{}, private_attrs \\ %{}) do
    upsert_changeset(%Seed{}, attrs, private_attrs) |> Repo.insert()
  end

  def new!(attrs \\ %{}, private_attrs \\ %{}) do
    {:ok, seed} = new(attrs, private_attrs)
    seed
  end

  def edit(%Seed{} = seed, attrs, private_attrs \\ %{}) do
    change = upsert_changeset(seed, attrs, private_attrs)

    case Repo.update(change) do
      {:ok, seed} = res ->
        Seed.replaced_images(change) |> Enum.each(&Images.delete(:seed, &1))
        Plants.update_names_from_seed(seed.id, seed.name)
        res

      {:error, _} = res ->
        Seed.new_images(change) |> Enum.each(&Images.delete(:seed, &1))
        res
    end
  end

  def store_seed_image(src), do: Images.store(:seeds, src)
  def seed_image(id, size), do: Images.url(:seeds, id, size)
end
