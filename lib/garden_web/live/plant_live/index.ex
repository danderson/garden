defmodule GardenWeb.PlantLive.Index do
  use GardenWeb, :live_view

  alias Garden.Plants
  alias GardenWeb.PlantLive.CreateForm
  alias GardenWeb.PlantLive.EditForm

  @impl true
  def mount(_params, _session, socket) do
    {:ok,
     socket
     |> stream(:plants, Plants.list(locations: :current))
     |> assign(:form_initial_params, %{})}
  end

  @impl true
  def handle_params(params, _url, socket) do
    {:noreply, apply_action(socket, socket.assigns.live_action, params)}
  end

  defp apply_action(socket, :edit, %{"id" => id}) do
    socket
    |> assign(:page_title, "Edit Plant")
    |> assign(:plant, Plants.get!(id, locations: :current, seed: true))
  end

  defp apply_action(socket, :new, params) do
    plant_params = Map.take(params, [:location_id, :seed_id])

    socket
    |> assign(:page_title, "Add Plant")
    |> assign(:form_initial_params, plant_params)
  end

  defp apply_action(socket, :index, _params) do
    socket
    |> assign(:page_title, "Plants")
    |> assign(:plant, nil)
  end

  @impl true
  def handle_info({src, {:saved, plant}}, socket) when src in [CreateForm, EditForm] do
    {:noreply, stream_insert(socket, :plants, Plants.get!(plant.id, locations: :current))}
  end
end
