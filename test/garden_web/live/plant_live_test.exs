defmodule GardenWeb.PlantLiveTest do
  use GardenWeb.ConnCase

  import Phoenix.LiveViewTest
  import Garden.PlantsFixtures
  import Garden.LocationsFixtures

  @update_attrs %{name: "some updated name"}
  @invalid_attrs %{name: nil}

  defp create_plant(_) do
    location = location_fixture()
    plant = plant_fixture(%{location_id: location.id})
    %{location: location, plant: plant}
  end

  describe "Index" do
    setup [:create_plant]

    test "lists all plants", %{conn: conn, plant: plant} do
      {:ok, _index_live, html} = live(conn, ~p"/plants")

      assert html =~ "My Plants"
      assert html =~ plant.name
    end

    test "saves new plant", %{conn: conn, location: location} do
      {:ok, index_live, _html} = live(conn, ~p"/plants")

      assert index_live
             |> element("a", "Add")
             |> render_click() =~ "Add"

      assert_patch(index_live, ~p"/plants/new")

      invalid_form_resp =
        index_live
        |> form("#plant-form", plant: @invalid_attrs)
        |> render_change()

      assert invalid_form_resp =~ "can&#39;t be blank"
      assert invalid_form_resp =~ "plants have to be planted somewhere"

      attrs = %{
        name: "some name",
        location_id: location.id
      }

      assert index_live
             |> form("#plant-form", plant: attrs)
             |> render_submit()

      assert_patch(index_live, ~p"/plants")

      html = render(index_live)
      assert html =~ "Plant created successfully"
      assert html =~ "some name"
    end

    test "updates plant in listing", %{conn: conn, plant: plant} do
      {:ok, index_live, _html} = live(conn, ~p"/plants")

      assert index_live
             |> element("#plants-#{plant.id} a", "Edit")
             |> render_click() =~ "Edit Plant"

      assert_patch(index_live, ~p"/plants/#{plant}/edit")

      assert index_live
             |> form("#plant-form", plant: @invalid_attrs)
             |> render_change() =~ "can&#39;t be blank"

      assert index_live
             |> form("#plant-form", plant: @update_attrs)
             |> render_submit()

      assert_patch(index_live, ~p"/plants")

      html = render(index_live)
      assert html =~ "Plant updated successfully"
      assert html =~ "some updated name"
    end

    test "can't remove a location by editing", %{conn: conn, plant: plant} do
      {:ok, index_live, _html} = live(conn, ~p"/plants")

      assert index_live
             |> element("#plants-#{plant.id} a", "Edit")
             |> render_click() =~ "Edit Plant"

      assert_patch(index_live, ~p"/plants/#{plant}/edit")

      remove_location = %{
        location_id: ""
      }

      assert index_live
             |> form("#plant-form", plant: remove_location)
             |> render_change() =~ "plants have to be planted somewhere"
    end

    test "deletes plant in listing", %{conn: conn, plant: plant} do
      {:ok, index_live, _html} = live(conn, ~p"/plants")

      assert index_live |> element("#plants-#{plant.id} a", "Delete") |> render_click()
      refute has_element?(index_live, "#plants-#{plant.id}")
    end
  end

  describe "Show" do
    setup [:create_plant]

    test "displays plant", %{conn: conn, plant: plant} do
      {:ok, _show_live, html} = live(conn, ~p"/plants/#{plant}")

      assert html =~ "Show Plant"
      assert html =~ plant.name
    end

    test "updates plant within modal", %{conn: conn, plant: plant} do
      {:ok, show_live, _html} = live(conn, ~p"/plants/#{plant}")

      assert show_live |> element("a", "Edit") |> render_click() =~
               "Edit Plant"

      assert_patch(show_live, ~p"/plants/#{plant}/show/edit")

      assert show_live
             |> form("#plant-form", plant: @invalid_attrs)
             |> render_change() =~ "can&#39;t be blank"

      assert show_live
             |> form("#plant-form", plant: @update_attrs)
             |> render_submit()

      assert_patch(show_live, ~p"/plants/#{plant}")

      html = render(show_live)
      assert html =~ "Plant updated successfully"
      assert html =~ "some updated name"
    end
  end
end
