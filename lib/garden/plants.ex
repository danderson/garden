defmodule Garden.Plants do
  @moduledoc """
  The Plants context.
  """

  import Ecto.Query, warn: false
  alias Garden.Repo

  alias Garden.Plants.Plant

  @doc """
  Returns the list of plants.

  ## Examples

      iex> list_plants()
      [%Plant{}, ...]

  """
  def list_plants do
    Repo.all(Plant)
  end

  @doc """
  Gets a single plant.

  Raises `Ecto.NoResultsError` if the Plant does not exist.

  ## Examples

      iex> get_plant!(123)
      %Plant{}

      iex> get_plant!(456)
      ** (Ecto.NoResultsError)

  """
  def get_plant!(id), do: Repo.get!(Plant, id)

  @doc """
  Creates a plant.

  ## Examples

      iex> create_plant(%{field: value})
      {:ok, %Plant{}}

      iex> create_plant(%{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def create_plant(attrs \\ %{}) do
    %Plant{}
    |> Plant.changeset(attrs)
    |> Repo.insert()
  end

  @doc """
  Updates a plant.

  ## Examples

      iex> update_plant(plant, %{field: new_value})
      {:ok, %Plant{}}

      iex> update_plant(plant, %{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def update_plant(%Plant{} = plant, attrs) do
    plant
    |> Plant.changeset(attrs)
    |> Repo.update()
  end

  @doc """
  Deletes a plant.

  ## Examples

      iex> delete_plant(plant)
      {:ok, %Plant{}}

      iex> delete_plant(plant)
      {:error, %Ecto.Changeset{}}

  """
  def delete_plant(%Plant{} = plant) do
    Repo.delete(plant)
  end

  @doc """
  Returns an `%Ecto.Changeset{}` for tracking plant changes.

  ## Examples

      iex> change_plant(plant)
      %Ecto.Changeset{data: %Plant{}}

  """
  def change_plant(%Plant{} = plant, attrs \\ %{}) do
    Plant.changeset(plant, attrs)
  end
end
