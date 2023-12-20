defmodule Garden.Library do
  @moduledoc """
  The Library context.
  """

  import Ecto.Query, warn: false
  alias Garden.Repo

  alias Garden.Library.Seed

  @doc """
  Returns the list of seeds.

  ## Examples

      iex> list_seeds()
      [%Seed{}, ...]

  """
  def list_seeds do
    Repo.all(Seed)
  end

  @doc """
  Gets a single seed.

  Raises `Ecto.NoResultsError` if the Seed does not exist.

  ## Examples

      iex> get_seed!(123)
      %Seed{}

      iex> get_seed!(456)
      ** (Ecto.NoResultsError)

  """
  def get_seed!(id), do: Repo.get!(Seed, id)

  @doc """
  Creates a seed.

  ## Examples

      iex> create_seed(%{field: value})
      {:ok, %Seed{}}

      iex> create_seed(%{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def create_seed(attrs \\ %{}, front_image_id \\ nil, back_image_id \\ nil) do
    %Seed{}
    |> Seed.changeset(attrs, front_image_id, back_image_id)
    |> Repo.insert()
  end

  @doc """
  Updates a seed.

  ## Examples

      iex> update_seed(seed, %{field: new_value})
      {:ok, %Seed{}}

      iex> update_seed(seed, %{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def update_seed(%Seed{} = seed, attrs, front_image_id \\ nil, back_image_id \\ nil) do
    seed
    |> Seed.changeset(attrs, front_image_id, back_image_id)
    |> Repo.update()
  end

  @doc """
  Deletes a seed.

  ## Examples

      iex> delete_seed(seed)
      {:ok, %Seed{}}

      iex> delete_seed(seed)
      {:error, %Ecto.Changeset{}}

  """
  def delete_seed(%Seed{:front_image_id => front_id, :back_image_id => back_id} = seed) do
    res = Repo.delete(seed)
    Images.delete(front_id)
    Images.delete(back_id)
    res
  end

  @doc """
  Returns an `%Ecto.Changeset{}` for tracking seed changes.

  ## Examples

      iex> change_seed(seed)
      %Ecto.Changeset{data: %Seed{}}

  """
  def change_seed(%Seed{} = seed, attrs \\ %{}) do
    Seed.changeset(seed, attrs)
  end

  alias Garden.Library.Location

  @doc """
  Returns the list of locations.

  ## Examples

      iex> list_locations()
      [%Location{}, ...]

  """
  def list_locations do
    Repo.all(Location)
  end

  @doc """
  Gets a single location.

  Raises `Ecto.NoResultsError` if the Location does not exist.

  ## Examples

      iex> get_location!(123)
      %Location{}

      iex> get_location!(456)
      ** (Ecto.NoResultsError)

  """
  def get_location!(id), do: Repo.get!(Location, id)

  @doc """
  Creates a location.

  ## Examples

      iex> create_location(%{field: value})
      {:ok, %Location{}}

      iex> create_location(%{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def create_location(attrs \\ %{}) do
    %Location{}
    |> Location.changeset(attrs)
    |> Repo.insert()
  end

  @doc """
  Updates a location.

  ## Examples

      iex> update_location(location, %{field: new_value})
      {:ok, %Location{}}

      iex> update_location(location, %{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def update_location(%Location{} = location, attrs) do
    location
    |> Location.changeset(attrs)
    |> Repo.update()
  end

  @doc """
  Deletes a location.

  ## Examples

      iex> delete_location(location)
      {:ok, %Location{}}

      iex> delete_location(location)
      {:error, %Ecto.Changeset{}}

  """
  def delete_location(%Location{} = location) do
    Repo.delete(location)
  end

  @doc """
  Returns an `%Ecto.Changeset{}` for tracking location changes.

  ## Examples

      iex> change_location(location)
      %Ecto.Changeset{data: %Location{}}

  """
  def change_location(%Location{} = location, attrs \\ %{}) do
    Location.changeset(location, attrs)
  end
end
