defmodule GardenWeb.PlantLive.Show do
  use GardenWeb, :live_view

  alias Garden.Plants
  alias GardenWeb.PlantLive.EditForm
  alias GardenWeb.PlantLive.MoveForm

  @impl true
  def mount(_params, _session, socket) do
    {:ok, socket}
  end

  @impl true
  def handle_params(%{"id" => id}, _, socket) do
    {:noreply,
     socket
     |> assign(:page_title, page_title(socket.assigns.live_action))
     |> assign(:plant, plant(id))}
  end

  defp page_title(:show), do: "Show Plant"
  defp page_title(:edit), do: "Edit Plant"
  defp page_title(:move), do: "Transplant"

  @impl true
  def handle_info({EditForm, {:saved, plant}}, socket) do
    {:noreply, assign(socket, :plant, plant(plant.id))}
  end

  def handle_info({MoveForm, {:moved, plant_id}}, socket) do
    {:noreply, assign(socket, :plant, plant(plant_id))}
  end

  defp plant(id), do: Plants.get!(id, locations: :all, seed: true)
end
