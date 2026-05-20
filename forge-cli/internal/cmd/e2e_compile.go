package cmd

import (
	"fmt"
	"os"

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
	Run: runE2ECompile,
}

func runE2ECompile(_ *cobra.Command, _ []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	if err := e2e.Compile(projectRoot); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
