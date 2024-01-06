package main

import (
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"

	"go.universe.tf/garden/gogarden/db"
	"go.universe.tf/garden/gogarden/migrations"
	"go.universe.tf/garden/gogarden/views"

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
	r.Get("/seeds", s.listSeeds)
	// r.Get("/api/locations", htu.ErrHandler(s.listLocations))
	// r.Post("/api/locations/{id}", htu.ErrHandler(s.updateLocation))
	r.Handle("/static/{hash}/*", http.HandlerFunc(s.static))
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

func (s *Server) listSeeds(w http.ResponseWriter, r *http.Request) {
	views.Seeds().Render(r.Context(), w)
}

// func (s *Server) listLocations(w http.ResponseWriter, r *http.Request) error {
// 	ctx := r.Context()
// 	tx, err := s.db.ReadTx(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	locs, err := tx.GetLocations(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	htu.RespondJSON(w, locs)
// 	return nil
// }

// func (s *Server) updateLocation(w http.ResponseWriter, r *http.Request) error {
// 	ctx := r.Context()
// 	id, err := htu.Int64Param(r, "id")
// 	if err != nil {
// 		return htu.BadRequest("invalid id", err)
// 	}
// 	var args struct {
// 		Name    string        `json:"name"`
// 		QRState types.QRState `json:"qr_state"`
// 	}
// 	if err := htu.JSONBody(r, &args); err != nil {
// 		return htu.BadRequest("failed to parse body", err)
// 	}

// 	tx, err := s.db.ReadTx(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	err = tx.UpdateLocation(ctx, db.UpdateLocationParams{
// 		ID:      id,
// 		Name:    &args.Name,
// 		QRState: args.QRState,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	if err := tx.Commit(); err != nil {
// 		return err
// 	}
// 	return nil
// }
