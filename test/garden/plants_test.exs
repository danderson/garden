defmodule Garden.PlantsTest do
  use Garden.DataCase

  alias Garden.Plants

  describe "plants" do
    alias Garden.Plants.Plant

    import Garden.PlantsFixtures

    @invalid_attrs %{name: nil}

    test "list_plants/0 returns all plants" do
      plant = plant_fixture()
      assert Plants.list_plants() == [plant]
    end

    test "get_plant!/1 returns the plant with given id" do
      plant = plant_fixture()
      assert Plants.get_plant!(plant.id) == plant
    end

    test "create_plant/1 with valid data creates a plant" do
      valid_attrs = %{name: "some name"}

      assert {:ok, %Plant{} = plant} = Plants.create_plant(valid_attrs)
      assert plant.name == "some name"
    end

    test "create_plant/1 with invalid data returns error changeset" do
      assert {:error, %Ecto.Changeset{}} = Plants.create_plant(@invalid_attrs)
    end

    test "update_plant/2 with valid data updates the plant" do
      plant = plant_fixture()
      update_attrs = %{name: "some updated name"}

      assert {:ok, %Plant{} = plant} = Plants.update_plant(plant, update_attrs)
      assert plant.name == "some updated name"
    end

    test "update_plant/2 with invalid data returns error changeset" do
      plant = plant_fixture()
      assert {:error, %Ecto.Changeset{}} = Plants.update_plant(plant, @invalid_attrs)
      assert plant == Plants.get_plant!(plant.id)
    end

    test "delete_plant/1 deletes the plant" do
      plant = plant_fixture()
      assert {:ok, %Plant{}} = Plants.delete_plant(plant)
      assert_raise Ecto.NoResultsError, fn -> Plants.get_plant!(plant.id) end
    end

    test "change_plant/1 returns a plant changeset" do
      plant = plant_fixture()
      assert %Ecto.Changeset{} = Plants.change_plant(plant)
    end
  end
end
