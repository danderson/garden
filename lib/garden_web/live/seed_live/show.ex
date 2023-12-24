defmodule GardenWeb.SeedLive.Show do
  use GardenWeb, :live_view

  alias Garden.Seeds

  @impl true
  def mount(_params, _session, socket) do
    {:ok, socket}
  end

  @impl true
  def handle_params(%{"id" => id}, _, socket) do
    {:noreply,
     socket
     |> assign(:page_title, page_title(socket.assigns.live_action))
     |> assign(:seed, Seeds.get!(id, plants: true, locations: :current))}
  end

  def tribool_fields() do
    Enum.filter(Garden.Seeds.Seed.__schema__(:fields), fn f ->
      Garden.Seeds.Seed.__schema__(:type, f) == Garden.Tribool
    end)
  end

  def set_tribools(seed) do
    tribool_fields()
    |> Enum.map(fn f ->
      v = Map.get(seed, f)

      if v != nil do
        {tribool_label(f), v}
      else
        nil
      end
    end)
    |> Enum.reject(&is_nil/1)
  end

  def has_unset_tribools(seed) do
    tribool_fields() |> Enum.any?(fn f -> Map.get(seed, f) == nil end)
  end

  def unset_tribools(seed) do
    tribool_fields()
    |> Enum.filter(fn f -> Map.get(seed, f) == nil end)
    |> Enum.map(&tribool_label/1)
    |> Enum.map(&to_string/1)
    |> Enum.sort()
    |> Enum.join(", ")
  end

  def tribool_label(:edible), do: "Edible"
  def tribool_label(:needs_trellis), do: "Needs trellis"
  def tribool_label(:needs_bird_netting), do: "Needs bird netting"
  def tribool_label(:is_keto), do: "Keto"
  def tribool_label(:is_native), do: "Native"
  def tribool_label(:is_invasive), do: "Invasive"
  def tribool_label(:is_cover_crop), do: "Cover crop"
  def tribool_label(:grows_well_from_seed), do: "Grows well from seeds"
  def tribool_label(:is_bad_for_cats), do: "Bad for cats"
  def tribool_label(:is_deer_resistant), do: "Deer resistant"
  def tribool_label(_), do: "????"

  defp page_title(:show), do: "Show Seed"
  defp page_title(:edit), do: "Edit Seed"
end
