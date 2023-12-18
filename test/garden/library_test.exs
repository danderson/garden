defmodule Garden.LibraryTest do
  use Garden.DataCase

  alias Garden.Library

  describe "seeds" do
    alias Garden.Library.Seed

    import Garden.LibraryFixtures

    @invalid_attrs %{name: nil}

    test "list_seeds/0 returns all seeds" do
      seed = seed_fixture()
      assert Library.list_seeds() == [seed]
    end

    test "get_seed!/1 returns the seed with given id" do
      seed = seed_fixture()
      assert Library.get_seed!(seed.id) == seed
    end

    test "create_seed/1 with valid data creates a seed" do
      valid_attrs = %{name: "some name"}

      assert {:ok, %Seed{} = seed} = Library.create_seed(valid_attrs)
      assert seed.name == "some name"
    end

    test "create_seed/1 with invalid data returns error changeset" do
      assert {:error, %Ecto.Changeset{}} = Library.create_seed(@invalid_attrs)
    end

    test "update_seed/2 with valid data updates the seed" do
      seed = seed_fixture()
      update_attrs = %{name: "some updated name"}

      assert {:ok, %Seed{} = seed} = Library.update_seed(seed, update_attrs)
      assert seed.name == "some updated name"
    end

    test "update_seed/2 with invalid data returns error changeset" do
      seed = seed_fixture()
      assert {:error, %Ecto.Changeset{}} = Library.update_seed(seed, @invalid_attrs)
      assert seed == Library.get_seed!(seed.id)
    end

    test "delete_seed/1 deletes the seed" do
      seed = seed_fixture()
      assert {:ok, %Seed{}} = Library.delete_seed(seed)
      assert_raise Ecto.NoResultsError, fn -> Library.get_seed!(seed.id) end
    end

    test "change_seed/1 returns a seed changeset" do
      seed = seed_fixture()
      assert %Ecto.Changeset{} = Library.change_seed(seed)
    end
  end

  describe "locations" do
    alias Garden.Library.Location

    import Garden.LibraryFixtures

    @invalid_attrs %{name: nil}

    test "list_locations/0 returns all locations" do
      location = location_fixture()
      assert Library.list_locations() == [location]
    end

    test "get_location!/1 returns the location with given id" do
      location = location_fixture()
      assert Library.get_location!(location.id) == location
    end

    test "create_location/1 with valid data creates a location" do
      valid_attrs = %{name: "some name"}

      assert {:ok, %Location{} = location} = Library.create_location(valid_attrs)
      assert location.name == "some name"
    end

    test "create_location/1 with invalid data returns error changeset" do
      assert {:error, %Ecto.Changeset{}} = Library.create_location(@invalid_attrs)
    end

    test "update_location/2 with valid data updates the location" do
      location = location_fixture()
      update_attrs = %{name: "some updated name"}

      assert {:ok, %Location{} = location} = Library.update_location(location, update_attrs)
      assert location.name == "some updated name"
    end

    test "update_location/2 with invalid data returns error changeset" do
      location = location_fixture()
      assert {:error, %Ecto.Changeset{}} = Library.update_location(location, @invalid_attrs)
      assert location == Library.get_location!(location.id)
    end

    test "delete_location/1 deletes the location" do
      location = location_fixture()
      assert {:ok, %Location{}} = Library.delete_location(location)
      assert_raise Ecto.NoResultsError, fn -> Library.get_location!(location.id) end
    end

    test "change_location/1 returns a location changeset" do
      location = location_fixture()
      assert %Ecto.Changeset{} = Library.change_location(location)
    end
  end
end
