defmodule Garden.LocationsFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `Garden.Locations` context.
  """

  alias Garden.Locations

  @doc """
  Generate a location.
  """
  def location_fixture(attrs \\ %{}) do
    {:ok, location} =
      attrs
      |> Enum.into(%{
        name: "some name"
      })
      |> Locations.new()

    Locations.get!(location.id)
  end
end
