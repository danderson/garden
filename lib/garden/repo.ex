defmodule Garden.Repo do
  use Ecto.Repo,
    otp_app: :garden,
    adapter: Ecto.Adapters.SQLite3
end
