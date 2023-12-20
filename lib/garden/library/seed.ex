defmodule Garden.Library.Seed do
  use Ecto.Schema
  import Ecto.Changeset

  schema "seeds" do
    field :name, :string
    field :front_image_id, :string, default: ""
    field :back_image_id, :string, default: ""
    field :year, :integer

    timestamps(type: :utc_datetime)
  end

  @doc false
  def changeset(seed, attrs, front_image_id \\ nil, back_image_id \\ nil) do
    seed
    |> cast(attrs, [:name, :year])
    |> validate_number(:year, greater_than_or_equal_to: 2020)
    |> front_image(front_image_id)
    |> back_image(back_image_id)
    |> validate_required([:name])
  end

  defp front_image(changeset, nil), do: changeset
  defp front_image(changeset, id), do: put_change(changeset, :front_image_id, id)

  defp back_image(changeset, nil), do: changeset
  defp back_image(changeset, id), do: put_change(changeset, :back_image_id, id)
end
