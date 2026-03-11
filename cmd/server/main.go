package main

import (
	"flag"
	"log"
	"net/http"

	"dt-tunnel/internal/server"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address for tunnel manager API")
	caddyAdmin := flag.String("caddy", "http://localhost:2019", "Caddy admin API base URL")
	baseDomain := flag.String("domain", "tunnel.doubletick.dev", "Base domain for tunnels")
	portMin := flag.Int("port-min", 10000, "Minimum port for tunnel backends")
	portMax := flag.Int("port-max", 20000, "Maximum port for tunnel backends")
	flag.Parse()

	cfg := server.Config{
		Addr:         *addr,
		CaddyAdmin:   *caddyAdmin,
		BaseDomain:   *baseDomain,
		PortRangeMin: *portMin,
		PortRangeMax: *portMax,
	}

	sm := server.NewSessionManager(cfg)
	mux := server.NewMux(sm)

	log.Printf("Tunnel manager listening on %s", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, mux))
}
