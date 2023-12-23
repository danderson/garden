defmodule Garden.Tribool do
  use Ecto.Type
  def type, do: :boolean

  def cast("y"), do: {:ok, true}
  def cast("n"), do: {:ok, false}
  def cast("?"), do: {:ok, nil}
  def cast(_), do: :error

  def load(x), do: {:ok, x}

  def dump(x), do: {:ok, x}

  def form_value(true), do: "y"
  def form_value("y"), do: "y"
  def form_value(false), do: "n"
  def form_value("n"), do: "n"
  def form_value(nil), do: "?"
  def form_value("?"), do: "?"

  def form_options(), do: [yes: "y", no: "n", unknown: "?"]

  def to_str(true), do: "Yes"
  def to_str(false), do: "No"
  def to_str(nil), do: "Unknown"
end
