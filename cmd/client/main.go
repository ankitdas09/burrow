package main

import (
	"flag"
	"fmt"
	"os"

	"dt-tunnel/internal/client"
)

func main() {
	port := flag.Int("port", 0, "Local port to tunnel (required)")
	named := flag.String("named", "", "Named tunnel for stable URL (e.g. --named user1)")
	key := flag.String("key", "", "Path to SSH private key (required)")
	server := flag.String("server", client.DefaultServer, "Tunnel server hostname")
	apiPort := flag.String("api-port", client.DefaultAPIPort, "Tunnel manager API port")
	sshUser := flag.String("ssh-user", client.DefaultSSHUser, "SSH user on tunnel server")
	flag.Parse()

	cfg := client.Config{
		Port:       *port,
		Named:      *named,
		SSHKeyPath: *key,
		Server:     *server,
		APIPort:    *apiPort,
		SSHUser:    *sshUser,
	}

	if err := client.Run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		fmt.Fprintln(os.Stderr, "Usage: dt-tunnel-client --port 3000 --key ./key.pem [--named user1]")
		os.Exit(1)
	}
}
