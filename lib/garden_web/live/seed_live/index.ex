defmodule GardenWeb.SeedLive.Index do
  use GardenWeb, :live_view

  alias Garden.Library
  alias Garden.Library.Seed

  @impl true
  def mount(_params, _session, socket) do
    {:ok, stream(socket, :seeds, Library.list_seeds())}
  end

  @impl true
  def handle_params(params, _url, socket) do
    {:noreply, apply_action(socket, socket.assigns.live_action, params)}
  end

  defp apply_action(socket, :edit, %{"id" => id}) do
    socket
    |> assign(:page_title, "Edit Seed")
    |> assign(:seed, Library.get_seed!(id))
  end

  defp apply_action(socket, :new, _params) do
    socket
    |> assign(:page_title, "New Seed")
    |> assign(:seed, %Seed{})
  end

  defp apply_action(socket, :index, _params) do
    socket
    |> assign(:page_title, "Listing Seeds")
    |> assign(:seed, nil)
  end

  @impl true
  def handle_info({GardenWeb.SeedLive.FormComponent, {:saved, seed}}, socket) do
    {:noreply, stream_insert(socket, :seeds, seed)}
  end

  @impl true
  def handle_event("delete", %{"id" => id}, socket) do
    seed = Library.get_seed!(id)
    {:ok, _} = Library.delete_seed(seed)

    {:noreply, stream_delete(socket, :seeds, seed)}
  end
end
