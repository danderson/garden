defmodule GardenWeb.PlantLive.CreateForm do
  use GardenWeb, :live_component

  alias Garden.Plants

  @impl true
  def render(assigns) do
    ~H"""
    <div>
      <.header>New plant</.header>

      <.simple_form
        class="px-4 pb-4"
        for={@form}
        id="new-plant-form"
        phx-target={@myself}
        phx-change="validate"
        phx-submit="save"
      >
        <.inputs_for :let={p} field={@form[:plant]}>
          <.input field={p[:name]} type="text" label="Name" data-1p-ignore />
          <.input field={p[:seed_id]} type="select" label="Seed" prompt="None" options={@seeds} />
        </.inputs_for>
        <.input
          field={@form[:location_id]}
          type="select"
          label="Location"
          prompt=""
          options={@locations}
        />

        <:actions>
          <.button phx-disable-with="Saving...">Add</.button>
        </:actions>
      </.simple_form>
    </div>
    """
  end

  @impl true
  def update(assigns, socket) do
    changeset = Plants.new_changeset(assigns[:initial_params])

    {:ok,
     socket
     |> assign(assigns)
     |> assign_locations()
     |> assign_seeds()
     |> assign_form(changeset)}
  end

  @impl true
  def handle_event("validate", %{"plant" => params}, socket) do
    changeset =
      Plants.new_changeset(params)
      |> Map.put(:action, :validate)

    {:noreply, assign_form(socket, changeset)}
  end

  def handle_event("save", %{"plant" => params}, socket) do
    case Plants.new(params) do
      {:ok, plant} ->
        notify_parent({:saved, plant})

        {:noreply,
         socket
         |> push_patch(to: socket.assigns.patch)}

      {:error, %Ecto.Changeset{} = changeset} ->
        {:noreply, assign_form(socket, changeset)}
    end
  end

  defp assign_locations(socket) do
    locs =
      Garden.Locations.list()
      |> Enum.map(fn loc -> {loc.name, loc.id} end)

    assign(socket, :locations, locs)
  end

  defp assign_seeds(socket) do
    seeds =
      Garden.Seeds.list()
      |> Enum.map(fn seed -> {seed.name, seed.id} end)

    assign(socket, :seeds, seeds)
  end

  defp assign_form(socket, %Ecto.Changeset{} = changeset) do
    assign(socket, :form, to_form(changeset, as: "plant"))
  end

  defp notify_parent(msg), do: send(self(), {__MODULE__, msg})
end
