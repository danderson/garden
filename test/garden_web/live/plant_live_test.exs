defmodule GardenWeb.PlantLiveTest do
  use GardenWeb.ConnCase

  import Phoenix.LiveViewTest
  import Garden.PlantsFixtures
  import Garden.LocationsFixtures

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
        |> form("#new-plant-form", plant: %{plant: %{name: nil}})
        |> render_change()

      assert invalid_form_resp =~ "can&#39;t be blank"
      assert invalid_form_resp =~ "plants have to be planted somewhere"

      attrs = %{
        plant: %{
          name: "some name"
        },
        location_id: location.id
      }

      assert index_live
             |> form("#new-plant-form", plant: attrs)
             |> render_submit()

      assert_patch(index_live, ~p"/plants")

      html = render(index_live)
      assert html =~ "some name"
    end

    test "updates plant in listing", %{conn: conn, plant: plant} do
      {:ok, index_live, _html} = live(conn, ~p"/plants")

      assert index_live
             |> element("#plants-#{plant.id} a", "Edit")
             |> render_click() =~ plant.name

      assert_patch(index_live, ~p"/plants/#{plant}/edit")

      assert index_live
             |> form("#edit-plant-form", plant: %{name: nil})
             |> render_change() =~ "can&#39;t be blank"

      assert index_live
             |> form("#edit-plant-form", plant: %{name: "some updated name"})
             |> render_submit()

      assert_patch(index_live, ~p"/plants")

      html = render(index_live)
      assert html =~ "some updated name"
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
               plant.name

      assert_patch(show_live, ~p"/plants/#{plant}/show/edit")

      assert show_live
             |> form("#edit-plant-form", plant: %{name: nil})
             |> render_change() =~ "can&#39;t be blank"

      assert show_live
             |> form("#edit-plant-form", plant: %{name: "some updated name"})
             |> render_submit()

      assert_patch(show_live, ~p"/plants/#{plant}")

      html = render(show_live)
      assert html =~ "some updated name"
    end
  end
end
