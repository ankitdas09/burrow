package models

import "time"

type Session struct {
	ID         string    `json:"id"`
	Subdomain  string    `json:"subdomain"`
	ServerPort int       `json:"server_port"`
	ClientPort int       `json:"client_port"`
	PublicURL  string    `json:"public_url"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateTunnelRequest struct {
	ClientPort int    `json:"client_port"`
	Named      string `json:"named,omitempty"`
}
