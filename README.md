# dt-tunnel

Tunnel manager and client: expose a local port via a public HTTPS URL using a central server, Caddy, and SSH reverse forwarding.

## Layout

```
dt-tunnel/
├── cmd/
│   ├── client/          # CLI to create and maintain a tunnel
│   └── server/          # Tunnel manager API + Caddy integration
├── internal/
│   ├── client/          # Client logic (API, SSH, banner)
│   ├── models/          # Shared API types (Session, requests)
│   └── server/          # Server logic (sessions, Caddy, proxy, handlers)
├── go.mod
├── Makefile
└── README.md
```

## Build

```bash
make build
# or
go build -o bin/dt-tunnel-client ./cmd/client
go build -o bin/dt-tunnel-server ./cmd/server
```

## Server (tunnel manager)

Runs the HTTP API that allocates ports and registers routes with Caddy. Caddy must be running with the admin API enabled.

```bash
./bin/dt-tunnel-server
# or with options:
./bin/dt-tunnel-server -addr :8080 -caddy http://localhost:2019 -domain tunnel.example.com -port-min 10000 -port-max 20000
```

Flags:

- `-addr` — listen address (default `:8080`)
- `-caddy` — Caddy admin API URL (default `http://localhost:2019`)
- `-domain` — base domain for tunnels (default `tunnel.doubletick.dev`)
- `-port-min` / `-port-max` — port range for tunnel backends

## Client

Registers a tunnel with the server and opens an SSH reverse tunnel so traffic to the allocated port reaches your local service.

```bash
./bin/dt-tunnel-client -port 3000 -key ./key.pem
# Named tunnel (stable subdomain):
./bin/dt-tunnel-client -port 3000 -key ./key.pem -named user1
```

Flags:

- `-port` — local port to expose (required)
- `-key` — path to SSH private key (required)
- `-named` — optional stable subdomain
- `-server` — tunnel server hostname
- `-api-port` — API port on server (default `8080`)
- `-ssh-user` — SSH user on server (default `ubuntu`)

## API (server)

- `POST /tunnels` — create tunnel (`{"client_port": 3000, "named": "user1"}` optional)
- `GET /tunnels` — list tunnels
- `GET /tunnels/{id}` — get tunnel
- `DELETE /tunnels/{id}` — remove tunnel
- `DELETE /tunnels` — remove all tunnels
