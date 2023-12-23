defmodule Garden.LocationsTest do
  use Garden.DataCase

  alias Garden.Locations

  describe "locations" do
    alias Garden.Locations.Location

    import Garden.LocationsFixtures

    @invalid_attrs %{name: nil}

    test "list/0 returns all locations" do
      location = location_fixture()
      assert Locations.list() == [location]
    end

    test "get!/1 returns the location with given id" do
      location = location_fixture()
      assert Locations.get!(location.id) == location
    end

    test "new/1 with valid data creates a location" do
      valid_attrs = %{name: "some name"}

      assert {:ok, %Location{} = location} = Locations.new(valid_attrs)
      assert location.name == "some name"
    end

    test "new/1 with invalid data returns error changeset" do
      assert {:error, %Ecto.Changeset{}} = Locations.new(@invalid_attrs)
    end

    test "edit/2 with valid data updates the location" do
      location = location_fixture()
      update_attrs = %{name: "some updated name"}

      assert {:ok, %Location{} = location} = Locations.edit(location, update_attrs)
      assert location.name == "some updated name"
    end

    test "edit/2 with invalid data returns error changeset" do
      location = location_fixture()
      assert {:error, %Ecto.Changeset{}} = Locations.edit(location, @invalid_attrs)
      assert location == Locations.get!(location.id)
    end

    test "edit_changeset/1 returns a location changeset" do
      location = location_fixture()
      assert %Ecto.Changeset{} = Locations.upsert_changeset(location)
    end
  end
end
