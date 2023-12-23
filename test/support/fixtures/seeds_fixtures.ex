defmodule Garden.SeedsFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `Garden.Seeds` context.
  """

  @doc """
  Generate a seed.
  """
  def seed_fixture(attrs \\ %{}) do
    {:ok, seed} =
      attrs
      |> Enum.into(%{
        name: "some name",
        year: 2023
      })
      |> Garden.Seeds.new()

    Garden.Seeds.get!(seed.id)
  end
end
