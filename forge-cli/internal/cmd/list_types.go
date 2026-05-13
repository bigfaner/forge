package cmd

import (
	"fmt"

	"forge-cli/pkg/task"

	"github.com/spf13/cobra"
)

var listTypesCmd = &cobra.Command{
	Use:   "list-types",
	Short: "List all supported task types",
	Args:  cobra.NoArgs,
	Run:   runListTypes,
}

func runListTypes(_ *cobra.Command, _ []string) {
	for _, entry := range task.TaskTypeRegistry {
		fmt.Printf("%s  %s\n", entry.Name, entry.Description)
	}
}
