package cmd

import (
	"github.com/spf13/cobra"
)

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Manage agent execution prompts",
	Args:  cobra.NoArgs,
}
