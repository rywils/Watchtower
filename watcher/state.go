package watcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"watchtower/internal/util"
)

type Device struct {
	IP       string `json:"ip"`
	MAC      string `json:"mac"`
	LastSeen int64  `json:"last_seen"`
}

type State struct {
	Devices map[string]Device `json:"devices"`
}

func NewState() *State {
	return &State{Devices: map[string]Device{}}
}

func statePath() (string, error) {
	dir, err := util.StateDir("watchtower")
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "state.json"), nil
}

func LoadState() (*State, error) {
	path, err := statePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read state file %q: %w", path, err)
	}

	if len(data) == 0 {
		return nil, nil
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse state file %q: %w", path, err)
	}

	if s.Devices == nil {
		s.Devices = map[string]Device{}
	}
	return &s, nil
}

func SaveState(s *State) error {
	path, err := statePath()
	if err != nil {
		return err
	}
	data, _ := json.MarshalIndent(s, "", "  ")
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write state file %q: %w", path, err)
	}
	return nil
}
