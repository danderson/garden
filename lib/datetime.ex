defmodule Garden.DateTime do
  use Ecto.Type
  def type, do: :utc_datetime

  def cast(s) when is_binary(s) do
    case DateTime.from_iso8601(s) do
      {:ok, datetime, _} -> {:ok, datetime}
      {:error, _} -> :error
    end
  end

  def cast(%DateTime{} = d), do: d

  def cast(_), do: :error

  def load(%DateTime{} = d) do
    DateTime.shift_zone(d, "America/Vancouver")
  end

  def dump(%DateTime{} = d) do
    DateTime.shift_zone(d, "Etc/UTC")
  end

  def dump(_), do: :error

  def now!(), do: DateTime.now!("America/Vancouver")
end
