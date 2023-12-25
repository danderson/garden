defmodule Garden.Seeds.Seed do
  use Ecto.Schema
  import Ecto.Changeset

  schema "seeds" do
    field :name, :string
    field :front_image_id, :string
    field :back_image_id, :string
    field :year, :integer

    field :family, Ecto.Enum,
      values: [
        :Adoxaceae,
        :Allium,
        :Amaranthaceae,
        :Apiaceae,
        :Apocynaceae,
        :Asparagaceae,
        :Asteraceae,
        :Boraginaceae,
        :Brassicaceae,
        :Campanulaceae,
        :Caprifoliaceae,
        :Caryophyllaceae,
        :Convolvulaceae,
        :Cucurbitaceae,
        :Fabaceae,
        :Lamiaceae,
        :Linaceae,
        :Malvaceae,
        :Onagraceae,
        :Papaveraceae,
        :Poaceae,
        :Polygonaceae,
        :Ranunculaceae,
        :Rosaceae,
        :Solanaceae,
        :Tropaeolaceae,
        :Violaceae,
        :Wildflower
      ]

    field :is_native, Garden.Tribool
    field :is_invasive, Garden.Tribool

    field :edible, Garden.Tribool
    field :is_keto, Garden.Tribool

    field :needs_bird_netting, Garden.Tribool
    field :is_deer_resistant, Garden.Tribool

    field :needs_trellis, Garden.Tribool
    field :is_cover_crop, Garden.Tribool
    field :grows_well_from_seed, Garden.Tribool
    field :is_bad_for_cats, Garden.Tribool

    field :type, Ecto.Enum,
      values: [vegetable: "V", fruit: "F", herb: "H", flower: "L", green: "G"]

    field :lifespan, Ecto.Enum, values: [annual: "A", biennial: "B", perennial: "P", unknown: "U"]

    has_many :plants, Garden.Plants.Plant, preload_order: [asc: :name]

    timestamps(type: :utc_datetime)
  end

  def family_mappings() do
    Ecto.Enum.mappings(__MODULE__, :family)
  end

  @doc false
  def upsert_changeset(seed, attrs, private_attrs) do
    seed
    |> cast(attrs, [
      :name,
      :year,
      :family,
      :is_native,
      :is_invasive,
      :edible,
      :is_keto,
      :needs_bird_netting,
      :is_deer_resistant,
      :needs_trellis,
      :is_cover_crop,
      :grows_well_from_seed,
      :is_bad_for_cats
    ])
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
