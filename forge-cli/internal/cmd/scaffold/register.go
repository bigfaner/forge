package scaffold

import (
	"fmt"
	"sort"

	"forge-cli/pkg/forgeconfig"
	"forge-cli/pkg/project"
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
All recipes include [unix] and [windows] dual-platform variants.

Use --aggregate to generate cross-surface aggregate recipes (install, ci, clean)
by reading the project's surface configuration.`,
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runScaffold,
}

var scaffoldType string
var scaffoldKey string
var scaffoldAggregate bool

func init() {
	Cmd.Flags().StringVar(&scaffoldType, "type", "", "surface type (cli, tui, api, web, mobile)")
	Cmd.Flags().StringVar(&scaffoldKey, "key", "", "surface key (required for named surfaces)")
	Cmd.Flags().BoolVar(&scaffoldAggregate, "aggregate", false, "generate cross-surface aggregate recipes (install, ci, clean)")
}

// Register is a no-op placeholder for consistent sub-package convention.
// The parent cmd package adds Cmd directly.
func Register() {}

func runScaffold(cmd *cobra.Command, _ []string) error {
	if scaffoldAggregate {
		return runAggregate(cmd)
	}

	if scaffoldType == "" {
		return fmt.Errorf("required flag(s) \"type\" not set")
	}

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

// ReadSurfacesFunc is the function used to read surface configuration.
// Exported as a variable for testability — tests can override with a mock.
var ReadSurfacesFunc = defaultReadSurfaces

// defaultReadSurfaces reads surfaces from forgeconfig using the project root.
func defaultReadSurfaces(projectRoot string) (map[string]string, error) {
	return forgeconfig.ReadSurfaces(projectRoot)
}

// runAggregate reads surface configuration and generates aggregate recipes.
func runAggregate(cmd *cobra.Command) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		projectRoot = "."
	}

	surfaces, err := ReadSurfacesFunc(projectRoot)
	if err != nil {
		return fmt.Errorf("failed to read surfaces: %w", err)
	}

	entries := surfacesToEntries(surfaces)

	out, err := GenerateAggregate(entries)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprint(cmd.OutOrStdout(), out)
	return nil
}

// surfacesToEntries converts a map[string]string (from forgeconfig.ReadSurfaces)
// to a sorted slice of SurfaceEntry. Scalar surfaces (key ".") get an empty Key.
func surfacesToEntries(surfaces map[string]string) []SurfaceEntry {
	if len(surfaces) == 0 {
		return nil
	}

	entries := make([]SurfaceEntry, 0, len(surfaces))
	for k, v := range surfaces {
		key := k
		if key == "." {
			key = ""
		}
		entries = append(entries, SurfaceEntry{Key: key, Type: types.SurfaceType(v)})
	}

	// Sort by key for deterministic output
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return entries
}
