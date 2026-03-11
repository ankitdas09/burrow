package client

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg Config) error {
	if cfg.Port == 0 {
		return fmt.Errorf("port is required")
	}
	if cfg.SSHKeyPath == "" {
		return fmt.Errorf("SSH key path is required")
	}
	if _, err := os.Stat(cfg.SSHKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("key file not found: %s", cfg.SSHKeyPath)
	}

	apiURL := fmt.Sprintf("http://%s:%s", cfg.Server, cfg.APIPort)

	if cfg.Named != "" {
		fmt.Printf("→ Registering named tunnel [%s] for localhost:%d ...\n", cfg.Named, cfg.Port)
	} else {
		fmt.Printf("→ Registering tunnel for localhost:%d ...\n", cfg.Port)
	}

	session, err := CreateTunnel(apiURL, cfg.Port, cfg.Named)
	if err != nil {
		return fmt.Errorf("create tunnel: %w", err)
	}

	fmt.Printf("→ Opening SSH tunnel on port %d ...\n", session.ServerPort)
	sshCmd, err := OpenSSHTunnel(cfg.SSHUser, cfg.Server, session.ServerPort, cfg.Port, cfg.SSHKeyPath)
	if err != nil {
		_ = DeleteTunnel(apiURL, session.ID)
		return fmt.Errorf("ssh: %w", err)
	}

	time.Sleep(500 * time.Millisecond)
	if sshCmd.ProcessState != nil && sshCmd.ProcessState.Exited() {
		_ = DeleteTunnel(apiURL, session.ID)
		return fmt.Errorf("SSH tunnel failed to establish")
	}

	PrintBanner(session, cfg.Named)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sshDone := make(chan error, 1)
	go func() { sshDone <- sshCmd.Wait() }()

	select {
	case <-sigCh:
		fmt.Println("\n→ Closing tunnel...")
	case err := <-sshDone:
		if err != nil {
			fmt.Printf("\n→ SSH connection lost: %v\n", err)
		}
	}

	if sshCmd.Process != nil {
		_ = sshCmd.Process.Kill()
	}
	fmt.Printf("→ Deregistering tunnel %s ...\n", session.ID)
	if err := DeleteTunnel(apiURL, session.ID); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to deregister: %v\n", err)
	}
	fmt.Println("→ Done.")
	return nil
}
