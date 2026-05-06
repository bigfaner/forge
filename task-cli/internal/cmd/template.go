package cmd

import (
	"fmt"
	"strings"

	"task-cli/pkg/template"

	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template [name]",
	Short: "Show template content by name, or list all templates",
	Long: `View task templates managed by task-cli.

Without arguments: list all available template names.
With a name argument: print the template content.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runTemplate,
}

func init() {
	rootCmd.AddCommand(templateCmd)
}

func runTemplate(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		listTemplates()
		return
	}
	showTemplate(args[0])
}

func listTemplates() {
	names := template.List()
	if len(names) == 0 {
		fmt.Println("No templates available.")
		return
	}
	fmt.Println("Available templates:")
	for _, n := range names {
		fmt.Printf("  %s\n", n)
	}
	fmt.Printf("\nUsage: task template <name> | task add --template <name> ...\n")
}

func showTemplate(name string) {
	content, err := template.Get(name)
	if err != nil {
		Exit(err)
	}
	fmt.Print(content)
	if !strings.HasSuffix(content, "\n") {
		fmt.Println()
	}
}
