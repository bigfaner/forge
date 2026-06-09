package scaffold

import (
	"strings"
	"testing"

	"forge-cli/pkg/types"
)

// --- AC-1: recipe completeness and placeholder syntax ---

func TestGenerate_CLI_ContainsAllRecipes(t *testing.T) {
	out, err := Generate(types.SurfaceCLI, "")
	if err != nil {
		t.Fatalf("Generate(cli): %v", err)
	}

	required := []string{"test", "teardown", "compile", "fmt", "lint", "unit-test"}
	for _, name := range required {
		// recipe name must appear as a line start (e.g. "test [unix]:")
		if !strings.Contains(out, name+" [unix]:") {
			t.Errorf("missing recipe %q in output", name)
		}
	}
}

func TestGenerate_CLI_PlaceholderSyntax(t *testing.T) {
	out, err := Generate(types.SurfaceCLI, "")
	if err != nil {
		t.Fatalf("Generate(cli): %v", err)
	}

	// Must contain <<...>> placeholders
	if !strings.Contains(out, "<<") || !strings.Contains(out, ">>") {
		t.Error("output missing <<PLACEHOLDER>> syntax")
	}

	// Must NOT contain {{...}} template syntax
	if strings.Contains(out, "{{") || strings.Contains(out, "}}") {
		t.Error("output contains {{...}} syntax, must use <<...>> instead")
	}
}

// --- AC-2: argument validation ---

func TestValidateArgs_UnknownType(t *testing.T) {
	err := ValidateArgs("unknown", "")
	if err == nil {
		t.Fatal("expected error for unknown surface type")
	}
	if !strings.Contains(err.Error(), "unknown surface type") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateArgs_ScalarWithKey(t *testing.T) {
	err := ValidateArgs(types.SurfaceCLI, "myapp")
	if err == nil {
		t.Fatal("expected error when scalar surface gets --key")
	}
	if !strings.Contains(err.Error(), "scalar") || !strings.Contains(err.Error(), "--key") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateArgs_ScalarTUIWithKey(t *testing.T) {
	err := ValidateArgs(types.SurfaceTUI, "myapp")
	if err == nil {
		t.Fatal("expected error when tui scalar surface gets --key")
	}
}

func TestValidateArgs_ScalarWithoutKey(t *testing.T) {
	if err := ValidateArgs(types.SurfaceCLI, ""); err != nil {
		t.Errorf("cli without key should be valid: %v", err)
	}
	if err := ValidateArgs(types.SurfaceTUI, ""); err != nil {
		t.Errorf("tui without key should be valid: %v", err)
	}
}

// --- AC-3: cli and tui have identical recipe sets, no dev/probe ---

func TestGenerate_TUI_SameRecipeSetAsCLI(t *testing.T) {
	cliOut, err := Generate(types.SurfaceCLI, "")
	if err != nil {
		t.Fatalf("Generate(cli): %v", err)
	}
	tuiOut, err := Generate(types.SurfaceTUI, "")
	if err != nil {
		t.Fatalf("Generate(tui): %v", err)
	}

	required := []string{"test", "teardown", "compile", "fmt", "lint", "unit-test"}
	for _, name := range required {
		if !strings.Contains(tuiOut, name+" [unix]:") {
			t.Errorf("tui missing recipe %q", name)
		}
	}

	forbidden := []string{"dev", "probe"}
	for _, name := range forbidden {
		if strings.Contains(cliOut, name+" [unix]:") {
			t.Errorf("cli should not have recipe %q", name)
		}
		if strings.Contains(tuiOut, name+" [unix]:") {
			t.Errorf("tui should not have recipe %q", name)
		}
	}
}

// --- AC-4: user-customized marker + scalar no prefix ---

func TestGenerate_AllRecipesMarkedUserCustomized(t *testing.T) {
	out, err := Generate(types.SurfaceCLI, "")
	if err != nil {
		t.Fatalf("Generate(cli): %v", err)
	}

	recipes := []string{"test", "teardown", "compile", "fmt", "lint", "unit-test"}
	for _, name := range recipes {
		// Each recipe should be preceded by "# user-customized"
		marker := "# user-customized\n" + name + " [unix]:"
		if !strings.Contains(out, marker) {
			t.Errorf("recipe %q missing '# user-customized' marker", name)
		}
	}
}

func TestGenerate_ScalarNoPrefix(t *testing.T) {
	out, err := Generate(types.SurfaceCLI, "")
	if err != nil {
		t.Fatalf("Generate(cli): %v", err)
	}

	// Scalar surface should have bare recipe names like "test [unix]:"
	// NOT prefixed names like "cli-test [unix]:" or "myapp-test [unix]:"
	// Check for key-prefixed forms that start with a non-trivial key prefix.
	// "unit-test" is a valid verb, not a key prefix — so we check that
	// the known verbs appear WITHOUT a preceding key segment.
	forbiddenPrefixes := []string{
		"cli-test ", "cli-teardown ", "cli-compile ", "cli-fmt ", "cli-lint ", "cli-unit-test ",
	}
	for _, fp := range forbiddenPrefixes {
		if strings.Contains(out, fp) {
			t.Errorf("scalar surface should not have key-prefixed recipe: %q found", fp)
		}
	}
}

func TestGenerate_NamedWithPrefix(t *testing.T) {
	// This test validates the prefix logic even though api is not yet a full spec.
	// We'll test with a mock: manually call recipeName.
	name := recipeName("backend", "test")
	if name != "backend-test" {
		t.Errorf("recipeName(backend, test) = %q, want backend-test", name)
	}
	name = recipeName("", "test")
	if name != "test" {
		t.Errorf("recipeName('', test) = %q, want test", name)
	}
}

// --- AC-5: dual platform variants ---

func TestGenerate_DualPlatformVariants(t *testing.T) {
	out, err := Generate(types.SurfaceCLI, "")
	if err != nil {
		t.Fatalf("Generate(cli): %v", err)
	}

	recipes := []string{"test", "teardown", "compile", "fmt", "lint", "unit-test"}
	for _, name := range recipes {
		if !strings.Contains(out, name+" [unix]:") {
			t.Errorf("recipe %q missing [unix] variant", name)
		}
		if !strings.Contains(out, name+" [windows]:") {
			t.Errorf("recipe %q missing [windows] variant", name)
		}
	}
}

// --- Table-driven: recipeName ---

func TestRecipeName(t *testing.T) {
	tests := []struct {
		key  string
		verb string
		want string
	}{
		{"", "test", "test"},
		{"", "compile", "compile"},
		{"backend", "test", "backend-test"},
		{"backend", "compile", "backend-compile"},
		{"frontend", "lint", "frontend-lint"},
	}
	for _, tt := range tests {
		got := recipeName(tt.key, tt.verb)
		if got != tt.want {
			t.Errorf("recipeName(%q, %q) = %q, want %q", tt.key, tt.verb, got, tt.want)
		}
	}
}

// --- Coverage: Generate error path and needsKey branches ---

func TestGenerate_UnknownType(t *testing.T) {
	_, err := Generate("unknown", "")
	if err == nil {
		t.Fatal("expected error for unknown surface type in Generate")
	}
	if !strings.Contains(err.Error(), "unknown surface type") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNeedsKey_NamedTypes(t *testing.T) {
	named := []types.SurfaceType{types.SurfaceAPI, types.SurfaceWeb, types.SurfaceMobile}
	for _, typ := range named {
		if !needsKey(typ) {
			t.Errorf("needsKey(%q) = false, want true", typ)
		}
	}
}

func TestNeedsKey_ScalarTypes(t *testing.T) {
	scalars := []types.SurfaceType{types.SurfaceCLI, types.SurfaceTUI}
	for _, typ := range scalars {
		if needsKey(typ) {
			t.Errorf("needsKey(%q) = true, want false", typ)
		}
	}
}

func TestValidateArgs_NamedWithoutKey(t *testing.T) {
	// api requires --key but none provided
	err := ValidateArgs(types.SurfaceAPI, "")
	if err == nil {
		t.Fatal("expected error when named surface missing --key")
	}
	if !strings.Contains(err.Error(), "requires --key") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateArgs_NamedWithKey_Supported(t *testing.T) {
	// api with key is now valid and supported for generation
	if err := ValidateArgs(types.SurfaceAPI, "backend"); err != nil {
		t.Errorf("api with key should be valid: %v", err)
	}
	if err := ValidateArgs(types.SurfaceWeb, "frontend"); err != nil {
		t.Errorf("web with key should be valid: %v", err)
	}
	if err := ValidateArgs(types.SurfaceMobile, "app"); err != nil {
		t.Errorf("mobile with key should be valid: %v", err)
	}
}

func TestGenerate_TUI_PlaceholderSyntax(t *testing.T) {
	out, err := Generate(types.SurfaceTUI, "")
	if err != nil {
		t.Fatalf("Generate(tui): %v", err)
	}
	if strings.Contains(out, "{{") {
		t.Error("tui output contains {{...}} syntax")
	}
	if !strings.Contains(out, "<<") {
		t.Error("tui output missing <<PLACEHOLDER>> syntax")
	}
}

func TestGenerate_TUI_DualPlatformVariants(t *testing.T) {
	out, err := Generate(types.SurfaceTUI, "")
	if err != nil {
		t.Fatalf("Generate(tui): %v", err)
	}
	for _, name := range []string{"test", "teardown", "compile", "fmt", "lint", "unit-test"} {
		if !strings.Contains(out, name+" [unix]:") {
			t.Errorf("tui recipe %q missing [unix] variant", name)
		}
		if !strings.Contains(out, name+" [windows]:") {
			t.Errorf("tui recipe %q missing [windows] variant", name)
		}
	}
}

// --- Table-driven: ValidateArgs ---

func TestValidateArgs_Table(t *testing.T) {
	tests := []struct {
		name    string
		typ     types.SurfaceType
		key     string
		wantErr bool
		errMsg  string
	}{
		{"cli scalar ok", types.SurfaceCLI, "", false, ""},
		{"tui scalar ok", types.SurfaceTUI, "", false, ""},
		{"cli with key error", types.SurfaceCLI, "app", true, "scalar"},
		{"tui with key error", types.SurfaceTUI, "app", true, "scalar"},
		{"unknown type error", "unknown", "", true, "unknown surface type"},
		{"api with key ok", types.SurfaceAPI, "backend", false, ""},
		{"web with key ok", types.SurfaceWeb, "frontend", false, ""},
		{"mobile with key ok", types.SurfaceMobile, "app", false, ""},
		{"api without key error", types.SurfaceAPI, "", true, "requires --key"},
		{"web without key error", types.SurfaceWeb, "", true, "requires --key"},
		{"mobile without key error", types.SurfaceMobile, "", true, "requires --key"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateArgs(tt.typ, tt.key)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("error %q should contain %q", err.Error(), tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// ============================================================================
// Task 2: api/web/mobile service surface templates
// ============================================================================

// --- AC-1: api surface recipe completeness with backend- prefix ---

func TestGenerate_API_RecipeCompleteness(t *testing.T) {
	out, err := Generate(types.SurfaceAPI, "backend")
	if err != nil {
		t.Fatalf("Generate(api, backend): %v", err)
	}

	// Lifecycle recipes
	lifecycle := []string{"backend-dev", "backend-probe", "backend-test", "backend-teardown", "backend"}
	for _, name := range lifecycle {
		if !strings.Contains(out, name+" [unix]:") {
			t.Errorf("api missing lifecycle recipe %q", name)
		}
		if !strings.Contains(out, name+" [windows]:") {
			t.Errorf("api missing [windows] variant for %q", name)
		}
	}

	// Quality recipes
	quality := []string{"backend-compile", "backend-fmt", "backend-lint", "backend-unit-test"}
	for _, name := range quality {
		if !strings.Contains(out, name+" [unix]:") {
			t.Errorf("api missing quality recipe %q", name)
		}
		if !strings.Contains(out, name+" [windows]:") {
			t.Errorf("api missing [windows] variant for %q", name)
		}
	}
}

func TestGenerate_API_OrchestrationSequence(t *testing.T) {
	out, err := Generate(types.SurfaceAPI, "backend")
	if err != nil {
		t.Fatalf("Generate(api, backend): %v", err)
	}

	// The orchestration recipe "backend" should chain: dev→probe→test→teardown
	// Pattern: just backend-dev && just backend-probe && just backend-test; rc=$?; just backend-teardown; exit $rc
	if !strings.Contains(out, "just backend-dev") {
		t.Error("api orchestration missing 'just backend-dev'")
	}
	if !strings.Contains(out, "just backend-probe") {
		t.Error("api orchestration missing 'just backend-probe'")
	}
	if !strings.Contains(out, "just backend-test") {
		t.Error("api orchestration missing 'just backend-test'")
	}
	if !strings.Contains(out, "just backend-teardown") {
		t.Error("api orchestration missing 'just backend-teardown'")
	}
}

func TestGenerate_API_NoForbiddenRecipes(t *testing.T) {
	out, err := Generate(types.SurfaceAPI, "backend")
	if err != nil {
		t.Fatalf("Generate(api, backend): %v", err)
	}

	// api should NOT have test-setup
	if strings.Contains(out, "backend-test-setup") {
		t.Error("api should not have test-setup recipe")
	}
}

// --- AC-2: web same as api; mobile has extra test-setup ---

func TestGenerate_WEB_SameRecipesAsAPI(t *testing.T) {
	webOut, err := Generate(types.SurfaceWeb, "frontend")
	if err != nil {
		t.Fatalf("Generate(web, frontend): %v", err)
	}

	// Same lifecycle + quality recipes as api
	required := []string{
		"frontend-dev", "frontend-probe", "frontend-test",
		"frontend-teardown", "frontend",
		"frontend-compile", "frontend-fmt", "frontend-lint", "frontend-unit-test",
	}
	for _, name := range required {
		if !strings.Contains(webOut, name+" [unix]:") {
			t.Errorf("web missing recipe %q", name)
		}
	}

	// web orchestration should also chain dev→probe→test→teardown
	if !strings.Contains(webOut, "just frontend-dev") {
		t.Error("web orchestration missing dev step")
	}
	if !strings.Contains(webOut, "just frontend-probe") {
		t.Error("web orchestration missing probe step")
	}

	// No test-setup
	if strings.Contains(webOut, "frontend-test-setup") {
		t.Error("web should not have test-setup recipe")
	}
}

func TestGenerate_MOBILE_HasTestSetup(t *testing.T) {
	out, err := Generate(types.SurfaceMobile, "app")
	if err != nil {
		t.Fatalf("Generate(mobile, app): %v", err)
	}

	// mobile has all api/web recipes PLUS test-setup
	required := []string{
		"app-test-setup", "app-dev", "app-probe", "app-test",
		"app-teardown", "app",
		"app-compile", "app-fmt", "app-lint", "app-unit-test",
	}
	for _, name := range required {
		if !strings.Contains(out, name+" [unix]:") {
			t.Errorf("mobile missing recipe %q", name)
		}
		if !strings.Contains(out, name+" [windows]:") {
			t.Errorf("mobile missing [windows] variant for %q", name)
		}
	}
}

func TestGenerate_MOBILE_OrchestrationIncludesTestSetup(t *testing.T) {
	out, err := Generate(types.SurfaceMobile, "app")
	if err != nil {
		t.Fatalf("Generate(mobile, app): %v", err)
	}

	// mobile orchestration: test-setup→dev→probe→test→teardown
	if !strings.Contains(out, "just app-test-setup") {
		t.Error("mobile orchestration missing test-setup step")
	}
	if !strings.Contains(out, "just app-dev") {
		t.Error("mobile orchestration missing dev step")
	}
}

// --- AC-3: PID file management and health check retry ---

func TestGenerate_API_PIDFileManagement(t *testing.T) {
	out, err := Generate(types.SurfaceAPI, "backend")
	if err != nil {
		t.Fatalf("Generate(api, backend): %v", err)
	}

	// dev recipe must contain PID file reference
	if !strings.Contains(out, "<<URL_KEY>>.pid") {
		t.Error("dev recipe missing PID file path with <<URL_KEY>> placeholder")
	}

	// teardown recipe must clean up PID file
	// Look for "rm -f" in context of pid file
	if !strings.Contains(out, "rm -f") {
		t.Error("teardown recipe missing PID file cleanup (rm -f)")
	}
}

func TestGenerate_API_HealthCheckRetry(t *testing.T) {
	out, err := Generate(types.SurfaceAPI, "backend")
	if err != nil {
		t.Fatalf("Generate(api, backend): %v", err)
	}

	// probe recipe must contain health check URL placeholder
	if !strings.Contains(out, "<<HEALTH_URL>>") {
		t.Error("probe recipe missing <<HEALTH_URL>> placeholder")
	}

	// Must contain retry loop pattern
	if !strings.Contains(out, "max_retries") && !strings.Contains(out, "_max_retries") && !strings.Contains(out, "retry") {
		t.Error("probe recipe missing health check retry logic")
	}
}

func TestGenerate_API_StartCommandPlaceholder(t *testing.T) {
	out, err := Generate(types.SurfaceAPI, "backend")
	if err != nil {
		t.Fatalf("Generate(api, backend): %v", err)
	}

	// dev recipe must have <<START_CMD>> placeholder
	if !strings.Contains(out, "<<START_CMD>>") {
		t.Error("dev recipe missing <<START_CMD>> placeholder")
	}
}

func TestGenerate_API_IdempotentStart(t *testing.T) {
	out, err := Generate(types.SurfaceAPI, "backend")
	if err != nil {
		t.Fatalf("Generate(api, backend): %v", err)
	}

	// dev recipe must check if process is already running (idempotent start)
	if !strings.Contains(out, "already running") {
		t.Error("dev recipe missing idempotent start check ('already running')")
	}
}

func TestGenerate_API_WindowsTeardownUsesTaskkill(t *testing.T) {
	out, err := Generate(types.SurfaceAPI, "backend")
	if err != nil {
		t.Fatalf("Generate(api, backend): %v", err)
	}

	// windows teardown variant must use taskkill instead of kill
	if !strings.Contains(out, "taskkill") {
		t.Error("windows teardown missing taskkill command")
	}
}

// --- AC-4: all recipes marked user-customized with dual-platform ---

func TestGenerate_API_AllRecipesMarkedUserCustomized(t *testing.T) {
	out, err := Generate(types.SurfaceAPI, "backend")
	if err != nil {
		t.Fatalf("Generate(api, backend): %v", err)
	}

	recipes := []string{
		"backend-dev", "backend-probe", "backend-test", "backend-teardown",
		"backend", // orchestration recipe
		"backend-compile", "backend-fmt", "backend-lint", "backend-unit-test",
	}
	for _, name := range recipes {
		marker := "# user-customized\n" + name + " [unix]:"
		if !strings.Contains(out, marker) {
			t.Errorf("recipe %q missing '# user-customized' marker before [unix] variant", name)
		}
	}
}

func TestGenerate_WEB_AllRecipesMarkedUserCustomized(t *testing.T) {
	out, err := Generate(types.SurfaceWeb, "frontend")
	if err != nil {
		t.Fatalf("Generate(web, frontend): %v", err)
	}

	recipes := []string{
		"frontend-dev", "frontend-probe", "frontend-test", "frontend-teardown",
		"frontend",
		"frontend-compile", "frontend-fmt", "frontend-lint", "frontend-unit-test",
	}
	for _, name := range recipes {
		marker := "# user-customized\n" + name + " [unix]:"
		if !strings.Contains(out, marker) {
			t.Errorf("recipe %q missing '# user-customized' marker", name)
		}
	}
}

func TestGenerate_MOBILE_AllRecipesMarkedUserCustomized(t *testing.T) {
	out, err := Generate(types.SurfaceMobile, "app")
	if err != nil {
		t.Fatalf("Generate(mobile, app): %v", err)
	}

	recipes := []string{
		"app-test-setup", "app-dev", "app-probe", "app-test", "app-teardown",
		"app",
		"app-compile", "app-fmt", "app-lint", "app-unit-test",
	}
	for _, name := range recipes {
		marker := "# user-customized\n" + name + " [unix]:"
		if !strings.Contains(out, marker) {
			t.Errorf("recipe %q missing '# user-customized' marker", name)
		}
	}
}

func TestGenerate_API_PlaceholderSyntax(t *testing.T) {
	out, err := Generate(types.SurfaceAPI, "backend")
	if err != nil {
		t.Fatalf("Generate(api, backend): %v", err)
	}

	if strings.Contains(out, "{{") || strings.Contains(out, "}}") {
		t.Error("api output contains {{...}} syntax, must use <<...>>")
	}
}

func TestGenerate_Mobile_PlaceholderSyntax(t *testing.T) {
	out, err := Generate(types.SurfaceMobile, "app")
	if err != nil {
		t.Fatalf("Generate(mobile, app): %v", err)
	}

	if strings.Contains(out, "{{") || strings.Contains(out, "}}") {
		t.Error("mobile output contains {{...}} syntax, must use <<...>>")
	}
}

// ============================================================================
// Task 3: scaffold --aggregate mode
// ============================================================================

// --- AC-1: aggregate generates install, ci, clean without user-customized marker ---

func TestGenerateAggregate_ThreeAggregates(t *testing.T) {
	surfaces := []SurfaceEntry{
		{Key: "backend", Type: types.SurfaceAPI},
		{Key: "frontend", Type: types.SurfaceWeb},
	}
	out, err := GenerateAggregate(surfaces)
	if err != nil {
		t.Fatalf("GenerateAggregate: %v", err)
	}

	for _, name := range []string{"install", "ci", "clean"} {
		if !strings.Contains(out, name+" [unix]:") {
			t.Errorf("missing aggregate recipe %q [unix]", name)
		}
		if !strings.Contains(out, name+" [windows]:") {
			t.Errorf("missing aggregate recipe %q [windows]", name)
		}
	}
}

func TestGenerateAggregate_NoUserCustomizedMarker(t *testing.T) {
	surfaces := []SurfaceEntry{
		{Key: "backend", Type: types.SurfaceAPI},
		{Key: "frontend", Type: types.SurfaceWeb},
	}
	out, err := GenerateAggregate(surfaces)
	if err != nil {
		t.Fatalf("GenerateAggregate: %v", err)
	}

	// Aggregate recipes must NOT have # user-customized marker
	for _, name := range []string{"install", "ci", "clean"} {
		marker := "# user-customized\n" + name + " [unix]:"
		if strings.Contains(out, marker) {
			t.Errorf("aggregate recipe %q should NOT have '# user-customized' marker", name)
		}
	}
}

// --- AC-2: ci recipe aggregates lint + compile + unit-test, no surface-level test ---

func TestGenerateAggregate_CIRecipe_ContainsLintCompileUnitTest(t *testing.T) {
	surfaces := []SurfaceEntry{
		{Key: "backend", Type: types.SurfaceAPI},
		{Key: "frontend", Type: types.SurfaceWeb},
	}
	out, err := GenerateAggregate(surfaces)
	if err != nil {
		t.Fatalf("GenerateAggregate: %v", err)
	}

	// ci should contain all lint, compile, unit-test calls for each surface
	for _, key := range []string{"backend", "frontend"} {
		for _, verb := range []string{"lint", "compile", "unit-test"} {
			recipeRef := "just " + key + "-" + verb
			if !strings.Contains(out, recipeRef) {
				t.Errorf("ci recipe missing %q", recipeRef)
			}
		}
	}
}

func TestGenerateAggregate_CIRecipe_ExcludesSurfaceTest(t *testing.T) {
	surfaces := []SurfaceEntry{
		{Key: "backend", Type: types.SurfaceAPI},
		{Key: "frontend", Type: types.SurfaceWeb},
	}
	out, err := GenerateAggregate(surfaces)
	if err != nil {
		t.Fatalf("GenerateAggregate: %v", err)
	}

	// ci recipe must NOT reference surface-level test (e.g. "just backend-test")
	// We need to check specifically within the ci recipe section.
	// Since the output is sequential, we extract the ci recipe block.
	ciStart := strings.Index(out, "ci [unix]:")
	if ciStart == -1 {
		t.Fatal("ci recipe not found")
	}
	// Extract the ci block: from "ci [unix]:" up to the next top-level recipe
	var ciBlock string
	if idx := strings.Index(out[ciStart+1:], "\nclean"); idx != -1 {
		ciBlock = out[ciStart : ciStart+1+idx]
	} else {
		ciBlock = out[ciStart:]
	}

	for _, key := range []string{"backend", "frontend"} {
		testRef := "just " + key + "-test"
		if strings.Contains(ciBlock, testRef) {
			t.Errorf("ci recipe should NOT contain surface-level test %q", testRef)
		}
	}
}

// --- AC-3: multiple service-type surfaces generate test-setup aggregate ---

func TestGenerateAggregate_MultiServiceGeneratesTestSetup(t *testing.T) {
	surfaces := []SurfaceEntry{
		{Key: "backend", Type: types.SurfaceAPI},
		{Key: "frontend", Type: types.SurfaceWeb},
	}
	out, err := GenerateAggregate(surfaces)
	if err != nil {
		t.Fatalf("GenerateAggregate: %v", err)
	}

	if !strings.Contains(out, "test-setup [unix]:") {
		t.Error("multi-service aggregate should contain test-setup recipe")
	}
	if !strings.Contains(out, "test-setup [windows]:") {
		t.Error("multi-service aggregate test-setup missing [windows] variant")
	}
}

func TestGenerateAggregate_TestSetup_DependencyOrder(t *testing.T) {
	surfaces := []SurfaceEntry{
		{Key: "frontend", Type: types.SurfaceWeb},
		{Key: "backend", Type: types.SurfaceAPI},
	}
	out, err := GenerateAggregate(surfaces)
	if err != nil {
		t.Fatalf("GenerateAggregate: %v", err)
	}

	// Startup order: api → web → mobile (regardless of input order)
	// teardown is reverse: mobile → web → api
	// So "just backend-dev" must appear before "just frontend-dev"
	backendIdx := strings.Index(out, "just backend-dev")
	frontendIdx := strings.Index(out, "just frontend-dev")
	if backendIdx == -1 || frontendIdx == -1 {
		t.Fatal("test-setup missing service dev references")
	}
	if backendIdx >= frontendIdx {
		t.Error("test-setup startup order wrong: api (backend) should come before web (frontend)")
	}

	// Teardown: reverse order — frontend teardown before backend teardown
	frontendTeardownIdx := strings.LastIndex(out, "just frontend-teardown")
	backendTeardownIdx := strings.LastIndex(out, "just backend-teardown")
	if frontendTeardownIdx >= backendTeardownIdx {
		t.Error("test-setup teardown order wrong: web (frontend) teardown should come before api (backend)")
	}
}

func TestGenerateAggregate_TestSetup_WithMobileAPIMultiService(t *testing.T) {
	surfaces := []SurfaceEntry{
		{Key: "backend", Type: types.SurfaceAPI},
		{Key: "app", Type: types.SurfaceMobile},
	}
	out, err := GenerateAggregate(surfaces)
	if err != nil {
		t.Fatalf("GenerateAggregate: %v", err)
	}

	// startup order: api → mobile
	backendIdx := strings.Index(out, "just backend-dev")
	appIdx := strings.Index(out, "just app-dev")
	if backendIdx == -1 || appIdx == -1 {
		t.Fatal("test-setup missing service references")
	}
	if backendIdx >= appIdx {
		t.Error("test-setup startup order wrong: api should come before mobile")
	}
}

// --- AC-4: pure cli/tui combo does NOT generate test-setup ---

func TestGenerateAggregate_PureCLI_TUI_NoTestSetup(t *testing.T) {
	surfaces := []SurfaceEntry{
		{Key: "", Type: types.SurfaceCLI},
		{Key: "tools", Type: types.SurfaceTUI},
	}
	out, err := GenerateAggregate(surfaces)
	if err != nil {
		t.Fatalf("GenerateAggregate: %v", err)
	}

	if strings.Contains(out, "test-setup") {
		t.Error("pure cli/tui aggregate should NOT contain test-setup recipe")
	}
}

func TestGenerateAggregate_SingleService_NoTestSetup(t *testing.T) {
	surfaces := []SurfaceEntry{
		{Key: "backend", Type: types.SurfaceAPI},
	}
	out, err := GenerateAggregate(surfaces)
	if err != nil {
		t.Fatalf("GenerateAggregate: %v", err)
	}

	if strings.Contains(out, "test-setup") {
		t.Error("single service surface should NOT generate test-setup")
	}
}

// --- Edge cases ---

func TestGenerateAggregate_EmptySurfaces(t *testing.T) {
	out, err := GenerateAggregate(nil)
	if err != nil {
		t.Fatalf("GenerateAggregate(nil): %v", err)
	}

	// Should still produce install, ci, clean (empty bodies)
	for _, name := range []string{"install", "ci", "clean"} {
		if !strings.Contains(out, name+" [unix]:") {
			t.Errorf("empty aggregate missing recipe %q", name)
		}
	}
	if strings.Contains(out, "test-setup") {
		t.Error("empty surfaces should not generate test-setup")
	}
}

func TestGenerateAggregate_ScalarCLI_NoPrefix(t *testing.T) {
	// Scalar CLI surface: key is "", recipes have no prefix
	surfaces := []SurfaceEntry{
		{Key: "", Type: types.SurfaceCLI},
	}
	out, err := GenerateAggregate(surfaces)
	if err != nil {
		t.Fatalf("GenerateAggregate: %v", err)
	}

	// ci should reference "just lint", "just compile", "just unit-test" (no prefix)
	for _, verb := range []string{"lint", "compile", "unit-test"} {
		if !strings.Contains(out, "just "+verb) {
			t.Errorf("scalar ci recipe missing 'just %s'", verb)
		}
	}
}

func TestGenerateAggregate_InstallAggregatesAllInstalls(t *testing.T) {
	surfaces := []SurfaceEntry{
		{Key: "backend", Type: types.SurfaceAPI},
		{Key: "frontend", Type: types.SurfaceWeb},
	}
	out, err := GenerateAggregate(surfaces)
	if err != nil {
		t.Fatalf("GenerateAggregate: %v", err)
	}

	if !strings.Contains(out, "just backend-install") {
		t.Error("install recipe missing 'just backend-install'")
	}
	if !strings.Contains(out, "just frontend-install") {
		t.Error("install recipe missing 'just frontend-install'")
	}
}

func TestGenerateAggregate_CleanAggregatesAllCleans(t *testing.T) {
	surfaces := []SurfaceEntry{
		{Key: "backend", Type: types.SurfaceAPI},
		{Key: "frontend", Type: types.SurfaceWeb},
	}
	out, err := GenerateAggregate(surfaces)
	if err != nil {
		t.Fatalf("GenerateAggregate: %v", err)
	}

	if !strings.Contains(out, "just backend-clean") {
		t.Error("clean recipe missing 'just backend-clean'")
	}
	if !strings.Contains(out, "just frontend-clean") {
		t.Error("clean recipe missing 'just frontend-clean'")
	}
}

func TestGenerateAggregate_TestSetup_NotMarkedUserCustomized(t *testing.T) {
	surfaces := []SurfaceEntry{
		{Key: "backend", Type: types.SurfaceAPI},
		{Key: "frontend", Type: types.SurfaceWeb},
	}
	out, err := GenerateAggregate(surfaces)
	if err != nil {
		t.Fatalf("GenerateAggregate: %v", err)
	}

	marker := "# user-customized\ntest-setup [unix]:"
	if strings.Contains(out, marker) {
		t.Error("test-setup aggregate should NOT have '# user-customized' marker")
	}
}

// --- Coverage: surfacesToEntries ---

func TestSurfacesToEntries_ScalarDot(t *testing.T) {
	surfaces := map[string]string{".": "cli"}
	entries := surfacesToEntries(surfaces)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Key != "" {
		t.Errorf("scalar key should be empty, got %q", entries[0].Key)
	}
	if entries[0].Type != types.SurfaceCLI {
		t.Errorf("expected cli type, got %q", entries[0].Type)
	}
}

func TestSurfacesToEntries_Empty(t *testing.T) {
	entries := surfacesToEntries(nil)
	if entries != nil {
		t.Errorf("expected nil for nil input, got %v", entries)
	}
	entries = surfacesToEntries(map[string]string{})
	if entries != nil {
		t.Errorf("expected nil for empty map, got %v", entries)
	}
}

func TestSurfacesToEntries_Multiple(t *testing.T) {
	surfaces := map[string]string{
		"frontend": "web",
		"backend":  "api",
	}
	entries := surfacesToEntries(surfaces)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	// Should be sorted by key: backend < frontend
	if entries[0].Key != "backend" {
		t.Errorf("first entry key should be 'backend', got %q", entries[0].Key)
	}
	if entries[1].Key != "frontend" {
		t.Errorf("second entry key should be 'frontend', got %q", entries[1].Key)
	}
}

// --- Coverage: isServiceType / hasMultipleServiceSurfaces ---

func TestIsServiceType(t *testing.T) {
	serviceTypes := []types.SurfaceType{types.SurfaceAPI, types.SurfaceWeb, types.SurfaceMobile}
	for _, typ := range serviceTypes {
		if !isServiceType(typ) {
			t.Errorf("isServiceType(%q) = false, want true", typ)
		}
	}
	nonServiceTypes := []types.SurfaceType{types.SurfaceCLI, types.SurfaceTUI}
	for _, typ := range nonServiceTypes {
		if isServiceType(typ) {
			t.Errorf("isServiceType(%q) = true, want false", typ)
		}
	}
}

// --- Coverage: mixed surfaces (cli + api) only 1 service → no test-setup ---

func TestGenerateAggregate_MixedCLI_API_NoTestSetup(t *testing.T) {
	surfaces := []SurfaceEntry{
		{Key: "", Type: types.SurfaceCLI},
		{Key: "backend", Type: types.SurfaceAPI},
	}
	out, err := GenerateAggregate(surfaces)
	if err != nil {
		t.Fatalf("GenerateAggregate: %v", err)
	}
	if strings.Contains(out, "test-setup") {
		t.Error("1 service + 1 cli should NOT generate test-setup")
	}
}
