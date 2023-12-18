defmodule GardenWeb.PlantLive.FormComponent do
  use GardenWeb, :live_component

  alias Garden.Plants

  @impl true
  def render(assigns) do
    ~H"""
    <div>
      <.header>
        <%= @title %>
        <:subtitle>Use this form to manage plant records in your database.</:subtitle>
      </.header>

      <.simple_form
        for={@form}
        id="plant-form"
        phx-target={@myself}
        phx-change="validate"
        phx-submit="save"
      >
        <.input field={@form[:name]} type="text" label="Name" />
        <:actions>
          <.button phx-disable-with="Saving...">Save Plant</.button>
        </:actions>
      </.simple_form>
    </div>
    """
  end

  @impl true
  def update(%{plant: plant} = assigns, socket) do
    changeset = Plants.change_plant(plant)

    {:ok,
     socket
     |> assign(assigns)
     |> assign_form(changeset)}
  end

  @impl true
  def handle_event("validate", %{"plant" => plant_params}, socket) do
    changeset =
      socket.assigns.plant
      |> Plants.change_plant(plant_params)
      |> Map.put(:action, :validate)

    {:noreply, assign_form(socket, changeset)}
  end

  def handle_event("save", %{"plant" => plant_params}, socket) do
    save_plant(socket, socket.assigns.action, plant_params)
  end

  defp save_plant(socket, :edit, plant_params) do
    case Plants.update_plant(socket.assigns.plant, plant_params) do
      {:ok, plant} ->
        notify_parent({:saved, plant})

        {:noreply,
         socket
         |> put_flash(:info, "Plant updated successfully")
         |> push_patch(to: socket.assigns.patch)}

      {:error, %Ecto.Changeset{} = changeset} ->
        {:noreply, assign_form(socket, changeset)}
    end
  end

  defp save_plant(socket, :new, plant_params) do
    case Plants.create_plant(plant_params) do
      {:ok, plant} ->
        notify_parent({:saved, plant})

        {:noreply,
         socket
         |> put_flash(:info, "Plant created successfully")
         |> push_patch(to: socket.assigns.patch)}

      {:error, %Ecto.Changeset{} = changeset} ->
        {:noreply, assign_form(socket, changeset)}
    end
  end

  defp assign_form(socket, %Ecto.Changeset{} = changeset) do
    assign(socket, :form, to_form(changeset))
  end

  defp notify_parent(msg), do: send(self(), {__MODULE__, msg})
end
