package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"

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
				Host:   "localhost:8000",
				Path:   "/",
			})
			r.Out.Host = r.In.Host
			r.SetXForwarded()
		},
	}

	if *dev {
		addDevAuth := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set("X-Tailscale-User", "dave@github")
			p.ServeHTTP(w, r)
		})
		if err := http.ListenAndServe(":8080", addDevAuth); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	s := tsnet.Server{
		Dir:      *stateDir,
		Hostname: *hostname,
	}
	lc, err := s.LocalClient()
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("gunicorn", "garden.wsgi")
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

	addAuth := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Del("X-Tailscale-User")
		who, err := lc.WhoIs(r.Context(), r.RemoteAddr)
		if err == nil && who.UserProfile != nil && who.UserProfile.LoginName != "" {
			r.Header.Set("X-Tailscale-User", who.UserProfile.LoginName)
		}
		p.ServeHTTP(w, r)
	})

	if err := http.Serve(ln, addAuth); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
