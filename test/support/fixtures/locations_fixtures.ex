defmodule Garden.LocationsFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `Garden.Locations` context.
  """

  @doc """
  Generate a location.
  """
  def location_fixture(attrs \\ %{}) do
    {:ok, location} =
      attrs
      |> Enum.into(%{
        name: "some name"
      })
      |> Garden.Locations.create_location()

    location
  end
end
