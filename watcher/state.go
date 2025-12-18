package watcher

import (
	"encoding/json"
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
		return nil, nil
	}

	var s State
	if json.Unmarshal(data, &s) != nil {
		return nil, nil
	}
	return &s, nil
}

func SaveState(s *State) {
	path, err := statePath()
	if err != nil {
		return
	}
	data, _ := json.MarshalIndent(s, "", "  ")
	_ = os.WriteFile(path, data, 0600)
}

