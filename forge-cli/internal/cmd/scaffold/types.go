// Package scaffold implements the `forge justfile scaffold` command.
//
// It generates just recipes for a given surface type, outputting
// placeholder-templated justfile code to stdout. Only the recipe
// skeleton is produced — the agent fills project-specific values.
package scaffold

import "forge-cli/pkg/types"

// RecipeKind categorizes recipes into lifecycle and quality groups.
type RecipeKind string

const (
	// KindLifecycle marks recipes that manage surface lifecycle (test, teardown, dev, probe).
	KindLifecycle RecipeKind = "lifecycle"
	// KindQuality marks recipes for code quality checks (compile, fmt, lint, unit-test).
	KindQuality RecipeKind = "quality"
)

// RecipeSpec describes a single recipe to generate for a surface type.
type RecipeSpec struct {
	// Name is the recipe verb (e.g. "test", "teardown", "compile").
	Name string
	// Kind is the recipe category: lifecycle or quality.
	Kind RecipeKind
	// IsScript is true when the body is a multi-line bash script (with shebang).
	// When false, the body is a simple one-liner recipe.
	IsScript bool
	// UnixBody is the [unix] recipe body template with <<PLACEHOLDER>> slots.
	UnixBody string
	// WindowsBody is the [windows] recipe body template with <<PLACEHOLDER>> slots.
	WindowsBody string
}

// SurfaceSpec defines the recipe set for a surface type.
type SurfaceSpec struct {
	// Type is the surface type this spec covers.
	Type types.SurfaceType
	// RequiresKey is true when this surface type must have a --key (named surface).
	// When false, passing --key is an error (scalar surface).
	RequiresKey bool
	// Recipes is the ordered list of recipe specs for this surface.
	Recipes []RecipeSpec
}

// needsKey returns true if the surface type requires a named key (api, web, mobile).
func needsKey(surfaceType types.SurfaceType) bool {
	switch surfaceType {
	case types.SurfaceAPI, types.SurfaceWeb, types.SurfaceMobile:
		return true
	default:
		return false
	}
}

// simpleRecipes returns the shared recipe set for cli and tui surface types.
// These are simple surfaces: no dev/probe lifecycle, only test + teardown + quality.
func simpleRecipes() []RecipeSpec {
	return []RecipeSpec{
		{
			Name:        "test",
			Kind:        KindLifecycle,
			UnixBody:    `cd <<PROJECT_DIR>> && <<TEST_CMD>>`,
			WindowsBody: `cd <<PROJECT_DIR>> && <<TEST_CMD>>`,
		},
		{
			Name:        "teardown",
			Kind:        KindLifecycle,
			UnixBody:    `echo "No teardown needed for CLI/TUI surface"`,
			WindowsBody: `echo No teardown needed for CLI/TUI surface`,
		},
		{
			Name:        "compile",
			Kind:        KindQuality,
			UnixBody:    `cd <<PROJECT_DIR>> && <<COMPILE_CMD>>`,
			WindowsBody: `cd <<PROJECT_DIR>> && <<COMPILE_CMD>>`,
		},
		{
			Name:        "fmt",
			Kind:        KindQuality,
			UnixBody:    `cd <<PROJECT_DIR>> && <<FMT_CMD>>`,
			WindowsBody: `cd <<PROJECT_DIR>> && <<FMT_CMD>>`,
		},
		{
			Name:        "lint",
			Kind:        KindQuality,
			UnixBody:    `cd <<PROJECT_DIR>> && <<LINT_CMD>>`,
			WindowsBody: `cd <<PROJECT_DIR>> && <<LINT_CMD>>`,
		},
		{
			Name:        "unit-test",
			Kind:        KindQuality,
			UnixBody:    `cd <<PROJECT_DIR>> && <<UNIT_TEST_CMD>>`,
			WindowsBody: `cd <<PROJECT_DIR>> && <<UNIT_TEST_CMD>>`,
		},
	}
}

// surfaceSpecs maps each surface type to its spec.
// cli and tui share simpleRecipes; api/web/mobile use serviceRecipes.
var surfaceSpecs = map[types.SurfaceType]SurfaceSpec{
	types.SurfaceCLI: {
		Type:    types.SurfaceCLI,
		Recipes: simpleRecipes(),
	},
	types.SurfaceTUI: {
		Type:    types.SurfaceTUI,
		Recipes: simpleRecipes(),
	},
	types.SurfaceAPI: {
		Type:    types.SurfaceAPI,
		Recipes: serviceRecipes(false),
	},
	types.SurfaceWeb: {
		Type:    types.SurfaceWeb,
		Recipes: serviceRecipes(false),
	},
	types.SurfaceMobile: {
		Type:    types.SurfaceMobile,
		Recipes: serviceRecipes(true),
	},
}

// serviceRecipes returns the recipe set for api/web/mobile surface types.
// When hasTestSetup is true (mobile), an extra test-setup recipe is included
// and the orchestration recipe chains test-setup→dev→probe→test→teardown.
func serviceRecipes(hasTestSetup bool) []RecipeSpec {
	var recipes []RecipeSpec

	if hasTestSetup {
		recipes = append(recipes, RecipeSpec{
			Name:        "test-setup",
			Kind:        KindLifecycle,
			IsScript:    true,
			UnixBody:    serviceTestSetupBody,
			WindowsBody: serviceTestSetupBody,
		})
	}

	recipes = append(recipes,
		RecipeSpec{
			Name:        "dev",
			Kind:        KindLifecycle,
			IsScript:    true,
			UnixBody:    serviceDevUnixBody,
			WindowsBody: serviceDevWindowsBody,
		},
		RecipeSpec{
			Name:        "probe",
			Kind:        KindLifecycle,
			IsScript:    true,
			UnixBody:    serviceProbeBody,
			WindowsBody: serviceProbeBody,
		},
		RecipeSpec{
			Name:        "test",
			Kind:        KindLifecycle,
			IsScript:    true,
			UnixBody:    serviceTestBody,
			WindowsBody: serviceTestBody,
		},
		RecipeSpec{
			Name:        "teardown",
			Kind:        KindLifecycle,
			IsScript:    true,
			UnixBody:    serviceTeardownUnixBody,
			WindowsBody: serviceTeardownWindowsBody,
		},
		RecipeSpec{
			Name:        "compile",
			Kind:        KindQuality,
			UnixBody:    `cd <<PROJECT_DIR>> && <<COMPILE_CMD>>`,
			WindowsBody: `cd <<PROJECT_DIR>> && <<COMPILE_CMD>>`,
		},
		RecipeSpec{
			Name:        "fmt",
			Kind:        KindQuality,
			UnixBody:    `cd <<PROJECT_DIR>> && <<FMT_CMD>>`,
			WindowsBody: `cd <<PROJECT_DIR>> && <<FMT_CMD>>`,
		},
		RecipeSpec{
			Name:        "lint",
			Kind:        KindQuality,
			UnixBody:    `cd <<PROJECT_DIR>> && <<LINT_CMD>>`,
			WindowsBody: `cd <<PROJECT_DIR>> && <<LINT_CMD>>`,
		},
		RecipeSpec{
			Name:        "unit-test",
			Kind:        KindQuality,
			UnixBody:    `cd <<PROJECT_DIR>> && <<UNIT_TEST_CMD>>`,
			WindowsBody: `cd <<PROJECT_DIR>> && <<UNIT_TEST_CMD>>`,
		},
	)

	return recipes
}

// --- Service surface recipe template bodies ---
// These are multi-line bash scripts migrated from server-lifecycle.md.
// They use <<PLACEHOLDER>> syntax for project-specific values.

const serviceDevUnixBody = `#!/usr/bin/env bash
set -euo pipefail
_pid_file=".forge/<<URL_KEY>>.pid"
mkdir -p .forge
# Layer 1: tracked process alive?
if [ -f "$_pid_file" ] && kill -0 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null; then
    echo "<<URL_KEY>>: already running (PID $(tr -d '\r' < "$_pid_file"))"
    exit 0
fi
[ -f "$_pid_file" ] && rm -f "$_pid_file"
# Layer 2: start
<<START_CMD>> &
printf '%s\n' "$!" > "$_pid_file"
_cleanup() { [ -f "$_pid_file" ] && { kill "$(tr -d '\r' < "$_pid_file")" 2>/dev/null || true; rm -f "$_pid_file"; }; }
trap _cleanup EXIT INT TERM
wait`

const serviceDevWindowsBody = `#!/usr/bin/env bash
set -euo pipefail
_pid_file=".forge/<<URL_KEY>>.pid"
mkdir -p .forge
if [ -f "$_pid_file" ] && kill -0 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null; then
    echo "<<URL_KEY>>: already running (PID $(tr -d '\r' < "$_pid_file"))"
    exit 0
fi
[ -f "$_pid_file" ] && rm -f "$_pid_file"
<<START_CMD>> &
printf '%s\n' "$!" > "$_pid_file"
_cleanup() { [ -f "$_pid_file" ] && { kill "$(tr -d '\r' < "$_pid_file")" 2>/dev/null || true; rm -f "$_pid_file"; }; }
trap _cleanup EXIT INT TERM
wait`

const serviceProbeBody = `#!/usr/bin/env bash
set -euo pipefail
_url=<<HEALTH_URL>>
_max_retries=3
_retry_interval=5
_timeout=5
_is_healthy() {
    local status
    status=$(curl -s -o /dev/null -w '%{http_code}' --max-time $_timeout "$_url" 2>/dev/null || echo "000")
    [ "$status" != "000" ] && [ "$status" -lt 500 ]
}
for _i in $(seq 1 $_max_retries); do
    if _is_healthy; then
        echo "OK: <<URL_KEY>> ($_url)"
        exit 0
    fi
    [ "$_i" -lt "$_max_retries" ] && sleep $_retry_interval
done
echo "FAIL: <<URL_KEY>> ($_url) not healthy after ${_max_retries} attempts" >&2
exit 1`

const serviceTestBody = `#!/usr/bin/env bash
set -euo pipefail
cd <<PROJECT_DIR>>
# --- Run Tests ---
<<TEST_CMD>>`

const serviceTeardownUnixBody = `#!/usr/bin/env bash
set -euo pipefail
_pid_file=".forge/<<URL_KEY>>.pid"
if [ -f "$_pid_file" ]; then
    kill "$(tr -d '\r' < "$_pid_file")" 2>/dev/null || true
    rm -f "$_pid_file"
fi`

const serviceTeardownWindowsBody = `#!/usr/bin/env bash
set -euo pipefail
_pid_file=".forge/<<URL_KEY>>.pid"
if [ -f "$_pid_file" ]; then
    _pid="$(tr -d '\r' < "$_pid_file")"
    taskkill //PID "$_pid" //F 2>/dev/null || true
    rm -f "$_pid_file"
fi`

const serviceTestSetupBody = `#!/usr/bin/env bash
set -euo pipefail
# --- Test Setup ---
# TODO: replace with project-specific test setup (e.g., emulator start, mock server)
echo "Running test setup for <<URL_KEY>>..."`
