package scaffold

import (
	"fmt"

	"forge-cli/pkg/types"

	"github.com/spf13/cobra"
)

// Cmd is the scaffold sub-command, exported for parent registration.
var Cmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Generate justfile recipe scaffolding for a surface type",
	Long: `Generate just recipes for a given surface type, outputting
placeholder-templated justfile code to stdout.

Each recipe uses <<PLACEHOLDER>> syntax for project-specific values.
All recipes include [unix] and [windows] dual-platform variants.`,
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runScaffold,
}

var scaffoldType string
var scaffoldKey string

func init() {
	Cmd.Flags().StringVar(&scaffoldType, "type", "", "surface type (cli, tui, api, web, mobile)")
	Cmd.Flags().StringVar(&scaffoldKey, "key", "", "surface key (required for named surfaces)")
	_ = Cmd.MarkFlagRequired("type")
}

// Register is a no-op placeholder for consistent sub-package convention.
// The parent cmd package adds Cmd directly.
func Register() {}

func runScaffold(cmd *cobra.Command, _ []string) error {
	surfaceType := types.SurfaceType(scaffoldType)

	if err := ValidateArgs(surfaceType, scaffoldKey); err != nil {
		return err
	}

	out, err := Generate(surfaceType, scaffoldKey)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprint(cmd.OutOrStdout(), out)
	return nil
}
