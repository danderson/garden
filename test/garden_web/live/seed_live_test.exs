defmodule GardenWeb.SeedLiveTest do
  use GardenWeb.ConnCase

  import Phoenix.LiveViewTest
  import Garden.SeedsFixtures
  import Garden.LocationsFixtures

  @create_attrs %{name: "some name"}
  @update_attrs %{name: "some updated name"}
  @invalid_attrs %{name: nil}

  defp create_seed(_) do
    location = location_fixture()
    seed = seed_fixture()
    %{location: location, seed: seed}
  end

  describe "Index" do
    setup [:create_seed]

    test "lists all seeds", %{conn: conn, seed: seed} do
      {:ok, _index_live, html} = live(conn, ~p"/seeds")

      assert html =~ "Listing Seeds"
      assert html =~ seed.name
    end

    test "saves new seed", %{conn: conn} do
      {:ok, index_live, _html} = live(conn, ~p"/seeds")

      assert index_live |> element("a", "Add") |> render_click() =~
               "Add"

      assert_patch(index_live, ~p"/seeds/new")

      assert index_live
             |> form("#seed-form", seed: @invalid_attrs)
             |> render_change() =~ "can&#39;t be blank"

      assert index_live
             |> form("#seed-form", seed: @create_attrs)
             |> render_submit()

      assert_patch(index_live, ~p"/seeds")

      html = render(index_live)
      assert html =~ "Seed created successfully"
      assert html =~ "some name"
    end

    test "updates seed in listing", %{conn: conn, seed: seed} do
      {:ok, index_live, _html} = live(conn, ~p"/seeds")

      assert index_live |> element("#seeds-#{seed.id} a", "Edit") |> render_click() =~
               "Edit Seed"

      assert_patch(index_live, ~p"/seeds/#{seed}/edit")

      assert index_live
             |> form("#seed-form", seed: @invalid_attrs)
             |> render_change() =~ "can&#39;t be blank"

      assert index_live
             |> form("#seed-form", seed: @update_attrs)
             |> render_submit()

      assert_patch(index_live, ~p"/seeds")

      html = render(index_live)
      assert html =~ "Seed updated successfully"
      assert html =~ "some updated name"
    end

    test "deletes seed in listing", %{conn: conn, seed: seed} do
      {:ok, index_live, _html} = live(conn, ~p"/seeds")

      assert index_live |> element("#seeds-#{seed.id} a", "Delete") |> render_click()
      refute has_element?(index_live, "#seeds-#{seed.id}")
    end

    test "plant seed from listing", %{conn: conn, seed: seed, location: location} do
      {:ok, index_live, _html} = live(conn, ~p"/seeds")

      {:ok, plant_live, _html} =
        index_live
        |> element("#seeds-#{seed.id} a", "Plant")
        |> render_click()
        |> follow_redirect(conn)

      assert plant_live
             |> form("#plant-form",
               plant: %{
                 location_id: location.id,
                 name: "plant from seed"
               }
             )
             |> render_submit()

      assert_patch(plant_live, ~p"/plants")

      assert render(plant_live) =~ "plant from seed"

      new_plant = Garden.Plants.get_plant!(1)
      assert new_plant.seed_id == seed.id
    end
  end

  describe "Show" do
    setup [:create_seed]

    test "displays seed", %{conn: conn, seed: seed} do
      {:ok, _show_live, html} = live(conn, ~p"/seeds/#{seed}")

      assert html =~ "Show Seed"
      assert html =~ seed.name
    end

    test "updates seed within modal", %{conn: conn, seed: seed} do
      {:ok, show_live, _html} = live(conn, ~p"/seeds/#{seed}")

      assert show_live |> element("a", "Edit") |> render_click() =~
               "Edit Seed"

      assert_patch(show_live, ~p"/seeds/#{seed}/show/edit")

      assert show_live
             |> form("#seed-form", seed: @invalid_attrs)
             |> render_change() =~ "can&#39;t be blank"

      assert show_live
             |> form("#seed-form", seed: @update_attrs)
             |> render_submit()

      assert_patch(show_live, ~p"/seeds/#{seed}")

      html = render(show_live)
      assert html =~ "Seed updated successfully"
      assert html =~ "some updated name"
    end

    test "plant seed", %{conn: conn, seed: seed, location: location} do
      {:ok, index_live, _html} = live(conn, ~p"/seeds/#{seed}")

      {:ok, plant_live, _html} =
        index_live
        |> element("main a", "Plant")
        |> render_click()
        |> follow_redirect(conn)

      assert plant_live
             |> form("#plant-form",
               plant: %{
                 location_id: location.id,
                 name: "plant from seed"
               }
             )
             |> render_submit()

      assert_patch(plant_live, ~p"/plants")

      assert render(plant_live) =~ "plant from seed"

      new_plant = Garden.Plants.get_plant!(1)
      assert new_plant.seed_id == seed.id
    end
  end
end
