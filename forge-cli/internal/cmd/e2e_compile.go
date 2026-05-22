package cmd

import (
	"fmt"

	e2e "forge-cli/pkg/e2e"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var e2eCompileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile-check e2e test files",
	Long: `Run a compile-only check on e2e test files without executing them.
	Uses the project's compiler (e.g. tsc --noEmit for TypeScript,
	go build for Go, python -m compileall for Python).`,
	Args: cobra.NoArgs,
	RunE: runE2ECompile,
}

func runE2ECompile(_ *cobra.Command, _ []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return ErrProjectNotFound()
	}

	if err := e2e.Compile(projectRoot); err != nil {
		return fmt.Errorf("e2e compile: %w", err)
	}
	return nil
}
