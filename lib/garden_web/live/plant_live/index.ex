defmodule GardenWeb.PlantLive.Index do
  use GardenWeb, :live_view

  alias Garden.Plants

  @impl true
  def mount(_params, _session, socket) do
    {:ok,
     socket
     |> stream(:plants, Plants.list_plants())
     |> assign(:form_initial_params, %{})}
  end

  @impl true
  def handle_params(params, _url, socket) do
    {:noreply, apply_action(socket, socket.assigns.live_action, params)}
  end

  defp apply_action(socket, :edit, %{"id" => id}) do
    socket
    |> assign(:page_title, "Edit Plant")
    |> assign(:plant, Plants.get_plant!(id))
  end

  defp apply_action(socket, :new, params) do
    {init, _} = new_plant_initial_params(params)

    socket
    |> assign(:page_title, "Add Plant")
    |> assign(:plant, Plants.new_plant())
    |> assign(:form_initial_params, init)
  end

  defp apply_action(socket, :index, _params) do
    socket
    |> assign(:page_title, "Plants")
    |> assign(:plant, nil)
  end

  @impl true
  def handle_info({GardenWeb.PlantLive.FormComponent, {:saved, plant}}, socket) do
    {:noreply, stream_insert(socket, :plants, Plants.expand_plant(plant))}
  end

  @impl true
  def handle_event("delete", %{"id" => id}, socket) do
    plant = Plants.get_plant!(id)
    {:ok, _} = Plants.delete_plant(plant)

    {:noreply, stream_delete(socket, :plants, plant)}
  end

  defp new_plant_initial_params(params) do
    {location, params} = Map.pop(params, "location_id")
    location = if location != nil, do: %{location_id: String.to_integer(location)}, else: %{}

    {seed, params} = Map.pop(params, "seed_id")
    seed = if seed != nil, do: %{seed_id: String.to_integer(seed)}, else: %{}

    {Map.merge(location, seed), params}
  end
end
