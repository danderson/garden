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

alias Garden.Repo
alias Garden.Seeds.Seed
alias Garden.Locations.Location
alias Garden.Plants.Plant

seeds = %{
  cucumber:
    Repo.insert!(%Seed{
      name: "Cucumber",
      year: 2023
    }),
  broccoli:
    Repo.insert!(%Seed{
      name: "Broccoli",
      year: 2023
    }),
  lawn_mix:
    Repo.insert!(%Seed{
      name: "Lawn mix",
      year: 2022
    })
}

locations = %{
  planter_box: Repo.insert!(%Location{name: "Planter box 3"}),
  coffin: Repo.insert!(%Location{name: "Coffin 1"}),
  meadow_left: Repo.insert!(%Location{name: "Meadow left of garage"}),
  gully_right: Repo.insert!(%Location{name: "Gully right of house"})
}

plants = %{
  cucumber:
    Repo.insert!(%Plant{
      name: "Cucumber",
      location_id: locations[:planter_box].id,
      seed_id: seeds[:cucumber].id
    }),
  meadow:
    Repo.insert!(%Plant{
      name: "Meadow lawn",
      location_id: locations[:meadow_left].id,
      seed_id: seeds[:lawn_mix].id
    }),
  broom:
    Repo.insert!(%Plant{
      name: "Fucking broom",
      location_id: locations[:meadow_left].id
    }),
  oak:
    Repo.insert!(%Plant{
      name: "Oak tree",
      location_id: locations[:gully_right].id
    })
}
