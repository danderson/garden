defmodule Garden.Seeds do
  @moduledoc """
  The Seeds context.
  """

  import Ecto.Query, warn: false
  alias Garden.Repo

  alias Garden.Seeds.Seed

  defp base_query() do
    from s in Seed, order_by: [:name], preload: [:plants, plants: :location]
  end

  def list_seeds do
    base_query() |> Repo.all()
  end

  def get_seed!(id) do
    base_query() |> Repo.get!(id)
  end

  def expand_seed(%Seed{} = seed), do: seed |> Repo.preload([:plants, plants: :location])

  def new_seed() do
    expand_seed(%Seed{})
  end

  def new_seed_changeset() do
    new_seed() |> Seed.changeset(%{year: Date.utc_today().year})
  end

  def create_seed(attrs \\ %{}, private_attrs \\ %{}) do
    %Seed{}
    |> Seed.changeset(attrs, private_attrs)
    |> Repo.insert()
  end

  def update_seed(%Seed{} = seed, attrs, private_attrs \\ %{}) do
    change =
      seed
      |> Seed.changeset(attrs, private_attrs)

    case Repo.update(change) do
      {:ok, _} = res ->
        Seed.replaced_images(change) |> Enum.each(&Images.delete(:seed, &1))
        res

      {:error, _} = res ->
        Seed.new_images(change) |> Enum.each(&Images.delete(:seed, &1))
        res
    end
  end

  def store_seed_image(src), do: Images.store(:seeds, src)
  def seed_image(id, size), do: Images.url(:seeds, id, size)

  def delete_seed(%Seed{:front_image_id => front_id, :back_image_id => back_id} = seed) do
    res = Repo.delete(seed)
    Images.delete(:seeds, front_id)
    Images.delete(:seeds, back_id)
    res
  end

  def change_seed(%Seed{} = seed, attrs \\ %{}) do
    Seed.changeset(seed, attrs)
  end
end
