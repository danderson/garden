defmodule GardenWeb.PlantLive.EditForm do
  use GardenWeb, :live_component

  alias Garden.Plants

  @impl true
  def render(assigns) do
    ~H"""
    <div>
      <.header><%= @plant.name %></.header>

      <.simple_form
        class="px-4 pb-4"
        for={@form}
        id="edit-plant-form"
        phx-target={@myself}
        phx-change="validate"
        phx-submit="save"
      >
        <.input field={@form[:name]} type="text" label="Name" data-1p-ignore />
        <.input field={@form[:seed_id]} type="select" label="Seed" prompt="None" options={@seeds} />

        <:actions>
          <.button phx-disable-with="Saving...">Save</.button>
        </:actions>
      </.simple_form>
    </div>
    """
  end

  @impl true
  def update(%{plant: plant} = assigns, socket) do
    changeset = Plants.edit_changeset(plant)

    {:ok,
     socket
     |> assign(assigns)
     |> assign_seeds
     |> assign_form(changeset)}
  end

  @impl true
  def handle_event("validate", %{"plant" => params}, socket) do
    changeset =
      socket.assigns.plant
      |> Plants.edit_changeset(params)
      |> Map.put(:action, :validate)

    {:noreply, assign_form(socket, changeset)}
  end

  def handle_event("save", %{"plant" => params}, socket) do
    case Plants.edit(socket.assigns.plant, params) do
      {:ok, plant} ->
        notify_parent({:saved, plant})

        {:noreply,
         socket
         |> put_flash(:info, "Plant updated")
         |> push_patch(to: socket.assigns.patch)}

      {:error, %Ecto.Changeset{} = changeset} ->
        {:noreply, assign_form(socket, changeset)}
    end
  end

  defp assign_seeds(socket) do
    seeds =
      Garden.Seeds.list()
      |> Enum.map(fn seed -> {seed.name, seed.id} end)

    assign(socket, :seeds, seeds)
  end

  defp assign_form(socket, %Ecto.Changeset{} = changeset) do
    assign(socket, :form, to_form(changeset))
  end

  defp notify_parent(msg), do: send(self(), {__MODULE__, msg})
end
