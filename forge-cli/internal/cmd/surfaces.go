package cmd

import (
	"encoding/json"
	"sort"
	"strings"

	"forge-cli/pkg/forgeconfig"
	"forge-cli/pkg/types"

	"github.com/spf13/cobra"
)

// KnownSurfaceTypes re-exports the canonical set from forgeconfig.
// Unknown types are filtered from --types output and silently ignored in listings.
var KnownSurfaceTypes = forgeconfig.KnownSurfaceTypes

var surfacesTypesFlag bool
var surfacesJSONFlag bool

var surfacesCmd = &cobra.Command{
	Use:   "surfaces [path]",
	Short: "Query project surfaces configuration",
	Long: `Query the surfaces field from .forge/config.yaml.

Without arguments: scalar form outputs the single type; map form outputs
one "path=surface" line per entry. Exit 0 always.

With a path argument: returns the surface type for that path using
segment prefix matching. Exit 0 on match, exit 1 with stderr error on
no match.

With --types: outputs space-separated deduplicated list of known surface types.

With --json: outputs structured JSON. List mode: {"surfaces": [...]};
query mode: [{"key": ..., "type": ...}]; types mode: {"types": [...]}.`,
	Args:          cobra.MaximumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runSurfaces,
}

func init() {
	surfacesCmd.Flags().BoolVar(&surfacesTypesFlag, "types", false, "list deduplicated surface types")
	surfacesCmd.Flags().BoolVar(&surfacesJSONFlag, "json", false, "output in JSON format")
	surfacesCmd.Flags().String("project-root", "", "project root directory (defaults to auto-detection)")
}

// jsonError writes a structured JSON error to stderr when --json is active.
// Hard Rule: must use json.NewEncoder(cmd.ErrOrStderr()).Encode(), not fmt.Fprintf.
func jsonError(cmd *cobra.Command, message string) error {
	return json.NewEncoder(cmd.ErrOrStderr()).Encode(map[string]string{"error": message})
}

// runSurfaces implements the three sub-invocations of the surfaces command.
func runSurfaces(cmd *cobra.Command, args []string) error {
	projectRoot := resolveProjectRoot(cmd)

	surfaces, err := forgeconfig.ReadSurfaces(projectRoot)
	if err != nil {
		if surfacesJSONFlag {
			_ = jsonError(cmd, err.Error())
			return err
		}
		return err
	}

	// --types: output deduplicated known types
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
		if surfacesJSONFlag {
			// No surfaces configured -- output JSON error to stderr
			_ = jsonError(cmd, "no surface configured; run `forge init` to configure surfaces")
			return &surfacesPathError{message: "no surface configured"}
		}
		return nil
	}

	// JSON mode: output {"surfaces": [...]}
	if surfacesJSONFlag {
		return runSurfacesListJSON(cmd, surfaces)
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

// runSurfacesListJSON outputs all surfaces as JSON.
func runSurfacesListJSON(cmd *cobra.Command, surfaces map[string]string) error {
	type jsonEntry struct {
		Key  string `json:"key"`
		Type string `json:"type"`
	}

	entries := make([]jsonEntry, 0, len(surfaces))
	for k, v := range surfaces {
		entries = append(entries, jsonEntry{Key: k, Type: v})
	}
	// Sort by key for deterministic output
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return json.NewEncoder(cmd.OutOrStdout()).Encode(map[string][]jsonEntry{
		"surfaces": entries,
	})
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
	// Handle empty surfaces before calling MatchSurface for JSON mode distinction
	if len(surfaces) == 0 {
		if surfacesJSONFlag {
			_ = jsonError(cmd, "no surface configured; run `forge init` to configure surfaces")
			return &surfacesPathError{message: "no surface configured"}
		}
		write(cmd.ErrOrStderr(), "no surface configured; run `forge init` to configure surfaces\n")
		return &surfacesPathError{message: "no surface configured"}
	}

	match, err := forgeconfig.MatchSurface(surfaces, query)
	if err != nil {
		if surfacesJSONFlag {
			// Path-specific no-match: output empty array, exit 0
			return json.NewEncoder(cmd.OutOrStdout()).Encode([]forgeconfig.SurfaceMatch{})
		}
		// Write raw error to stderr (no "Error: " prefix per Hard Rules)
		write(cmd.ErrOrStderr(), "%s\n", err.Error())
		return &surfacesPathError{message: err.Error()}
	}

	if surfacesJSONFlag {
		return json.NewEncoder(cmd.OutOrStdout()).Encode([]forgeconfig.SurfaceMatch{match})
	}

	write(cmd.OutOrStdout(), "%s\n", match.Type)
	return nil
}

// runSurfacesTypes outputs deduplicated known surface types.
func runSurfacesTypes(cmd *cobra.Command, surfaces map[string]string) error {
	if len(surfaces) == 0 {
		if surfacesJSONFlag {
			_ = jsonError(cmd, "no surface configured; run `forge init` to configure surfaces")
			return &surfacesPathError{message: "no surface configured"}
		}
		return nil
	}

	seen := make(map[string]bool)
	var surfaceTypes []string
	for _, typ := range surfaces {
		if seen[typ] {
			continue
		}
		if !KnownSurfaceTypes[types.SurfaceType(typ)] {
			continue
		}
		seen[typ] = true
		surfaceTypes = append(surfaceTypes, typ)
	}

	if len(surfaceTypes) == 0 {
		if surfacesJSONFlag {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(map[string][]string{"types": {}})
		}
		return nil
	}

	sort.Strings(surfaceTypes)

	if surfacesJSONFlag {
		return json.NewEncoder(cmd.OutOrStdout()).Encode(map[string][]string{"types": surfaceTypes})
	}

	write(cmd.OutOrStdout(), "%s\n", strings.Join(surfaceTypes, " "))
	return nil
}
