package task

import (
	"github.com/spf13/cobra"
)

// Cmd is the parent task command, exported for use by the cmd package.
var Cmd = &cobra.Command{
	Use:   "task",
	Short: "Manage task lifecycle",
	Args:  cobra.NoArgs,
}
