defmodule Garden.LibraryFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `Garden.Library` context.
  """

  @doc """
  Generate a seed.
  """
  def seed_fixture(attrs \\ %{}) do
    {:ok, seed} =
      attrs
      |> Enum.into(%{
        name: "some name",
        year: 2023,
        front_image_id: "1234",
        back_image_id: "2345"
      })
      |> Garden.Library.create_seed()

    seed
  end

  @doc """
  Generate a location.
  """
  def location_fixture(attrs \\ %{}) do
    {:ok, location} =
      attrs
      |> Enum.into(%{
        name: "some name"
      })
      |> Garden.Library.create_location()

    location
  end
end
