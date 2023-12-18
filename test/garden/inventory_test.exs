defmodule Garden.InventoryTest do
  use Garden.DataCase

  alias Garden.Inventory

  describe "plants" do
    alias Garden.Inventory.Plant

    import Garden.InventoryFixtures

    @invalid_attrs %{name: nil}

    test "list_plants/0 returns all plants" do
      plant = plant_fixture()
      assert Inventory.list_plants() == [plant]
    end

    test "get_plant!/1 returns the plant with given id" do
      plant = plant_fixture()
      assert Inventory.get_plant!(plant.id) == plant
    end

    test "create_plant/1 with valid data creates a plant" do
      valid_attrs = %{name: "some name"}

      assert {:ok, %Plant{} = plant} = Inventory.create_plant(valid_attrs)
      assert plant.name == "some name"
    end

    test "create_plant/1 with invalid data returns error changeset" do
      assert {:error, %Ecto.Changeset{}} = Inventory.create_plant(@invalid_attrs)
    end

    test "update_plant/2 with valid data updates the plant" do
      plant = plant_fixture()
      update_attrs = %{name: "some updated name"}

      assert {:ok, %Plant{} = plant} = Inventory.update_plant(plant, update_attrs)
      assert plant.name == "some updated name"
    end

    test "update_plant/2 with invalid data returns error changeset" do
      plant = plant_fixture()
      assert {:error, %Ecto.Changeset{}} = Inventory.update_plant(plant, @invalid_attrs)
      assert plant == Inventory.get_plant!(plant.id)
    end

    test "delete_plant/1 deletes the plant" do
      plant = plant_fixture()
      assert {:ok, %Plant{}} = Inventory.delete_plant(plant)
      assert_raise Ecto.NoResultsError, fn -> Inventory.get_plant!(plant.id) end
    end

    test "change_plant/1 returns a plant changeset" do
      plant = plant_fixture()
      assert %Ecto.Changeset{} = Inventory.change_plant(plant)
    end
  end
end
