package client

type Config struct {
	Port       int
	Named      string
	SSHKeyPath string
	Server     string
	APIPort    string
	SSHUser    string
}

const (
	DefaultServer   = "your-server.com"
	DefaultAPIPort  = "8080"
	DefaultSSHUser  = "ubuntu"
)
