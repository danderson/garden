defmodule Garden.Library.Seed do
  use Ecto.Schema
  import Ecto.Changeset

  schema "seeds" do
    field :name, :string
    field :image_id, :string, default: ""

    timestamps(type: :utc_datetime)
  end

  @doc false
  def changeset(seed, attrs, image_id \\ nil) do
    seed
    |> cast(attrs, [:name])
    |> maybe_update_image(image_id)
    |> validate_required([:name])
  end

  defp maybe_update_image(changeset, nil), do: changeset
  defp maybe_update_image(changeset, id), do: put_change(changeset, :image_id, id)
end
