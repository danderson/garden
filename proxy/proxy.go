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
)

func main() {
	flag.Parse()

	s := tsnet.Server{
		Dir:      *stateDir,
		Hostname: *hostname,
	}
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
	if err := http.Serve(ln, p); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
