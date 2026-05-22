package cmd

import (
	"github.com/spf13/cobra"
)

var e2eCmd = &cobra.Command{
	Use:   "e2e",
	Short: "End-to-end test management",
	Args:  cobra.NoArgs,
}
