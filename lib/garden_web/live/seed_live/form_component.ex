defmodule GardenWeb.SeedLive.FormComponent do
  use GardenWeb, :live_component

  alias Garden.Seeds

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
        <.input
          field={@form[:year]}
          type="number"
          label="Packed for"
          min="2020"
          max="2099"
          step="1"
          placeholder="Year"
        />

        <.photo_upload
          upload={@uploads.front_image}
          label="Front photo"
          existing_id={@seed.front_image_id}
          target={@myself}
        />
        <.photo_upload
          upload={@uploads.back_image}
          label="Back photo"
          existing_id={@seed.back_image_id}
          target={@myself}
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
  attr :target, :any, required: true, doc: "target of events"

  def photo_upload(assigns) do
    ~H"""
    <div>
      <div class="block text-sm font-semibold leading-6 text-zinc-800"><%= @label %></div>
      <%= if @upload.entries == [] do %>
        <img :if={@existing_id} src={Seeds.seed_image(@existing_id, :medium)} />
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
          phx-target={@target}
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
      |> allow_upload(:front_image, accept: ["image/jpeg"])
      |> allow_upload(:back_image, accept: ["image/jpeg"])

    {:ok, socket}
  end

  def update(%{id: :new} = assigns, socket) do
    {:ok,
     socket
     |> assign(assigns)
     |> assign_form(Seeds.upsert_changeset(%Garden.Seeds.Seed{}))}
  end

  @impl true
  def update(%{seed: seed} = assigns, socket) do
    changeset = Seeds.upsert_changeset(seed)

    {:ok,
     socket
     |> assign(assigns)
     |> assign_form(changeset)}
  end

  @impl true
  def handle_event("cancel-upload", %{"upload" => upload, "ref" => ref}, socket) do
    {:noreply, cancel_upload(socket, String.to_existing_atom(upload), ref)}
  end

  @impl true
  def handle_event("validate", %{"seed" => seed_params}, socket) do
    changeset =
      socket.assigns.seed
      |> Seeds.upsert_changeset(seed_params)
      |> Map.put(:action, :validate)

    {:noreply, assign_form(socket, changeset)}
  end

  def handle_event("save", %{"seed" => seed_params}, socket) do
    save_seed(socket, socket.assigns.action, seed_params)
  end

  defp save_seed(socket, :edit, seed_params) do
    case Seeds.edit(
           socket.assigns.seed,
           seed_params,
           collect_images(socket)
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
    case Seeds.new(
           seed_params,
           collect_images(socket)
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

  defp collect_images(socket) do
    %{
      front_image_id: collect_image(socket, :front_image),
      back_image_id: collect_image(socket, :back_image)
    }
    |> Enum.filter(fn {_k, v} -> v end)
    |> Map.new()
  end

  defp collect_image(socket, name) do
    uploads =
      consume_uploaded_entries(socket, name, fn %{path: path}, _entry ->
        {:ok, Seeds.store_seed_image(path)}
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
