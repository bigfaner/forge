package cmd

import (
	"fmt"
	"os"

	"forge-cli/pkg/e2eprobe"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var probeCmd = &cobra.Command{
	Use:   "probe [path]",
	Short: "HTTP health check for e2e servers",
	Long: `Probe configured e2e test servers by performing HTTP health checks.

Reads tests/e2e/config.yaml for baseUrl and apiBaseUrl, then probes each
endpoint with an HTTP GET. Exits 0 if all endpoints respond, exit 1 if any
fails. If no config.yaml exists, prints "OK: CLI-only project" and exits 0.

The optional [path] argument specifies the health check path (default: /health).`,
	Args: cobra.MaximumNArgs(1),
	Run:  runProbe,
}

func runProbe(_ *cobra.Command, args []string) {
	path := "/health"
	if len(args) == 1 {
		path = args[0]
	}

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, "OK: CLI-only project")
		os.Exit(0)
	}

	if !e2eprobe.ProbeServers(projectRoot, path) {
		os.Exit(1)
	}
}
