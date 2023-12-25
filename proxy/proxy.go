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
				Host:   "localhost:4000",
				Path:   "/",
			})
			r.Out.Host = r.In.Host
			r.SetXForwarded()
		},
	}

	if *dev {
		if err := http.ListenAndServe(":8080", p); err != nil {
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
		if err := http.Serve(lns, p); err != nil {
			log.Fatal(err)
		}
	}()
	if err := http.Serve(ln, p); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
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
