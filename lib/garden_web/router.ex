defmodule GardenWeb.Router do
  use GardenWeb, :router

  pipeline :browser do
    plug(:accepts, ["html"])
    plug(:fetch_session)
    plug(:fetch_live_flash)
    plug(:put_root_layout, html: {GardenWeb.Layouts, :root})
    plug(:protect_from_forgery)
    plug(:put_secure_browser_headers)
  end

  pipeline :api do
    plug(:accepts, ["json"])
  end

  scope "/", GardenWeb do
    pipe_through(:browser)

    get("/", PageController, :home)

    live("/plants", PlantLive.Index, :index)
    live("/plants/new", PlantLive.Index, :new)
    live("/plants/:id/edit", PlantLive.Index, :edit)

    live("/plants/:id", PlantLive.Show, :show)
    live("/plants/:id/show/edit", PlantLive.Show, :edit)
  end

  # Other scopes may use custom stacks.
  # scope "/api", GardenWeb do
  #   pipe_through :api
  # end

  # Enable LiveDashboard in development
  if Application.compile_env(:garden, :dev_routes) do
    # If you want to use the LiveDashboard in production, you should put
    # it behind authentication and allow only admins to access it.
    # If your application does not have an admins-only section yet,
    # you can use Plug.BasicAuth to set up some basic authentication
    # as long as you are also using SSL (which you should anyway).
    import Phoenix.LiveDashboard.Router

    scope "/dev" do
      pipe_through(:browser)

      live_dashboard("/dashboard", metrics: GardenWeb.Telemetry)
    end
  end
end
