defmodule GardenWeb.PlantLive.MoveForm do
  use GardenWeb, :live_component

  alias Garden.Plants

  @impl true
  def render(assigns) do
    ~H"""
    <div>
      <.header>Transplant</.header>

      <.simple_form
        for={@form}
        id="move-plant-form"
        phx-target={@myself}
        phx-change="validate"
        phx-submit="save"
      >
        <.input
          field={@form[:location_id]}
          type="select"
          label="Location"
          prompt=""
          options={@locations}
        />

        <:actions>
          <.button phx-disable-with="Saving...">Transplant</.button>
        </:actions>
      </.simple_form>
    </div>
    """
  end

  defmodule Form do
    use Ecto.Schema
    import Ecto.Changeset

    embedded_schema do
      field :location_id, :id
    end

    def changeset(location_id, attrs \\ %{}) do
      %Form{location_id: location_id}
      |> cast(attrs, [:location_id])
      |> validate_required([:location_id])
    end
  end

  @impl true
  def update(%{location_id: id} = assigns, socket) do
    changeset = Form.changeset(id)

    {:ok,
     socket
     |> assign(assigns)
     |> assign_locations()
     |> assign_form(changeset)}
  end

  @impl true
  def handle_event("validate", %{"form" => params}, socket) do
    changeset =
      Form.changeset(socket.assigns.location_id, params)
      |> Map.put(:action, :validate)

    {:noreply, assign_form(socket, changeset)}
  end

  def handle_event("save", %{"form" => params}, socket) do
    change =
      Form.changeset(socket.assigns.location_id, params)

    case Ecto.Changeset.apply_action(change, :move) do
      {:ok, data} ->
        if Ecto.Changeset.changed?(change, :location_id) do
          Plants.move(socket.assigns.plant.id, data.location_id)
        end

        notify_parent({:moved, socket.assigns.plant.id})
        {:noreply, push_patch(socket, to: socket.assigns.patch)}

      {:error, changeset} ->
        {:noreply, assign_form(socket, changeset)}
    end
  end

  defp assign_locations(socket) do
    locs =
      Garden.Locations.list()
      |> Enum.map(fn loc -> {loc.name, loc.id} end)

    assign(socket, :locations, locs)
  end

  defp assign_form(socket, %Ecto.Changeset{} = changeset) do
    assign(socket, :form, to_form(changeset))
  end

  defp notify_parent(msg), do: send(self(), {__MODULE__, msg})
end
