package main

import (
	"log"
	"net/http"

	"go.universe.tf/garden/gogarden/db"
	"go.universe.tf/garden/gogarden/htu"
	"go.universe.tf/garden/gogarden/migrations"
	"go.universe.tf/garden/gogarden/types"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
)

func main() {
	logger := slog.Default()
	db, err := db.Open(logger, "garden_dev.db", migrations.FileMigrations(), migrations.GoMigrations())
	if err != nil {
		log.Fatalf("opening database: %v", err)
	}
	defer db.Close()

	s := &Server{db: db}

	r := chi.NewRouter()
	r.Get("/api/locations", htu.ErrHandler(s.listLocations))
	r.Post("/api/locations/{id}", htu.ErrHandler(s.updateLocation))
	http.ListenAndServe(":8000", r)
}

type Server struct {
	db *db.DB
}

func (s *Server) listLocations(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tx, err := s.db.ReadTx(ctx)
	if err != nil {
		return err
	}
	locs, err := tx.GetLocations(ctx)
	if err != nil {
		return err
	}
	htu.RespondJSON(w, locs)
	return nil
}

func (s *Server) updateLocation(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	id, err := htu.Int64Param(r, "id")
	if err != nil {
		return htu.BadRequest("invalid id", err)
	}
	var args struct {
		Name    string        `json:"name"`
		QRState types.QRState `json:"qr_state"`
	}
	if err := htu.JSONBody(r, &args); err != nil {
		return htu.BadRequest("failed to parse body", err)
	}

	tx, err := s.db.ReadTx(ctx)
	if err != nil {
		return err
	}
	err = tx.UpdateLocation(ctx, db.UpdateLocationParams{
		ID:      id,
		Name:    args.Name,
		QRState: args.QRState,
	})
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
