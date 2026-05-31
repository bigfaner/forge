package feature

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"forge-cli/pkg/index"
)

// ForgeState represents the session-level runtime state in .forge/state.json.
type ForgeState struct {
	Feature      string `json:"feature"`
	AllCompleted bool   `json:"allCompleted"`
	CompletedAt  string `json:"completedAt,omitempty"` // set by feature complete hook on success
	UpdatedAt    string `json:"updatedAt"`
}

// WriteForgeState writes .forge/state.json with allCompleted=true.
func WriteForgeState(projectRoot, featureSlug string) error {
	statePath := GetForgeStatePath(projectRoot)
	if err := os.MkdirAll(filepath.Dir(statePath), 0o755); err != nil {
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
	return index.AtomicWrite(statePath, data, 0o644)
}

// MarkFeatureCompleted sets completedAt on the existing state.json.
// Does not touch other fields. No-op if state.json doesn't exist.
func MarkFeatureCompleted(projectRoot string) error {
	statePath := GetForgeStatePath(projectRoot)
	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil // no state file — nothing to mark
	}

	var state ForgeState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil // malformed — don't touch
	}

	state.CompletedAt = time.Now().Format(time.RFC3339)
	state.UpdatedAt = time.Now().Format(time.RFC3339)

	updated, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return index.AtomicWrite(statePath, updated, 0o644)
}

// EnsureForgeDir creates the .forge/ directory at project root if it doesn't exist.
// This bootstraps the forge workspace marker so that subagents running from
// subdirectories can find the project root via FindProjectRoot().
func EnsureForgeDir(projectRoot string) error {
	forgeDir := filepath.Join(projectRoot, ForgeDir)
	return os.MkdirAll(forgeDir, 0o755)
}

// EnsureForgeState writes .forge/state.json with allCompleted=false.
// Called by task claim to create the workspace marker and active session state.
// Overwrites any existing file (e.g., after fix-e2e tasks reset completion).
func EnsureForgeState(projectRoot, featureSlug string) error {
	statePath := GetForgeStatePath(projectRoot)
	if err := os.MkdirAll(filepath.Dir(statePath), 0o755); err != nil {
		return err
	}
	state := ForgeState{
		Feature:      featureSlug,
		AllCompleted: false,
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return index.AtomicWrite(statePath, data, 0o644)
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

// ClearForgeState writes .forge/state.json with allCompleted=false.
// The file is preserved (not deleted) so the workspace marker remains.
func ClearForgeState(projectRoot string) error {
	statePath := GetForgeStatePath(projectRoot)
	state := ReadForgeState(projectRoot)
	if state == nil {
		return nil // no state file — nothing to clear
	}
	state.AllCompleted = false
	state.UpdatedAt = time.Now().Format(time.RFC3339)
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return index.AtomicWrite(statePath, data, 0o644)
}
