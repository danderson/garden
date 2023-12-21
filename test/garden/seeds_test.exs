defmodule Garden.SeedsTest do
  use Garden.DataCase

  alias Garden.Seeds

  describe "seeds" do
    alias Garden.Seeds.Seed

    import Garden.SeedsFixtures

    @invalid_attrs %{name: nil}

    test "list_seeds/0 returns all seeds" do
      seed = seed_fixture()
      assert Seeds.list_seeds() == [seed]
    end

    test "get_seed!/1 returns the seed with given id" do
      seed = seed_fixture()
      assert Seeds.get_seed!(seed.id) == seed
    end

    test "create_seed/1 with valid data creates a seed" do
      valid_attrs = %{
        name: "some name",
        year: 2023,
        front_image_id: "1234",
        back_image_id: "2345"
      }

      assert {:ok, %Seed{} = seed} = Seeds.create_seed(valid_attrs)
      assert seed.name == "some name"
      assert seed.year == 2023
    end

    test "create_seed/2 with valid data creates a seed" do
      attrs = %{
        name: "some name",
        year: 2023
      }

      private_attrs = %{
        front_image_id: "1234",
        back_image_id: "2345"
      }

      assert {:ok, %Seed{} = seed} = Seeds.create_seed(attrs, private_attrs)
      assert seed.name == "some name"
      assert seed.year == 2023
      assert seed.front_image_id == "1234"
      assert seed.back_image_id == "2345"
    end

    test "create_seed/1 with invalid data returns error changeset" do
      assert {:error, %Ecto.Changeset{}} = Seeds.create_seed(@invalid_attrs)
    end

    test "update_seed/2 with valid data updates the seed" do
      seed = seed_fixture()
      update_attrs = %{name: "some updated name"}

      assert {:ok, %Seed{} = seed} = Seeds.update_seed(seed, update_attrs)
      assert seed.name == "some updated name"
    end

    test "update_seed/2 with invalid data returns error changeset" do
      seed = seed_fixture()
      assert {:error, %Ecto.Changeset{}} = Seeds.update_seed(seed, @invalid_attrs)
      assert seed == Seeds.get_seed!(seed.id)
    end

    test "delete_seed/1 deletes the seed" do
      seed = seed_fixture()
      assert {:ok, %Seed{}} = Seeds.delete_seed(seed)
      assert_raise Ecto.NoResultsError, fn -> Seeds.get_seed!(seed.id) end
    end

    test "change_seed/1 returns a seed changeset" do
      seed = seed_fixture()
      assert %Ecto.Changeset{} = Seeds.change_seed(seed)
    end
  end
end
