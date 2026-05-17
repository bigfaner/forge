// Package e2e provides shared e2e execution logic, replacing justfile bash.
package e2e

import (
	"errors"
	"fmt"

	"forge-cli/pkg/profile"
)

// Sentinel errors for profile resolution.
var (
	ErrNoProfile       = errors.New("no e2e profile configured")
	ErrBadProfile      = errors.New("unknown profile")
	ErrFeatureNotFound = errors.New("feature not found")
)

// RunOpts holds options for e2e operations.
type RunOpts struct {
	ProjectRoot string
	Feature     string // empty = run all
	Force       bool   // for setup
}

// ResolveProfile reads config.yaml and validates the profile.
// Returns the profile name or an error.
func ResolveProfile(projectRoot string) (string, error) {
	profiles, err := profile.ReadLanguages(projectRoot)
	if err != nil {
		return "", fmt.Errorf("read languages: %w", err)
	}

	if len(profiles) == 0 {
		return "", ErrNoProfile
	}

	name := profiles[0]
	if !profile.IsKnownLanguage(name) {
		return "", fmt.Errorf("%w: %s", ErrBadProfile, name)
	}

	return name, nil
}
