defmodule Images do
  def base_url, do: "/user_images"

  def images_dir, do: Application.fetch_env!(:garden, :images_dir)

  def store(src_path) do
    File.mkdir_p!(images_dir())

    id = generate_id()

    full =
      src_path
      |> Image.open!()
      |> Image.autorotate!()
      |> Image.thumbnail!("1600x1600", resize: :down)
      |> Image.write!(disk(id, :full))

    full
    |> Image.thumbnail!("400x600", resize: :down)
    |> Image.write!(disk(id, :medium))

    full
    |> Image.thumbnail!("128x128", resize: :down)
    |> Image.write!(disk(id, :thumbnail))

    id
  end

  def delete(nil), do: nil

  def delete(id) do
    [:full, :medium, :thumbnail]
    |> Enum.each(fn size ->
      path = disk(id, size)
      if File.exists?(path), do: File.rm!(path)
    end)

    nil
  end

  def url(id, size), do: Path.join(base_url(), filename(id, size))
  defp disk(id, size), do: Path.join(images_dir(), filename(id, size))

  defp filename(id, :full), do: "#{id}.jpg"
  defp filename(id, :medium), do: "#{id}_med.jpg"
  defp filename(id, :thumbnail), do: "#{id}_thumb.jpg"

  defp generate_id, do: Ecto.UUID.generate()
end
