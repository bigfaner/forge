// Package cmd provides the CLI commands for the forge CLI tool.
// Output formatting functions are defined in the base sub-package.
// Callers should import base directly: "forge-cli/internal/cmd/base".
//
// Debugf is defined locally because the base package version does not
// correctly expand variadic args (passes args instead of args...).
package cmd

import (
	"fmt"
	"os"
)

// Debugf prints a debug line to stderr if verbose is true.
// Inlined from base to preserve variadic call semantics across package boundaries.
func Debugf(verbose bool, format string, args ...any) {
	if verbose {
		fmt.Fprintf(os.Stderr, "[debug] "+format+"\n", args...)
	}
}
