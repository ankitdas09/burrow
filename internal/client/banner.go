package client

import (
	"fmt"

	"dt-tunnel/internal/models"
)

func PrintBanner(session *models.Session, named string) {
	label := session.Subdomain
	if named != "" {
		label = named
	}
	fmt.Println()
	fmt.Println("  ┌─────────────────────────────────────────────┐")
	fmt.Println("  │  ✓ Tunnel active                            │")
	fmt.Println("  │                                             │")
	fmt.Println("  │  Public URL:                                │")
	fmt.Printf("  │  %-43s│\n", session.PublicURL)
	fmt.Println("  │                                             │")
	fmt.Printf("  │  Named      : %-28s │\n", label)
	fmt.Printf("  │  Local port : %-4d                          │\n", session.ClientPort)
	fmt.Printf("  │  Tunnel ID  : %-28s │\n", session.ID)
	fmt.Println("  └─────────────────────────────────────────────┘")
	fmt.Println()
	fmt.Println("  Press Ctrl+C to close the tunnel.")
	fmt.Println()
}
