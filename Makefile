.PHONY: build build-client build-server clean run-server run-client

BINDIR ?= bin

build: build-client build-server

build-client:
	@mkdir -p $(BINDIR)
	go build -o $(BINDIR)/dt-tunnel-client ./cmd/client

build-server:
	@mkdir -p $(BINDIR)
	go build -o $(BINDIR)/dt-tunnel-server ./cmd/server

clean:
	rm -rf $(BINDIR)

# Run server (tunnel manager API). Requires Caddy running with admin API.
run-server:
	go run ./cmd/server

# Run client. Example: make run-client ARGS="--port 3000 --key ./key.pem"
run-client:
	go run ./cmd/client $(ARGS)
