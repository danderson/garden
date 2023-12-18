defmodule Garden.InventoryFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `Garden.Inventory` context.
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
      |> Garden.Inventory.create_plant()

    plant
  end
end
