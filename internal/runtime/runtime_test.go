package runtime

import "testing"

func TestCapabilitiesAreStableAndCopied(t *testing.T) {
	rt := NewLocal("/")
	first := rt.Capabilities()
	if len(first) != 3 {
		t.Fatalf("expected 3 operational capabilities, got %d", len(first))
	}

	first[0] = "mutated"
	second := rt.Capabilities()
	if second[0] == "mutated" {
		t.Fatal("Capabilities returned internal storage")
	}

	for _, capability := range []string{
		CapabilityHeartbeatTelemetry,
		CapabilityLocalState,
		CapabilityLegacySource,
	} {
		found := false
		for _, got := range second {
			if got == capability {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("missing capability %q", capability)
		}
	}
}
