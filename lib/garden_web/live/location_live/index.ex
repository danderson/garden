defmodule GardenWeb.LocationLive.Index do
  use GardenWeb, :live_view

  alias Garden.Locations
  alias Garden.Locations.Location

  @impl true
  def mount(_params, _session, socket) do
    {:ok, stream(socket, :locations, Locations.list(plants: :current))}
  end

  @impl true
  def handle_params(params, _url, socket) do
    {:noreply, apply_action(socket, socket.assigns.live_action, params)}
  end

  defp apply_action(socket, :edit, %{"id" => id}) do
    socket
    |> assign(:page_title, "Edit Location")
    |> assign(:location, Locations.get!(id))
  end

  defp apply_action(socket, :new, params) do
    location_params = %{
      qr_id: params["qr_id"]
    }

    socket
    |> assign(:page_title, "New Location")
    |> assign(:location, %Location{})
    |> assign(:form_initial_params, location_params)
  end

  defp apply_action(socket, :index, _params) do
    socket
    |> assign(:page_title, "Listing Locations")
    |> assign(:location, nil)
  end

  @impl true
  def handle_info({GardenWeb.LocationLive.FormComponent, {:saved, location}}, socket) do
    {:noreply, stream_insert(socket, :locations, Locations.get!(location.id, plants: :current))}
  end
end
