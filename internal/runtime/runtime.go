package runtime

import (
	"runtime"
	"strings"
)

const (
	CapabilityHeartbeatTelemetry = "endpoint.telemetry.heartbeat"
	CapabilityLocalState         = "endpoint.state.local"
	CapabilityLegacySource       = "legacy.stargate-ui.source"
)

var advertisedCapabilities = []string{
	CapabilityHeartbeatTelemetry,
	CapabilityLocalState,
	CapabilityLegacySource,
}

type Snapshot struct {
	CPUPercent     float64
	MemoryPercent  float64
	DiskPercent    float64
	ActiveSessions int
	UptimeSeconds  int64
}

type Runtime interface {
	Capabilities() []string
	Snapshot() (Snapshot, error)
}

type Local struct {
	rootPath string
}

func NewLocal(rootPath string) *Local {
	if strings.TrimSpace(rootPath) == "" {
		rootPath = "/"
	}
	return &Local{rootPath: rootPath}
}

func (l *Local) Capabilities() []string {
	out := make([]string, len(advertisedCapabilities))
	copy(out, advertisedCapabilities)
	return out
}

func Platform() string {
	return runtime.GOOS + "/" + runtime.GOARCH
}
