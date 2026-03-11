package client

import (
	"fmt"
	"os"
	"os/exec"
)

func OpenSSHTunnel(sshUser, sshHost string, serverPort, clientPort int, sshKeyPath string) (*exec.Cmd, error) {
	args := []string{
		"-N", "-T",
		"-o", "StrictHostKeyChecking=no",
		"-o", "ServerAliveInterval=10",
		"-o", "ServerAliveCountMax=3",
		"-o", "ExitOnForwardFailure=yes",
		"-i", sshKeyPath,
		"-R", fmt.Sprintf("%d:localhost:%d", serverPort, clientPort),
		fmt.Sprintf("%s@%s", sshUser, sshHost),
	}
	cmd := exec.Command("ssh", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start ssh: %w", err)
	}
	return cmd, nil
}
