defmodule Importer do
  def sqlite3(sql) do
    {out, 0} = System.cmd("sqlite3", [Path.join(__DIR__, "../../../garden/db.sqlite3"), ".mode json", sql])
    Jason.decode!(out, keys: :atoms)
  end

  def adjust(map, adjusters) do
    Enum.map(map, fn {k, v} ->
      {k, Map.get(adjusters, k, &(&1)).(v)}
    end)
    |> Map.new
  end

  def to_id_table(list_of_maps, field \\ :id) do
    list_of_maps
    |> Enum.map(fn %{^field => k} = m ->
      {k, m}
    end)
    |> Map.new
  end

  def adjust_enum(type, field, v) do
    Ecto.Enum.mappings(type, field)
    |> Enum.map(fn {k, v} -> {v, k} end)
    |> Map.new
    |> Map.fetch!(v)
  end

  def rename(map, from_field, to_field, adjuster) do
    adj = Map.fetch!(map, from_field)    |> adjuster.()
    map
    |> Map.drop([from_field])
    |> Map.put(to_field, adj)
  end

  def adjust_tribool(nil), do: nil
  def adjust_tribool(0), do: false
  def adjust_tribool(1), do: true

  def families() do
    sqlite3("select * from herbarium_family")
    |> Enum.map(fn %{id: k, name: v} -> {k, String.to_atom(v)} end)
    |> Map.new()
  end

  def adjust_type(v), do: adjust_enum(Garden.Seeds.Seed, :type, v)
  def adjust_lifespan(v), do: adjust_enum(Garden.Seeds.Seed, :lifespan, v)

  def plants(families) do
    sqlite3("select * from herbarium_plant")
    |> Enum.map(fn plant ->
      plant
      |> rename(:family_id, :family, &(families[&1]))
      |> adjust(%{
        edible: &adjust_tribool/1,
        needs_trellis: &adjust_tribool/1,
        needs_bird_netting: &adjust_tribool/1,
        is_keto: &adjust_tribool/1,
        native: &adjust_tribool/1,
        invasive: &adjust_tribool/1,
        is_cover: &adjust_tribool/1,
        grow_from_seed: &adjust_tribool/1,
        bad_for_cats: &adjust_tribool/1,
        deer_resistant: &adjust_tribool/1,
        type: &adjust_type/1,
        lifespan: &adjust_lifespan/1,
      })
    end)
    |> to_id_table
  end

  def adjust_date(nil), do: nil
  def adjust_date(s), do: Date.from_iso8601!(s)

  def boxcontent() do
    sqlite3("select * from boxinventory_boxcontent")
    |> Enum.map(fn plant ->
      plant
      |> adjust(%{
        planted: &adjust_date/1,
        removed: &adjust_date/1,
      })
      |> Map.drop([:latin_name])
      |> then(fn %{box_id: id} = plant ->
        {id, Map.drop(plant, [:box_id])}
      end)
    end)
    |> Enum.group_by(fn {k, _} -> k end, fn {_, v} -> v end)
  end

  def boxes(content) do
    sqlite3("select * from boxinventory_box")
    |> Enum.map(fn location ->
      location
      |> adjust(%{
        qr_applied: &adjust_tribool/1,
        want_qr: &adjust_tribool/1,
      })
      |> Map.put(:contents, content[location.id])
      |> rename(:id, :qr_id, &(&1))
    end)
    |> to_id_table(:qr_id)
  end

  def run() do
    #families = families()
    #plants = plants(families)
    content = boxcontent()
    boxes(content)
  end
end
