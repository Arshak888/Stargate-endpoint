package main

import (
	"log"
	"time"

	"github.com/Arshak888/Stargate-endpoint/internal/config"
	"github.com/Arshak888/Stargate-endpoint/internal/manager"
	"github.com/Arshak888/Stargate-endpoint/internal/protocol"
	"github.com/Arshak888/Stargate-endpoint/internal/state"
)

func main() {
	cfg := config.Load()
	if cfg.ManagerURL == "" {
		log.Fatal("STARGATE_MANAGER_URL is required")
	}

	s, err := state.Load(cfg.ConfigPath)
	if err != nil {
		if cfg.EnrollmentID == "" || cfg.EnrollmentToken == "" {
			log.Fatalf("endpoint is not enrolled and enrollment credentials are missing: %v", err)
		}

		client := manager.NewClient(cfg.ManagerURL, "")
		result, err := client.Bootstrap(protocol.BootstrapRequest{
			EnrollmentID: cfg.EnrollmentID,
			Token:        cfg.EnrollmentToken,
			Name:         cfg.EndpointName,
			Region:       cfg.Region,
			Country:      cfg.Country,
			City:         cfg.City,
			Version:      cfg.Version,
		})
		if err != nil {
			log.Fatalf("bootstrap failed: %v", err)
		}

		s = state.State{
			EndpointID:      result.EndpointID,
			Credential:      result.EndpointCredential,
			ManagerURL:      cfg.ManagerURL,
			Name:           cfg.EndpointName,
			Region:          cfg.Region,
			Country:         cfg.Country,
			City:            cfg.City,
			ProtocolVersion: result.ProtocolVersion,
			Version:         cfg.Version,
			Capabilities:    []string{},
		}
		if err := state.Save(cfg.ConfigPath, s); err != nil {
			log.Fatalf("save endpoint state: %v", err)
		}
		log.Printf("endpoint enrolled: %s", s.EndpointID)
	}

	client := manager.NewClient(s.ManagerURL, s.Credential)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	sendHeartbeat := func() {
		result, err := client.Heartbeat(protocol.HeartbeatRequest{
			EndpointID:     s.EndpointID,
			Version:        s.Version,
			Capabilities:   s.Capabilities,
		})
		if err != nil {
			log.Printf("heartbeat failed: %v", err)
			return
		}
		log.Printf("heartbeat acknowledged at %s", result.ObservedAt)
	}

	sendHeartbeat()
	for range ticker.C {
		sendHeartbeat()
	}
}
