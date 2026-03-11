package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (sm *SessionManager) registerCaddyRoute(subdomain string, serverPort int) error {
	route := map[string]interface{}{
		"@id":   subdomain,
		"match": []map[string]interface{}{{"host": []string{fmt.Sprintf("%s.%s", subdomain, sm.cfg.BaseDomain)}}},
		"handle": []map[string]interface{}{{
			"handler": "subroute",
			"routes": []map[string]interface{}{{
				"handle": []map[string]interface{}{{
					"handler":   "reverse_proxy",
					"upstreams": []map[string]interface{}{{"dial": fmt.Sprintf("localhost:%d", serverPort)}},
				}},
			}},
		}},
	}
	body, err := json.Marshal(route)
	if err != nil {
		return fmt.Errorf("marshal route: %w", err)
	}
	resp, err := http.Post(
		fmt.Sprintf("%s/config/apps/http/servers/srv0/routes", sm.cfg.CaddyAdmin),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("caddy API request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("caddy returned %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func (sm *SessionManager) deregisterCaddyRoute(subdomain string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/id/%s", sm.cfg.CaddyAdmin, subdomain), nil)
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
