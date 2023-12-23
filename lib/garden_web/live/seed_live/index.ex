defmodule GardenWeb.SeedLive.Index do
  use GardenWeb, :live_view

  alias Garden.Seeds

  @impl true
  def mount(_params, _session, socket) do
    {:ok, stream(socket, :seeds, Seeds.list())}
  end

  @impl true
  def handle_params(params, _url, socket) do
    {:noreply, apply_action(socket, socket.assigns.live_action, params)}
  end

  defp apply_action(socket, :edit, %{"id" => id}) do
    socket
    |> assign(:page_title, "Edit Seed")
    |> assign(:seed, Seeds.get!(id))
  end

  defp apply_action(socket, :new, _params) do
    socket
    |> assign(:page_title, "New Seed")
    |> assign(:seed, %Garden.Seeds.Seed{})
  end

  defp apply_action(socket, :index, _params) do
    socket
    |> assign(:page_title, "Listing Seeds")
    |> assign(:seed, nil)
  end

  @impl true
  def handle_info({GardenWeb.SeedLive.FormComponent, {:saved, seed}}, socket) do
    {:noreply, stream_insert(socket, :seeds, Seeds.get!(seed))}
  end
end
