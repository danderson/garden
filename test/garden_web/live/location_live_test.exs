defmodule GardenWeb.LocationLiveTest do
  use GardenWeb.ConnCase

  import Phoenix.LiveViewTest
  import Garden.LocationsFixtures

  @create_attrs %{name: "some name"}
  @update_attrs %{name: "some updated name"}
  @invalid_attrs %{name: nil}

  defp create_location(_) do
    location = location_fixture()
    %{location: location}
  end

  describe "Index" do
    setup [:create_location]

    test "lists all locations", %{conn: conn, location: location} do
      {:ok, _index_live, html} = live(conn, ~p"/locations")

      assert html =~ "Listing Locations"
      assert html =~ location.name
    end

    test "saves new location", %{conn: conn} do
      {:ok, index_live, _html} = live(conn, ~p"/locations")

      assert index_live
             |> element("a", "Add")
             |> render_click() =~ "Add"

      assert_patch(index_live, ~p"/locations/new")

      assert index_live
             |> form("#location-form", location: @invalid_attrs)
             |> render_change() =~ "can&#39;t be blank"

      assert index_live
             |> form("#location-form", location: @create_attrs)
             |> render_submit()

      assert_patch(index_live, ~p"/locations")

      html = render(index_live)
      assert html =~ "Location created successfully"
      assert html =~ "some name"
    end

    test "updates location in listing", %{conn: conn, location: location} do
      {:ok, index_live, _html} = live(conn, ~p"/locations")

      assert index_live |> element("#locations-#{location.id} a", "Edit") |> render_click() =~
               "Edit Location"

      assert_patch(index_live, ~p"/locations/#{location.id}/edit")

      assert index_live
             |> form("#location-form", location: @invalid_attrs)
             |> render_change() =~ "can&#39;t be blank"

      assert index_live
             |> form("#location-form", location: @update_attrs)
             |> render_submit()

      assert_patch(index_live, ~p"/locations")

      html = render(index_live)
      assert html =~ "Location updated successfully"
      assert html =~ "some updated name"
    end

    test "plant from listing", %{conn: conn, location: location} do
      {:ok, index_live, _html} = live(conn, ~p"/locations")

      {:ok, plant_live, _html} =
        index_live
        |> element("#locations-#{location.id} a", "Plant")
        |> render_click()
        |> follow_redirect(conn)

      assert plant_live
             |> form("#new-plant-form",
               plant: %{
                 plant: %{
                   name: "plant from location"
                 },
                 location_id: location.id
               }
             )
             |> render_submit()

      assert_patch(plant_live, ~p"/plants")

      assert render(plant_live) =~ "plant from location"

      new_plant = Garden.Plants.get!(1, locations: :current)
      assert new_plant.current_location == location
    end
  end

  describe "Show" do
    setup [:create_location]

    test "displays location", %{conn: conn, location: location} do
      {:ok, _show_live, html} = live(conn, ~p"/locations/#{location.id}")

      assert html =~ "Show Location"
      assert html =~ location.name
    end

    test "updates location within modal", %{conn: conn, location: location} do
      {:ok, show_live, _html} = live(conn, ~p"/locations/#{location.id}")

      assert show_live |> element("a", "Edit") |> render_click() =~
               "Edit Location"

      assert_patch(show_live, ~p"/locations/#{location.id}/show/edit")

      assert show_live
             |> form("#location-form", location: @invalid_attrs)
             |> render_change() =~ "can&#39;t be blank"

      assert show_live
             |> form("#location-form", location: @update_attrs)
             |> render_submit()

      assert_patch(show_live, ~p"/locations/#{location.id}")

      html = render(show_live)
      assert html =~ "Location updated successfully"
      assert html =~ "some updated name"
    end

    test "plant here", %{conn: conn, location: location} do
      {:ok, index_live, _html} = live(conn, ~p"/locations/#{location.id}")

      {:ok, plant_live, _html} =
        index_live
        |> element("main a", "Plant")
        |> render_click()
        |> follow_redirect(conn)

      assert plant_live
             |> form("#new-plant-form",
               plant: %{
                 plant: %{
                   name: "plant from location"
                 },
                 location_id: location.id
               }
             )
             |> render_submit()

      assert_patch(plant_live, ~p"/plants")

      assert render(plant_live) =~ "plant from location"

      new_plant = Garden.Plants.get!(1, locations: :current)
      assert new_plant.current_location == location
    end
  end
end
