package cmd

import (
	"sort"
	"strings"

	"forge-cli/pkg/forgeconfig"

	"github.com/spf13/cobra"
)

// KnownSurfaceTypes lists all valid surface type values.
// Unknown types are filtered from --types output and silently ignored in listings.
var KnownSurfaceTypes = map[string]bool{
	"web":    true,
	"api":    true,
	"cli":    true,
	"tui":    true,
	"mobile": true,
}

var surfacesTypesFlag bool

var surfacesCmd = &cobra.Command{
	Use:   "surfaces [path]",
	Short: "Query project surfaces configuration",
	Long: `Query the surfaces field from .forge/config.yaml.

Without arguments: scalar form outputs the single type; map form outputs
one "path=surface" line per entry. Exit 0 always.

With a path argument: returns the surface type for that path using
segment prefix matching. Exit 0 on match, exit 1 with stderr error on
no match.

With --types: outputs space-separated deduplicated list of known surface types.`,
	Args:          cobra.MaximumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runSurfaces,
}

func init() {
	surfacesCmd.Flags().BoolVar(&surfacesTypesFlag, "types", false, "list deduplicated surface types")
	surfacesCmd.Flags().String("project-root", "", "project root directory (defaults to auto-detection)")
}

// runSurfaces implements the three sub-invocations of the surfaces command.
func runSurfaces(cmd *cobra.Command, args []string) error {
	projectRoot := resolveProjectRoot(cmd)

	surfaces, err := forgeconfig.ReadSurfaces(projectRoot)
	if err != nil {
		return err
	}

	// --types: output space-separated deduplicated known types
	if surfacesTypesFlag {
		return runSurfacesTypes(cmd, surfaces)
	}

	// No args: list all surfaces
	if len(args) == 0 {
		return runSurfacesList(cmd, surfaces)
	}

	// Path query: match surface for given path
	return runSurfacesQuery(cmd, surfaces, args[0])
}

// runSurfacesList outputs all surfaces.
// Scalar form: single type string.
// Map form: one "path=surface" line per entry.
func runSurfacesList(cmd *cobra.Command, surfaces map[string]string) error {
	if len(surfaces) == 0 {
		return nil
	}

	out := cmd.OutOrStdout()

	// Scalar form: single "." key
	if len(surfaces) == 1 {
		if v, ok := surfaces["."]; ok {
			write(out, "%s\n", v)
			return nil
		}
	}

	// Map form: output sorted path=surface lines
	paths := make([]string, 0, len(surfaces))
	for p := range surfaces {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	for _, p := range paths {
		write(out, "%s=%s\n", p, surfaces[p])
	}
	return nil
}

// surfacesPathError is a sentinel error type for unmatched path queries.
// The actual error message is written directly to stderr before returning
// this error, ensuring raw unformatted output on stderr (no "Error: " prefix).
type surfacesPathError struct {
	message string
}

func (e *surfacesPathError) Error() string { return e.message }

// runSurfacesQuery finds the surface type for a query path.
// Exit 0 on match (stdout), exit 1 on no match (stderr).
// Error message is written raw to stderr -- no formatting prefix.
func runSurfacesQuery(cmd *cobra.Command, surfaces map[string]string, query string) error {
	result, err := forgeconfig.MatchSurface(surfaces, query)
	if err != nil {
		// Write raw error to stderr (no "Error: " prefix per Hard Rules)
		write(cmd.ErrOrStderr(), "%s\n", err.Error())
		return &surfacesPathError{message: err.Error()}
	}

	write(cmd.OutOrStdout(), "%s\n", result)
	return nil
}

// runSurfacesTypes outputs space-separated deduplicated known surface types.
func runSurfacesTypes(cmd *cobra.Command, surfaces map[string]string) error {
	if len(surfaces) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	var types []string
	for _, typ := range surfaces {
		if seen[typ] {
			continue
		}
		if !KnownSurfaceTypes[typ] {
			continue
		}
		seen[typ] = true
		types = append(types, typ)
	}

	if len(types) == 0 {
		return nil
	}

	sort.Strings(types)
	write(cmd.OutOrStdout(), "%s\n", strings.Join(types, " "))
	return nil
}
