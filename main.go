package main

import (
	"archive/tar"
	"cmp"
	"compress/gzip"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"slices"
	"syscall"

	"go.universe.tf/garden/controllers"
	"go.universe.tf/garden/db"
	"go.universe.tf/garden/htu"
	"go.universe.tf/garden/migrations"
	"go.universe.tf/garden/views"
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
	dbPath := filepath.Join(*stateDir, dbFile)
	db, err := db.Open(dbPath, migrations.FileMigrations(), migrations.GoMigrations())
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

	imagesPath := filepath.Join(*stateDir, "images")

	s := &Server{
		dbPath:     dbPath,
		imagesPath: imagesPath,
		db:         db,
		assets:     http.FileServer(http.FS(assets)),
	}

	r := chi.NewRouter()
	controllers.Seeds(r, db)
	controllers.Locations(r, db)
	controllers.Plants(r, db)
	controllers.Home(r, db)
	r.Get("/csv", htu.HandlerFunc(s.serveCSV).ServeHTTP)
	r.Handle("/static/{hash}/*", http.HandlerFunc(s.static))
	r.Handle("/user_images/*", http.StripPrefix("/user_images", http.FileServer(http.Dir(imagesPath))))
	if *dev {
		r.Handle("/.live", &reload.Reloader{})
	} else {
		r.HandleFunc("/.live", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/javascript")
			w.WriteHeader(http.StatusOK)
		})
	}
	r.Handle("/.magic/db", htu.HandlerFunc(s.serveDB))
	r.HandleFunc("/.magic/images", s.serveImages)

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
	dbPath     string
	imagesPath string
	db         *db.DB
	assets     http.Handler
}

func (s *Server) serveCSV(w http.ResponseWriter, r *http.Request) error {
	seeds, err := s.db.SearchSeeds(r.Context(), "%")
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

func (s *Server) serveDB(w http.ResponseWriter, r *http.Request) error {
	d, err := os.MkdirTemp("", "backup")
	if err != nil {
		return fmt.Errorf("couldn't make tempdir: %w", err)
	}
	defer os.RemoveAll(d)

	backup := filepath.Join(d, "garden_backup.db")
	out, err := exec.Command("sqlite3", s.dbPath, ".backup "+backup).CombinedOutput()
	if err != nil {
		return fmt.Errorf("couldn't make backup: %w (%s)", err, string(out))
	}

	http.ServeFile(w, r, backup)
	return nil
}

func (s *Server) serveImages(w http.ResponseWriter, r *http.Request) {
	imgs := s.imagesPath

	g := gzip.NewWriter(w)
	defer g.Close()

	t := tar.NewWriter(g)
	defer t.Close()
	base := filepath.Dir(imgs)

	w.Header().Set("Content-Type", "application/tar+gzip")

	err := filepath.WalkDir(imgs, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(base, path)
		if err != nil {
			return err
		}

		fi, err := d.Info()
		if err != nil {
			return err
		}

		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return err
		}

		hdr.Name = rel

		if err := t.WriteHeader(hdr); err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := io.Copy(t, f); err != nil {
			return err
		}
		f.Close()
		return nil
	})
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) render(w http.ResponseWriter, r *http.Request, page templ.Component) {
	views.Root(page).Render(r.Context(), w)
}
