defmodule Garden.Seeds.Seed do
  use Ecto.Schema
  import Ecto.Changeset

  schema "seeds" do
    field :name, :string
    field :front_image_id, :string
    field :back_image_id, :string
    field :year, :integer

    timestamps(type: :utc_datetime)
  end

  @doc false
  def changeset(seed, attrs, private_attrs \\ %{}) do
    seed
    |> cast(attrs, [:name, :year])
    |> cast(private_attrs, [:front_image_id, :back_image_id])
    |> validate_number(:year, greater_than_or_equal_to: 2020)
    |> validate_required([:name])
  end

  def new_images(%Ecto.Changeset{} = change) do
    [
      get_change(change, :front_image_id),
      get_change(change, :back_image_id)
    ]
    |> Enum.reject(&is_nil/1)
  end

  def replaced_images(%Ecto.Changeset{} = change) do
    [
      if(changed?(change, :front_image_id), do: change.data.front_image_id),
      if(changed?(change, :back_image_id), do: change.data.back_image_id)
    ]
    |> Enum.reject(&is_nil/1)
  end
end
