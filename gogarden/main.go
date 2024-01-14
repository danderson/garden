package main

import (
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"

	"go.universe.tf/garden/gogarden/controllers"
	"go.universe.tf/garden/gogarden/db"
	"go.universe.tf/garden/gogarden/migrations"
	"go.universe.tf/garden/gogarden/views"

	"github.com/a-h/templ"
	"github.com/danderson/reload"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"

	"embed"
)

//go:embed static/*
var static embed.FS

var (
	dev = flag.Bool("dev", false, "development mode")
)

func main() {
	flag.Parse()

	logger := slog.Default()
	db, err := db.Open(logger, "garden_dev.db", migrations.FileMigrations(), migrations.GoMigrations())
	if err != nil {
		log.Fatalf("opening database: %v", err)
	}
	defer db.Close()

	assets, err := fs.Sub(static, "static")
	if err != nil {
		log.Fatalf("Getting subdir of embedded assets: %v", err)
	}
	if *dev {
		log.Printf("Running in dev mode")
		assets = os.DirFS("static")
	}

	s := &Server{
		db:     db,
		assets: http.FileServer(http.FS(assets)),
	}

	r := chi.NewRouter()
	controllers.Seeds(r, db)
	controllers.Locations(r, db)
	controllers.Plants(r, db)
	// r.Get("/api/locations", htu.ErrHandler(s.listLocations))
	// r.Post("/api/locations/{id}", htu.ErrHandler(s.updateLocation))
	r.Handle("/static/{hash}/*", http.HandlerFunc(s.static))
	r.Handle("/.live", &reload.Reloader{})
	log.Print("Server running")
	http.ListenAndServe(":8000", r)
}

type Server struct {
	db     *db.DB
	assets http.Handler
}

func (s *Server) static(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "*")
	r.URL.Path = "/" + slug
	r.URL.RawPath = "/" + slug
	s.assets.ServeHTTP(w, r)
}

func (s *Server) render(w http.ResponseWriter, r *http.Request, page templ.Component) {
	views.Root(page).Render(r.Context(), w)
}
