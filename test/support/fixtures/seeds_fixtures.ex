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
        year: 2023,
        front_image_id: "1234",
        back_image_id: "2345"
      })
      |> Garden.Seeds.create_seed()

    Garden.Seeds.expand_seed(seed)
  end
end
