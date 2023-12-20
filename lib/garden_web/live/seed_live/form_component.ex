defmodule GardenWeb.SeedLive.FormComponent do
  use GardenWeb, :live_component

  alias Garden.Library

  @impl true
  def render(assigns) do
    ~H"""
    <div>
      <.header>
        <%= @title %>
      </.header>

      <.simple_form
        for={@form}
        id="seed-form"
        phx-target={@myself}
        phx-change="validate"
        phx-submit="save"
      >
        <.input field={@form[:name]} type="text" label="Name" placeholder="Big potats" data-1p-ignore />
        <.input field={@form[:year]} type="number" label="Packed for" min="2020" max="2099" step="1" />

        <.photo_upload
          upload={@uploads.front_photo}
          label="Front photo"
          existing_id={@seed.front_image_id}
        />
        <.photo_upload
          upload={@uploads.back_photo}
          label="Back photo"
          existing_id={@seed.back_image_id}
        />

        <:actions>
          <.button phx-disable-with="Saving...">Save Seed</.button>
        </:actions>
      </.simple_form>
    </div>
    """
  end

  attr :upload, :any, required: true, doc: "the upload object for the photo"
  attr :label, :string, required: true, doc: "the label for the photo"
  attr :existing_id, :string, required: true, doc: "the existing database ID for the photo"

  def photo_upload(assigns) do
    ~H"""
    <div>
      <div class="block text-sm font-semibold leading-6 text-zinc-800"><%= @label %></div>
      <%= if @upload.entries == [] do %>
        <%= if @existing_id != "" do %>
          <img src={Images.url(@existing_id, :medium)} />
        <% end %>
      <% else %>
        <.live_img_preview entry={List.first(@upload.entries)} />
        <%= for err <- upload_errors(@upload, List.first(@upload.entries)) do %>
          <p class="alert alert-danger"><%= error_to_string(err) %></p>
        <% end %>
        <button
          type="button"
          phx-click="cancel-upload"
          phx-value-upload={@upload.name}
          phx-value-ref={List.first(@upload.entries).ref}
        >
          <.icon name="hero-trash" />
        </button>
      <% end %>
      <label for={@upload.ref}>
        <.live_file_input upload={@upload} style="display: none" />
        <.icon name="hero-camera" />
      </label>
    </div>
    """
  end

  @impl true
  def mount(socket) do
    socket =
      socket
      |> allow_upload(:front_photo, accept: ["image/jpeg"])
      |> allow_upload(:back_photo, accept: ["image/jpeg"])

    {:ok, socket}
  end

  @impl true
  def update(%{seed: seed} = assigns, socket) do
    changeset = Library.change_seed(seed)

    {:ok,
     socket
     |> assign(assigns)
     |> assign_form(changeset)}
  end

  @impl true
  def handle_event("cancel-upload", %{"upload" => upload, "ref" => ref}, socket) do
    {:noreply, cancel_upload(socket, upload, ref)}
  end

  @impl true
  def handle_event("validate", %{"seed" => seed_params}, socket) do
    changeset =
      socket.assigns.seed
      |> Library.change_seed(seed_params)
      |> Map.put(:action, :validate)

    {:noreply, assign_form(socket, changeset)}
  end

  def handle_event("save", %{"seed" => seed_params}, socket) do
    save_seed(socket, socket.assigns.action, seed_params)
  end

  defp save_seed(socket, :edit, seed_params) do
    case Library.update_seed(
           socket.assigns.seed,
           seed_params,
           uploaded_photo(socket, :front_photo),
           uploaded_photo(socket, :back_photo)
         ) do
      {:ok, seed} ->
        notify_parent({:saved, seed})

        {:noreply,
         socket
         |> put_flash(:info, "Seed updated successfully")
         |> push_patch(to: socket.assigns.patch)}

      {:error, %Ecto.Changeset{} = changeset} ->
        {:noreply, assign_form(socket, changeset)}
    end
  end

  defp save_seed(socket, :new, seed_params) do
    case Library.create_seed(
           seed_params,
           uploaded_photo(socket, :front_photo),
           uploaded_photo(socket, :back_photo)
         ) do
      {:ok, seed} ->
        notify_parent({:saved, seed})

        {:noreply,
         socket
         |> put_flash(:info, "Seed created successfully")
         |> push_patch(to: socket.assigns.patch)}

      {:error, %Ecto.Changeset{} = changeset} ->
        {:noreply, assign_form(socket, changeset)}
    end
  end

  defp uploaded_photo(socket, name) do
    uploads =
      consume_uploaded_entries(socket, name, fn %{path: path}, _entry ->
        {:ok, Images.store(path)}
      end)

    case uploads do
      [] -> nil
      [path] -> path
    end
  end

  defp assign_form(socket, %Ecto.Changeset{} = changeset) do
    assign(socket, :form, to_form(changeset))
  end

  defp notify_parent(msg), do: send(self(), {__MODULE__, msg})

  defp error_to_string(:too_many_files), do: "Only 1 file allowed"
  defp error_to_string(:not_accepted), do: "Unacceptable file type"
end
