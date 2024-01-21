package main

import (
	"cmp"
	"encoding/csv"
	"errors"
	"flag"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"syscall"

	"go.universe.tf/garden/gogarden/controllers"
	"go.universe.tf/garden/gogarden/db"
	"go.universe.tf/garden/gogarden/htu"
	"go.universe.tf/garden/gogarden/migrations"
	"go.universe.tf/garden/gogarden/views"
	"tailscale.com/tsnet"
	"tailscale.com/types/logger"

	"github.com/a-h/templ"
	"github.com/danderson/reload"
	"github.com/go-chi/chi/v5"

	"embed"
)

//go:embed static/*
var static embed.FS

var (
	hostname = flag.String("hostname", "garden", "hostname")
	stateDir = flag.String("state-dir", "", "state directory")
	dev      = flag.Bool("dev", false, "development mode")
)

func main() {
	flag.Parse()

	dbFile := "garden.db"
	if *dev {
		dbFile = "garden_dev.db"
	}

	if *stateDir != "" {
		if err := os.MkdirAll(*stateDir, 0700); err != nil {
			log.Fatal(err)
		}
	}
	db, err := db.Open(filepath.Join(*stateDir, dbFile), migrations.FileMigrations(), migrations.GoMigrations())
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
	controllers.Home(r, db)
	r.Get("/csv", htu.HandlerFunc(s.serveCSV).ServeHTTP)
	r.Handle("/static/{hash}/*", http.HandlerFunc(s.static))
	r.Handle("/.live", &reload.Reloader{})

	srv := http.Server{
		Handler: r,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sig
		log.Printf("received %v", s)
		srv.Close()
	}()

	if *dev {
		ln, err := net.Listen("tcp", ":8000")
		if err != nil {
			log.Fatalf("listening: %v", err)
		}
		log.Print("Server running on :8000")
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	} else {
		tnsrv := &tsnet.Server{
			Dir:      filepath.Join(*stateDir, "tsnet"),
			Hostname: *hostname,
			Logf:     logger.Discard,
		}
		p80, err := tnsrv.Listen("tcp", ":80")
		if err != nil {
			log.Fatal(err)
		}
		p443, err := tnsrv.ListenTLS("tcp", ":443")
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			if err := srv.Serve(p80); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}()
		log.Print("Server running over tailscale")
		if err := srv.Serve(p443); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}
}

type Server struct {
	db     *db.DB
	assets http.Handler
}

func (s *Server) serveCSV(w http.ResponseWriter, r *http.Request) error {
	seeds, err := s.db.ListSeeds(r.Context())
	if err != nil {
		return err
	}

	slices.SortFunc(seeds, func(a, b db.Seed) int {
		if v := cmp.Compare(a.Family, b.Family); v != 0 {
			return v
		}
		return cmp.Compare(a.Name, b.Name)
	})

	fams := make([]string, 0, len(seeds))
	names := make([]string, 0, len(seeds))
	blank := make([]string, 0, len(seeds))

	for _, s := range seeds {
		fams = append(fams, s.Family.String())
		names = append(names, s.Name)
		blank = append(blank, "")
	}

	w.Header().Set("Content-Type", "text/csv")
	wr := csv.NewWriter(w)
	rc := append([]string{"", "", "", "", "", ""}, fams...)
	if err := wr.Write(rc); err != nil {
		return err
	}

	rc = append([]string{"type", "family", "companions", "do not plant with", "name", ""}, names...)
	if err := wr.Write(rc); err != nil {
		return err
	}

	for i := range fams {
		rc = make([]string, len(blank)+6)
		rc[1] = fams[i]
		rc[4] = names[i]
		if err := wr.Write(rc); err != nil {
			return err
		}
	}

	wr.Flush()
	return nil
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
