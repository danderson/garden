defmodule Garden.SeedsTest do
  use Garden.DataCase

  alias Garden.Seeds

  describe "seeds" do
    alias Garden.Seeds.Seed

    import Garden.SeedsFixtures

    @invalid_attrs %{name: nil}

    test "list/0 returns all seeds" do
      seed = seed_fixture()
      assert Seeds.list() == [seed]
    end

    test "get!/1 returns the seed with given id" do
      seed = seed_fixture()
      assert Seeds.get!(seed.id) == seed
    end

    test "new/1 with valid data creates a seed" do
      valid_attrs = %{
        name: "some name",
        year: 2023
      }

      assert {:ok, %Seed{} = seed} = Seeds.new(valid_attrs)
      assert seed.name == "some name"
      assert seed.year == 2023
    end

    test "new/2 with valid data creates a seed" do
      attrs = %{
        name: "some name",
        year: 2023
      }

      private_attrs = %{
        front_image_id: "1234",
        back_image_id: "2345"
      }

      assert {:ok, %Seed{} = seed} = Seeds.new(attrs, private_attrs)
      assert seed.name == "some name"
      assert seed.year == 2023
      assert seed.front_image_id == "1234"
      assert seed.back_image_id == "2345"
    end

    test "new/1 with invalid data returns error changeset" do
      assert {:error, %Ecto.Changeset{}} = Seeds.new(@invalid_attrs)
    end

    test "edit/2 with valid data updates the seed" do
      seed = seed_fixture()
      update_attrs = %{name: "some updated name"}

      assert {:ok, %Seed{} = seed} = Seeds.edit(seed, update_attrs)
      assert seed.name == "some updated name"
    end

    test "edit/2 with invalid data returns error changeset" do
      seed = seed_fixture()
      assert {:error, %Ecto.Changeset{}} = Seeds.edit(seed, @invalid_attrs)
      assert seed == Seeds.get!(seed.id)
    end

    test "upsert_changeset/1 returns a seed changeset" do
      seed = seed_fixture()
      assert %Ecto.Changeset{} = Seeds.upsert_changeset(seed)
    end
  end
end
