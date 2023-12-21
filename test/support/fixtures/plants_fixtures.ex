defmodule Garden.PlantsFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `Garden.Plants` context.
  """

  @doc """
  Generate a plant.
  """
  def plant_fixture(attrs \\ %{}) do
    {:ok, plant} =
      attrs
      |> Enum.into(%{
        name: "some name"
      })
      |> Map.put_new_lazy(:location_id, fn ->
        Garden.LocationsFixtures.location_fixture().id
      end)
      |> Garden.Plants.create_plant()

    Garden.Plants.expand_plant(plant)
  end
end
