package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type State struct {
	EndpointID       string   `json:"endpoint_id"`
	Credential       string   `json:"credential"`
	ManagerURL       string   `json:"manager_url"`
	Name             string   `json:"name"`
	Region           string   `json:"region"`
	Country          string   `json:"country"`
	City             string   `json:"city"`
	ProtocolVersion  int      `json:"protocol_version"`
	Version          string   `json:"version"`
	Capabilities     []string `json:"capabilities"`
}

func Load(path string) (State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return State{}, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return State{}, err
	}
	if s.EndpointID == "" {
		return State{}, errors.New("endpoint_id missing")
	}
	return s, nil
}

func Save(path string, s State) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
