defmodule Images do
  def base_url, do: "/user_images"

  def images_dir, do: Application.fetch_env!(:garden, :images_dir)

  def store(kind, src_path) do
    File.mkdir_p!(disk_dir(kind))

    id = generate_id()

    full =
      src_path
      |> Image.open!()
      |> Image.autorotate!()
      |> Image.thumbnail!("1600x1600", resize: :down)
      |> Image.write!(disk(kind, id, :full))

    full
    |> Image.thumbnail!("400x600", resize: :down)
    |> Image.write!(disk(kind, id, :medium))

    full
    |> Image.thumbnail!("128x128", resize: :down)
    |> Image.write!(disk(kind, id, :thumbnail))

    id
  end

  def delete(_kind, nil), do: nil

  def delete(kind, id) do
    [:full, :medium, :thumbnail]
    |> Enum.each(fn size ->
      path = disk(kind, id, size)
      if File.exists?(path), do: File.rm!(path)
    end)

    nil
  end

  def url(kind, id, size), do: Path.join([base_url(), to_string(kind), filename(id, size)])
  defp disk_dir(kind), do: Path.join(images_dir(), to_string(kind))
  defp disk(kind, id, size), do: Path.join(disk_dir(kind), filename(id, size))

  defp filename(id, :full), do: "#{id}.jpg"
  defp filename(id, :medium), do: "#{id}_med.jpg"
  defp filename(id, :thumbnail), do: "#{id}_thumb.jpg"

  defp generate_id, do: Ecto.UUID.generate()
end
