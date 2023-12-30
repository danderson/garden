package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"flag"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"tailscale.com/tsnet"
)

var (
	hostname = flag.String("hostname", "garden", "hostname")
	stateDir = flag.String("state-dir", "/state/tailscale", "state directory")
	dev      = flag.Bool("dev", false, "development mode")
)

func main() {
	flag.Parse()

	p := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(&url.URL{
				Scheme: "http",
				Host:   "localhost:4000",
				Path:   "/",
			})
			r.Out.Host = r.In.Host
			r.SetXForwarded()
		},
	}

	mux := http.NewServeMux()
	mux.Handle("/", p)
	mux.HandleFunc("/.magic/backup", serveBackup)
	mux.HandleFunc("/.magic/images", serveImages)

	if *dev {
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	s := tsnet.Server{
		Dir:      *stateDir,
		Hostname: *hostname,
	}

	cmd := exec.Command("/app/bin/server")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	go func() {
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	ln, err := s.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
	}
	lns, err := s.ListenTLS("tcp", ":443")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := http.Serve(lns, mux); err != nil {
			log.Fatal(err)
		}
	}()
	if err := http.Serve(ln, mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func serveBackup(w http.ResponseWriter, r *http.Request) {
	db := os.Getenv("DATABASE_PATH")
	if db == "" {
		http.Error(w, "no database known", http.StatusInternalServerError)
		return
	}

	d, err := os.MkdirTemp("", "backup")
	if err != nil {
		http.Error(w, "couldn't make tempdir", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(d)

	backup := filepath.Join(d, "garden_backup.db")
	out, err := exec.Command("sqlite3", db, ".backup "+backup).CombinedOutput()
	if err != nil {
		http.Error(w, string(out), http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, backup)
}

func serveImages(w http.ResponseWriter, r *http.Request) {
	imgs := os.Getenv("IMAGES_PATH")
	if imgs == "" {
		http.Error(w, "no image dir known", http.StatusInternalServerError)
		return
	}

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

		hdr.Name = filepath.Join(rel, hdr.Name)

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

// type debugListener struct {
// 	name string
// 	net.Listener
// }

// func (d *debugListener) Accept() (net.Conn, error) {
// 	log.Printf("XXXXXXXXX %s running accept", d.name)
// 	ret, err := d.Listener.Accept()
// 	log.Printf("XXXXXXXXX %s ACCEPT %v %v", d.name, ret.RemoteAddr(), err)
// 	return &debugConn{d.name, ret}, err
// }

// type debugConn struct {
// 	name string
// 	net.Conn
// }

// func (d *debugConn) Read(b []byte) (n int, err error) {
// 	n, err = d.Conn.Read(b)
// 	log.Printf("RRRRRRRRRRRRRRRRRRR %s %d %v %s", d.name, n, err, string(b[:n]))
// 	return n, err
// }

// func (d *debugConn) Write(b []byte) (n int, err error) {
// 	log.Printf("WWWWWWWWWWWWWWWWWWW %s, %d %s", d.name, len(b), string(b))
// 	return d.Conn.Write(b)
// }
