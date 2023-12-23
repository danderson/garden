defmodule Garden.PlantsFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `Garden.Plants` context.
  """

  @doc """
  Generate a plant.
  """
  def plant_fixture(attrs \\ %{}) do
    {loc, attrs} =
      Map.pop_lazy(attrs, :location_id, fn -> Garden.LocationFixtures.location_fixture().id end)

    plant =
      attrs
      |> Enum.into(%{
        name: "some name"
      })

    {:ok, plant} = Garden.Plants.new(%{location_id: loc, plant: plant})

    Garden.Plants.get!(plant.id, locations: :current)
  end
end
