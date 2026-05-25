// Package cmd provides the CLI commands for the forge CLI tool.
// Output formatting functions are defined in the base sub-package.
// Callers should import base directly: "forge-cli/internal/cmd/base".
package cmd

import (
	"fmt"
	"os"
)

// Debugf prints a debug line to stderr if verbose is true.
// Defined locally for direct use within the cmd package without importing base.
func Debugf(verbose bool, format string, args ...any) {
	if verbose {
		fmt.Fprintf(os.Stderr, "[debug] "+format+"\n", args...)
	}
}
