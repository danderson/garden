defmodule Garden.PlantsTest do
  use Garden.DataCase

  alias Garden.Plants

  describe "plants" do
    alias Garden.Plants.Plant

    import Garden.PlantsFixtures
    import Garden.LocationsFixtures

    @invalid_attrs %{name: nil}

    test "list/0 returns all plants" do
      plant = plant_fixture()
      assert Plants.list() == [plant]
    end

    test "get!/1 returns the plant with given id" do
      plant = plant_fixture()
      assert Plants.get!(plant.id) == plant
    end

    test "new/1 with valid data creates a plant" do
      location = location_fixture()

      valid_attrs = %{
        plant: %{
          name: "some name"
        },
        location_id: location.id
      }

      assert {:ok, %Plant{} = plant} = Plants.new(valid_attrs)
      assert plant.name == "some name"
      assert plant.seed_id == nil
      plant = Plants.get!(plant.id, locations: :current)
      assert plant.current_location == location
    end

    test "new/1 with invalid data returns error changeset" do
      assert {:error, %Ecto.Changeset{}} = Plants.new(@invalid_attrs)
    end

    test "edit/2 with valid data updates the plant" do
      plant = plant_fixture()
      update_attrs = %{name: "some updated name"}

      assert {:ok, %Plant{} = plant} = Plants.edit(plant, update_attrs)
      assert plant.name == "some updated name"
    end

    test "edit/2 with invalid data returns error changeset" do
      plant = plant_fixture()
      assert {:error, %Ecto.Changeset{}} = Plants.edit(plant, @invalid_attrs)
      assert plant == Plants.get!(plant.id)
    end

    test "edit_changeset/1 returns a plant changeset" do
      plant = plant_fixture()
      assert %Ecto.Changeset{} = Plants.edit_changeset(plant)
    end
  end
end
