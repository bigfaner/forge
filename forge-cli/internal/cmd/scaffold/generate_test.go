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

func TestValidateArgs_NamedWithKey_NotYetSupported(t *testing.T) {
	// api with key is valid for arg checking but not yet supported for generation
	err := ValidateArgs(types.SurfaceAPI, "backend")
	if err == nil {
		t.Fatal("expected error for not-yet-supported type")
	}
	if !strings.Contains(err.Error(), "not yet supported") {
		t.Errorf("unexpected error: %v", err)
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
