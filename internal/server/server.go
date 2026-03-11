package server

import "net/http"

func NewMux(sm *SessionManager) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /tunnels", sm.CreateTunnel)
	mux.HandleFunc("GET /tunnels", sm.ListTunnels)
	mux.HandleFunc("GET /tunnels/{id}", sm.GetTunnel)
	mux.HandleFunc("DELETE /tunnels", sm.DeleteAllTunnels)
	mux.HandleFunc("DELETE /tunnels/{id}", sm.DeleteTunnel)
	return mux
}
