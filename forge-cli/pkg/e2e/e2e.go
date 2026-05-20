// Package e2e provides shared e2e execution logic, replacing justfile bash.
package e2e

import (
	"errors"
)

// Sentinel errors for e2e operations.
var (
	ErrFeatureNotFound = errors.New("feature not found")
)

// RunOpts holds options for e2e operations.
type RunOpts struct {
	ProjectRoot string
	Feature     string // empty = run all
	Force       bool   // for setup
}
