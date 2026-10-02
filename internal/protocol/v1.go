package protocol

const Version = 1

type BootstrapRequest struct {
	EnrollmentID string `json:"enrollment_id"`
	Token        string `json:"token"`
	Name         string `json:"name,omitempty"`
	Region       string `json:"region,omitempty"`
	Country      string `json:"country,omitempty"`
	City         string `json:"city,omitempty"`
	Version      string `json:"version,omitempty"`
}

type BootstrapResponse struct {
	EndpointID              string `json:"endpoint_id"`
	EndpointCredential      string `json:"endpoint_credential"`
	ProtocolVersion         int    `json:"protocol_version"`
	HeartbeatIntervalSecs   int    `json:"heartbeat_interval_seconds"`
}

type HeartbeatRequest struct {
	EndpointID    string   `json:"endpoint_id"`
	Version       string   `json:"version"`
	Capabilities  []string `json:"capabilities"`
	CPUPercent    float64  `json:"cpu_percent"`
	MemoryPercent float64  `json:"memory_percent"`
	DiskPercent   float64  `json:"disk_percent"`
	ActiveSessions int     `json:"active_sessions"`
}

type HeartbeatResponse struct {
	Status                    string `json:"status"`
	ObservedAt                string `json:"observed_at"`
	HeartbeatIntervalSecs     int    `json:"heartbeat_interval_seconds"`
}
