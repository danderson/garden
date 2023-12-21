defmodule GardenWeb.PageControllerTest do
  use GardenWeb.ConnCase

  test "GET /", %{conn: conn} do
    conn = get(conn, ~p"/")
    assert html_response(conn, 200) =~ "Welcome to the garden"
  end
end
