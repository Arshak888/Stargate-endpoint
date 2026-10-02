package main

import (
	"log"
	"time"

	"github.com/Arshak888/Stargate-endpoint/internal/config"
	"github.com/Arshak888/Stargate-endpoint/internal/manager"
	"github.com/Arshak888/Stargate-endpoint/internal/protocol"
	stargateruntime "github.com/Arshak888/Stargate-endpoint/internal/runtime"
	"github.com/Arshak888/Stargate-endpoint/internal/state"
)

const defaultHeartbeatInterval = 30 * time.Second

func main() {
	cfg := config.Load()
	if cfg.ManagerURL == "" {
		log.Fatal("STARGATE_MANAGER_URL is required")
	}

	rt := stargateruntime.NewLocal("/")
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
			Name:            cfg.EndpointName,
			Region:          cfg.Region,
			Country:         cfg.Country,
			City:            cfg.City,
			ProtocolVersion: result.ProtocolVersion,
			Version:         cfg.Version,
			Capabilities:    rt.Capabilities(),
		}
		if err := state.Save(cfg.ConfigPath, s); err != nil {
			log.Fatalf("save endpoint state: %v", err)
		}
		log.Printf("endpoint enrolled: %s", s.EndpointID)
	}

	if s.ManagerURL == "" {
		s.ManagerURL = cfg.ManagerURL
	}
	if s.Version == "" {
		s.Version = cfg.Version
	}
	if len(s.Capabilities) == 0 {
		s.Capabilities = rt.Capabilities()
		if err := state.Save(cfg.ConfigPath, s); err != nil {
			log.Fatalf("save endpoint capabilities: %v", err)
		}
	}

	client := manager.NewClient(s.ManagerURL, s.Credential)
	interval := defaultHeartbeatInterval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	sendHeartbeat := func() {
		snapshot, err := rt.Snapshot()
		if err != nil {
			log.Printf("runtime metrics unavailable: %v", err)
			snapshot = stargateruntime.Snapshot{}
		}

		result, err := client.Heartbeat(protocol.HeartbeatRequest{
			EndpointID:     s.EndpointID,
			Version:        s.Version,
			Capabilities:   s.Capabilities,
			CPUPercent:     snapshot.CPUPercent,
			MemoryPercent:  snapshot.MemoryPercent,
			DiskPercent:    snapshot.DiskPercent,
			ActiveSessions: snapshot.ActiveSessions,
		})
		if err != nil {
			log.Printf("heartbeat failed: %v", err)
			return
		}
		if result.HeartbeatIntervalSecs > 0 {
			next := time.Duration(result.HeartbeatIntervalSecs) * time.Second
			if next != interval {
				interval = next
				ticker.Reset(interval)
			}
		}
		log.Printf("heartbeat acknowledged at %s", result.ObservedAt)
	}

	sendHeartbeat()
	for range ticker.C {
		sendHeartbeat()
	}
}
