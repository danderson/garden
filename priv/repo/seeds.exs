# Script for populating the database. You can run it as:
#
#     mix run priv/repo/seeds.exs
#
# Inside the script, you can read and write to any of your
# repositories directly:
#
#     Garden.Repo.insert!(%Garden.SomeSchema{})
#
# We recommend using the bang functions (`insert!`, `update!`
# and so on) as they will fail if something goes wrong.

alias Garden.{Repo, Seeds, Locations, Plants}
alias Garden.Plants.PlantLocation

if Mix.env() == :dev do
  image = fn name ->
    Seeds.store_seed_image(Path.expand("./seed_images/#{name}.jpg", __DIR__))
  end

  seeds = %{
    lawn_mix:
      Seeds.new!(
        %{
          name: "Lawn mix",
          year: 2023
        },
        %{
          front_image_id: image.("lawn_mix_front"),
          back_image_id: image.("lawn_mix_back")
        }
      ),
    wildflowers:
      Seeds.new!(
        %{
          name: "Wildflowers",
          year: 2022
        },
        %{
          front_image_id: image.("wildflower_front"),
          back_image_id: image.("wildflower_back")
        }
      ),
    broccoli:
      Seeds.new!(%{
        name: "Broccoli",
        year: 2023
      }),
    cucumber:
      Seeds.new!(%{
        name: "Cucumber",
        year: 2022
      })
  }

  locations = %{
    freezer: Locations.new!(%{name: "Freezer"}),
    planter_box: Locations.new!(%{name: "Planter box 3"}),
    coffin: Locations.new!(%{name: "Coffin 1"}),
    meadow_left: Locations.new!(%{name: "Meadow left of garage"}),
    gully_right: Locations.new!(%{name: "Gully right of house"})
  }

  now = Garden.DateTime.now!()

  plants = %{
    cucumber:
      Plants.new!(%{
        plant: %{
          name: "Cucumber"
        },
        location_id: locations[:planter_box].id
      }),
    meadow:
      Plants.new!(%{
        plant: %{
          name: "Meadow lawn",
          seed_id: seeds[:lawn_mix].id
        },
        location_id: locations[:meadow_left].id
      }),
    broom:
      Plants.new!(%{
        plant: %{
          name: "Fucking broom"
        },
        location_id: locations[:meadow_left].id
      }),
    oak:
      Plants.new!(%{
        plant: %{
          name: "Oak tree"
        },
        location_id: locations[:gully_right].id
      })
  }

  last = DateTime.add(now, -10, :day)

  %PlantLocation{
    plant_id: plants[:cucumber].id,
    location_id: locations[:coffin].id,
    start: last,
    end: now
  }
  |> Repo.insert!()

  last2 = DateTime.add(last, -15, :day)

  %PlantLocation{
    plant_id: plants[:cucumber].id,
    location_id: locations[:freezer].id,
    start: last2,
    end: last
  }
  |> Repo.insert!()
end
