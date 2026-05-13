package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var e2eCompileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile-check e2e test files",
	Long: `Run a compile-only check on e2e test files without executing them.
Uses the active profile's compiler (e.g. tsc --noEmit for TypeScript,
go build for Go, python -m compileall for Python).`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Fprintln(os.Stderr, "not yet implemented: forge e2e compile")
		os.Exit(1)
	},
}
