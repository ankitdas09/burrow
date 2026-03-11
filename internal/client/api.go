package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"dt-tunnel/internal/models"
)

func CreateTunnel(apiURL string, port int, named string) (*models.Session, error) {
	req := models.CreateTunnelRequest{ClientPort: port, Named: named}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	resp, err := http.Post(apiURL+"/tunnels", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server %s: %s", resp.Status, string(b))
	}

	var session models.Session
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &session, nil
}

func DeleteTunnel(apiURL, id string) error {
	req, err := http.NewRequest("DELETE", apiURL+"/tunnels/"+id, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
