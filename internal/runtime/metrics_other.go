//go:build !linux

package runtime

import "errors"

func (l *Local) Snapshot() (Snapshot, error) {
	return Snapshot{}, errors.New("local system metrics are currently implemented for Linux endpoints only")
}
