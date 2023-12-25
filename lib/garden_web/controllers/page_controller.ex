defmodule GardenWeb.PageController do
  use GardenWeb, :controller

  alias Garden.Locations

  def home(conn, _params) do
    # The home page is often custom made,
    # so skip the default app layout.
    render(conn, :home)
  end

  def legacy_qr(conn, %{"id" => id}) do
    case Locations.get_from_qr(id) do
      nil -> GardenWeb.ErrorHTML.render("404.html", %{})
      location -> redirect(conn, to: ~p"/locations/#{location.id}")
    end
  end
end
