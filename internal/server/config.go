package server

type Config struct {
	Addr         string
	CaddyAdmin   string
	BaseDomain   string
	PortRangeMin int
	PortRangeMax int
}

func DefaultConfig() Config {
	return Config{
		Addr:         ":8080",
		CaddyAdmin:   "http://localhost:2019",
		BaseDomain:   "tunnel.doubletick.dev",
		PortRangeMin: 10000,
		PortRangeMax: 20000,
	}
}
