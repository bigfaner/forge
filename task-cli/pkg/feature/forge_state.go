package feature

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// ForgeState represents the session-level runtime state in .forge/state.json.
type ForgeState struct {
	Feature      string `json:"feature"`
	AllCompleted bool   `json:"allCompleted"`
	UpdatedAt    string `json:"updatedAt"`
}

// WriteForgeState writes .forge/state.json with allCompleted=true.
func WriteForgeState(projectRoot, featureSlug string) error {
	statePath := GetForgeStatePath(projectRoot)
	if err := os.MkdirAll(filepath.Dir(statePath), 0755); err != nil {
		return err
	}

	state := ForgeState{
		Feature:      featureSlug,
		AllCompleted: true,
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return os.WriteFile(statePath, data, 0644)
}

// ReadForgeState reads .forge/state.json. Returns nil if the file doesn't exist.
func ReadForgeState(projectRoot string) *ForgeState {
	statePath := GetForgeStatePath(projectRoot)
	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil
	}
	var state ForgeState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil
	}
	return &state
}

// ClearForgeState removes .forge/state.json.
func ClearForgeState(projectRoot string) error {
	statePath := GetForgeStatePath(projectRoot)
	return os.Remove(statePath)
}
