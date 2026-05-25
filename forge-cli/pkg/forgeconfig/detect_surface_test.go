package forgeconfig

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// --- AC1: Detects surface types from various manifest files ---

func TestDetectSurfaces_PackageJSON(t *testing.T) {
	tests := []struct {
		name     string
		deps     map[string]string // dependencies
		devDeps  map[string]string // devDependencies
		wantType string
	}{
		{
			name:     "react -> web",
			deps:     map[string]string{"react": "^18.0.0"},
			wantType: "web",
		},
		{
			name:     "vue -> web",
			deps:     map[string]string{"vue": "^3.0.0"},
			wantType: "web",
		},
		{
			name:     "svelte -> web",
			devDeps:  map[string]string{"svelte": "^4.0.0"},
			wantType: "web",
		},
		{
			name:     "express (no frontend) -> api",
			deps:     map[string]string{"express": "^4.18.0"},
			wantType: "api",
		},
		{
			name:     "fastify (no frontend) -> api",
			deps:     map[string]string{"fastify": "^4.0.0"},
			wantType: "api",
		},
		{
			name:     "commander (no frontend) -> cli",
			deps:     map[string]string{"commander": "^11.0.0"},
			wantType: "cli",
		},
		{
			name:     "yargs (no frontend) -> cli",
			deps:     map[string]string{"yargs": "^17.0.0"},
			wantType: "cli",
		},
		{
			name:     "blessed -> tui",
			deps:     map[string]string{"blessed": "^0.1.81"},
			wantType: "tui",
		},
		{
			name:     "ink -> tui",
			deps:     map[string]string{"ink": "^4.0.0"},
			wantType: "tui",
		},
		{
			name:     "react + express conflict -> web (higher priority)",
			deps:     map[string]string{"react": "^18.0.0", "express": "^4.18.0"},
			wantType: "web",
		},
		{
			name:     "react-native + react conflict -> mobile",
			deps:     map[string]string{"react-native": "^0.72.0", "react": "^18.0.0"},
			wantType: "mobile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writePackageJSON(t, dir, tt.deps, tt.devDeps)

			result, err := DetectSurfaces(dir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 {
				t.Fatalf("expected 1 surface, got %d: %v", len(result), result)
			}
			if result["."] != tt.wantType {
				t.Errorf("expected surfaces['.']='%s', got '%s'", tt.wantType, result["."])
			}
		})
	}
}

func TestDetectSurfaces_GoMod(t *testing.T) {
	tests := []struct {
		name     string
		require  []string // lines in require block
		wantType string
	}{
		{
			name:     "gin -> api",
			require:  []string{"github.com/gin-gonic/gin v1.9.0"},
			wantType: "api",
		},
		{
			name:     "echo -> api",
			require:  []string{"github.com/labstack/echo/v4 v4.11.0"},
			wantType: "api",
		},
		{
			name:     "cobra -> cli",
			require:  []string{"github.com/spf13/cobra v1.7.0"},
			wantType: "cli",
		},
		{
			name:     "bubbletea -> tui",
			require:  []string{"github.com/charmbracelet/bubbletea v0.24.0"},
			wantType: "tui",
		},
		{
			name:     "gin + cobra conflict -> api (higher priority)",
			require:  []string{"github.com/gin-gonic/gin v1.9.0", "github.com/spf13/cobra v1.7.0"},
			wantType: "api",
		},
		{
			name:     "bubbletea + cobra conflict -> tui (TUI is primary surface)",
			require:  []string{"github.com/charmbracelet/bubbletea v1.3.0", "github.com/spf13/cobra v1.10.0"},
			wantType: "tui",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writeGoMod(t, dir, tt.require)

			result, err := DetectSurfaces(dir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 {
				t.Fatalf("expected 1 surface, got %d: %v", len(result), result)
			}
			if result["."] != tt.wantType {
				t.Errorf("expected surfaces['.']='%s', got '%s'", tt.wantType, result["."])
			}
		})
	}
}

func TestDetectSurfaces_CargoToml(t *testing.T) {
	tests := []struct {
		name     string
		deps     []string // dependency lines
		wantType string
	}{
		{
			name:     "actix -> api",
			deps:     []string{"actix-web = \"4\""},
			wantType: "api",
		},
		{
			name:     "axum -> api",
			deps:     []string{"axum = \"0.7\""},
			wantType: "api",
		},
		{
			name:     "clap -> cli",
			deps:     []string{"clap = { version = \"4\", features = [\"derive\"] }"},
			wantType: "cli",
		},
		{
			name:     "ratatui -> tui",
			deps:     []string{"ratatui = \"0.25\""},
			wantType: "tui",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writeCargoToml(t, dir, tt.deps)

			result, err := DetectSurfaces(dir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 {
				t.Fatalf("expected 1 surface, got %d: %v", len(result), result)
			}
			if result["."] != tt.wantType {
				t.Errorf("expected surfaces['.']='%s', got '%s'", tt.wantType, result["."])
			}
		})
	}
}

func TestDetectSurfaces_MobileProjects(t *testing.T) {
	t.Run("AndroidManifest.xml -> mobile", func(t *testing.T) {
		dir := t.TempDir()
		writeAndroidManifest(t, dir)

		result, err := DetectSurfaces(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result["."] != "mobile" {
			t.Errorf("expected 'mobile', got '%s'", result["."])
		}
	})

	t.Run("xcodeproj -> mobile", func(t *testing.T) {
		dir := t.TempDir()
		xcodeDir := filepath.Join(dir, "MyApp.xcodeproj")
		if err := os.MkdirAll(xcodeDir, 0755); err != nil {
			t.Fatal(err)
		}
		// Create project.pbxproj to make it look real
		if err := os.WriteFile(filepath.Join(xcodeDir, "project.pbxproj"), []byte("{}"), 0644); err != nil {
			t.Fatal(err)
		}

		result, err := DetectSurfaces(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result["."] != "mobile" {
			t.Errorf("expected 'mobile', got '%s'", result["."])
		}
	})
}

func TestDetectSurfaces_PyProject(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantType string
	}{
		{
			name:     "flask -> api",
			content:  "[project]\ndependencies = [\"flask\"]\n",
			wantType: "api",
		},
		{
			name:     "fastapi -> api",
			content:  "[project]\ndependencies = [\"fastapi\"]\n",
			wantType: "api",
		},
		{
			name:     "click -> cli",
			content:  "[project]\ndependencies = [\"click\"]\n",
			wantType: "cli",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "pyproject.toml"), []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			result, err := DetectSurfaces(dir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 || result["."] != tt.wantType {
				t.Errorf("expected surfaces['.']='%s', got '%s'", tt.wantType, result["."])
			}
		})
	}
}

// --- AC2: Workspace detection ---

func TestDetectSurfaces_WorkspacePnpm(t *testing.T) {
	dir := t.TempDir()

	// Root package.json (workspace config only)
	writePackageJSON(t, dir, map[string]string{"express": "^4.18.0"}, nil)
	// pnpm-workspace.yaml triggers workspace mode
	writeTestFile(t, dir, "pnpm-workspace.yaml", "packages:\n  - 'apps/*'\n  - 'packages/*'\n")

	// Subdir: apps/web has react
	webDir := filepath.Join(dir, "apps", "web")
	mkdirAll(t, webDir)
	writePackageJSON(t, webDir, map[string]string{"react": "^18.0.0"}, nil)

	// Subdir: apps/api has express
	apiDir := filepath.Join(dir, "apps", "api")
	mkdirAll(t, apiDir)
	writePackageJSON(t, apiDir, map[string]string{"express": "^4.18.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Root deps should be skipped; only subdirs detected
	if _, ok := result["."]; ok {
		t.Error("root deps should be skipped in workspace mode")
	}
	webKey := filepath.ToSlash(filepath.Join("apps", "web"))
	apiKey := filepath.ToSlash(filepath.Join("apps", "api"))
	if result[webKey] != "web" {
		t.Errorf("expected apps/web=web, got %v", result)
	}
	if result[apiKey] != "api" {
		t.Errorf("expected apps/api=api, got %v", result)
	}
}

func TestDetectSurfaces_WorkspacePackageJSON(t *testing.T) {
	dir := t.TempDir()

	// Root package.json with workspaces field
	writePackageJSONWithWorkspaces(t, dir, []string{"apps/*", "packages/*"})

	// Subdir
	cliDir := filepath.Join(dir, "packages", "cli")
	mkdirAll(t, cliDir)
	writePackageJSON(t, cliDir, map[string]string{"commander": "^11.0.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cliKey := filepath.ToSlash(filepath.Join("packages", "cli"))
	if result[cliKey] != "cli" {
		t.Errorf("expected packages/cli=cli, got %v", result)
	}
}

// --- AC3: Non-workspace root detected as "." ---

func TestDetectSurfaces_NonWorkspace_ScalarOutput(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{"react": "^18.0.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Single type -> scalar form (key is ".")
	if len(result) != 1 {
		t.Fatalf("expected 1 surface, got %d", len(result))
	}
	if result["."] != "web" {
		t.Errorf("expected surfaces['.']='web', got '%s'", result["."])
	}
}

// --- AC4: Depth limit ---

func TestDetectSurfaces_DefaultDepth3(t *testing.T) {
	dir := t.TempDir()
	writePackageJSONWithWorkspaces(t, dir, []string{"apps/*"})

	// Create at depth 3: apps/web/deep (3 levels from root)
	deepDir := filepath.Join(dir, "apps", "web", "deep")
	mkdirAll(t, deepDir)
	writePackageJSON(t, deepDir, map[string]string{"react": "^18.0.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	deepKey := filepath.ToSlash(filepath.Join("apps", "web", "deep"))
	if result[deepKey] != "web" {
		t.Errorf("expected to detect at depth 3, got %v", result)
	}
}

func TestDetectSurfaces_DepthExceeded(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("FORGE_DETECT_DEPTH", "2")

	writePackageJSONWithWorkspaces(t, dir, []string{"apps/*"})

	// Create at depth 3: should NOT be detected
	deepDir := filepath.Join(dir, "apps", "web", "deep")
	mkdirAll(t, deepDir)
	writePackageJSON(t, deepDir, map[string]string{"react": "^18.0.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	deepKey := filepath.ToSlash(filepath.Join("apps", "web", "deep"))
	if _, ok := result[deepKey]; ok {
		t.Errorf("should NOT detect at depth 3 with FORGE_DETECT_DEPTH=2, got %v", result)
	}
}

func TestDetectSurfaces_DepthZeroError(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("FORGE_DETECT_DEPTH", "0")

	writePackageJSON(t, dir, map[string]string{"react": "^18.0.0"}, nil)

	_, err := DetectSurfaces(dir)
	if err == nil {
		t.Fatal("expected error for FORGE_DETECT_DEPTH=0")
	}
	if _, ok := err.(*ErrInvalidDepth); !ok {
		t.Errorf("expected ErrInvalidDepth, got %T: %v", err, err)
	}
}

func TestDetectSurfaces_DepthNegativeError(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("FORGE_DETECT_DEPTH", "-1")

	writePackageJSON(t, dir, map[string]string{"react": "^18.0.0"}, nil)

	_, err := DetectSurfaces(dir)
	if err == nil {
		t.Fatal("expected error for FORGE_DETECT_DEPTH=-1")
	}
}

func TestDetectSurfaces_DepthEnvOverride(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("FORGE_DETECT_DEPTH", "1")

	writePackageJSONWithWorkspaces(t, dir, []string{"apps/*"})

	// apps/web is depth 2 from root, exceeds FORGE_DETECT_DEPTH=1
	skipDir := filepath.Join(dir, "apps", "web")
	mkdirAll(t, skipDir)
	writePackageJSON(t, skipDir, map[string]string{"react": "^18.0.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	skipKey := filepath.ToSlash(filepath.Join("apps", "web"))
	if _, ok := result[skipKey]; ok {
		t.Errorf("should NOT detect at depth 2 with FORGE_DETECT_DEPTH=1, got %v", result)
	}
}

// --- AC5: Exclusion dirs ---

func TestDetectSurfaces_ExclusionDirs(t *testing.T) {
	dir := t.TempDir()
	writePackageJSONWithWorkspaces(t, dir, []string{"apps/*"})

	// node_modules should be skipped
	nmDir := filepath.Join(dir, "node_modules", "pkg")
	mkdirAll(t, nmDir)
	writePackageJSON(t, nmDir, map[string]string{"react": "^18.0.0"}, nil)

	// .git should be skipped
	gitDir := filepath.Join(dir, ".git", "subdir")
	mkdirAll(t, gitDir)
	writePackageJSON(t, gitDir, map[string]string{"express": "^4.18.0"}, nil)

	// dist should be skipped
	distDir := filepath.Join(dir, "dist", "pkg")
	mkdirAll(t, distDir)
	writePackageJSON(t, distDir, map[string]string{"commander": "^11.0.0"}, nil)

	// tests should be skipped (dedicated to advanced test cases per Forge conventions)
	testsDir := filepath.Join(dir, "tests", "e2e")
	mkdirAll(t, testsDir)
	writePackageJSON(t, testsDir, map[string]string{"react": "^18.0.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) > 0 {
		t.Errorf("expected 0 surfaces (all dirs excluded), got %d: %v", len(result), result)
	}
}

func TestDetectSurfaces_ExclusionDirList(t *testing.T) {
	expectedExclusions := []string{
		"node_modules", ".git", "vendor", "dist",
		"build", "__pycache__", ".next", "target",
		"tests",
	}

	for _, d := range expectedExclusions {
		if !isExcludedDir(d) {
			t.Errorf("expected %q to be excluded", d)
		}
	}
}

// --- AC6: Signal conflict resolution ---

func TestDetectSurfaces_ConflictResolution(t *testing.T) {
	tests := []struct {
		name     string
		deps     map[string]string
		wantType string
	}{
		{
			name:     "web > api (react + express)",
			deps:     map[string]string{"react": "^18.0.0", "express": "^4.18.0"},
			wantType: "web",
		},
		{
			name:     "web > cli (react + commander)",
			deps:     map[string]string{"react": "^18.0.0", "commander": "^11.0.0"},
			wantType: "web",
		},
		{
			name:     "api > cli (express + commander)",
			deps:     map[string]string{"express": "^4.18.0", "commander": "^11.0.0"},
			wantType: "api",
		},
		{
			name:     "api > tui (express + blessed)",
			deps:     map[string]string{"express": "^4.18.0", "blessed": "^0.1.81"},
			wantType: "api",
		},
		{
			name:     "tui > cli (blessed + commander)",
			deps:     map[string]string{"commander": "^11.0.0", "blessed": "^0.1.81"},
			wantType: "tui",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writePackageJSON(t, dir, tt.deps, nil)

			result, err := DetectSurfaces(dir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result["."] != tt.wantType {
				t.Errorf("expected '%s', got '%s'", tt.wantType, result["."])
			}
		})
	}
}

// --- AC7: Single-type -> scalar; multi-type -> map ---

func TestDetectSurfaces_SingleTypeScalar(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{"react": "^18.0.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 surface, got %d", len(result))
	}
	if result["."] != "web" {
		t.Errorf("expected scalar form with key '.', got %v", result)
	}
}

// bug: scalar key not normalized to "." when only subdir has signal
func TestDetectSurfaces_SubdirOnly_ScalarKeyNormalized(t *testing.T) {
	dir := t.TempDir()
	// Root has NO go.mod — only a subdirectory does
	cliDir := filepath.Join(dir, "my-cli")
	mkdirAll(t, cliDir)
	writeGoMod(t, cliDir, []string{"github.com/spf13/cobra v1.10.0"})

	result, err := DetectSurfacesWithConflicts(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsScalar {
		t.Fatalf("expected IsScalar=true, got false; surfaces=%v", result.Surfaces)
	}
	if result.Surfaces["."] != "cli" {
		t.Errorf("expected surfaces['.']='cli', got surfaces=%v", result.Surfaces)
	}
}

func TestDetectSurfaces_MultiTypeMap(t *testing.T) {
	dir := t.TempDir()
	writePackageJSONWithWorkspaces(t, dir, []string{"apps/*"})

	webDir := filepath.Join(dir, "apps", "frontend")
	mkdirAll(t, webDir)
	writePackageJSON(t, webDir, map[string]string{"react": "^18.0.0"}, nil)

	apiDir := filepath.Join(dir, "apps", "backend")
	mkdirAll(t, apiDir)
	writePackageJSON(t, apiDir, map[string]string{"express": "^4.18.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 surfaces, got %d: %v", len(result), result)
	}

	frontendKey := filepath.ToSlash(filepath.Join("apps", "frontend"))
	backendKey := filepath.ToSlash(filepath.Join("apps", "backend"))
	if result[frontendKey] != "web" {
		t.Errorf("expected apps/frontend=web, got %v", result)
	}
	if result[backendKey] != "api" {
		t.Errorf("expected apps/backend=api, got %v", result)
	}
}

// --- Edge cases ---

func TestDetectSurfaces_NoManifests(t *testing.T) {
	dir := t.TempDir()

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result for empty dir, got %v", result)
	}
}

func TestDetectSurfaces_UnknownDepsIgnored(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{"lodash": "^4.17.0", "axios": "^1.0.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result for unknown deps, got %v", result)
	}
}

// --- Structural inference tests ---

func TestInferGoSurface(t *testing.T) {
	tests := []struct {
		name       string
		setupDir   func(t *testing.T, dir string)
		wantType   string
		wantSource string
	}{
		{
			name: "cmd subdirectories -> cli",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				mkdirAll(t, filepath.Join(dir, "cmd", "myapp"))
				writeGoMod(t, dir, nil)
			},
			wantType:   "cli",
			wantSource: "inference:cmd-dir",
		},
		{
			name: "api directory -> api",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				mkdirAll(t, filepath.Join(dir, "api"))
				writeGoMod(t, dir, nil)
			},
			wantType:   "api",
			wantSource: "inference:api-dir",
		},
		{
			name: "handler directory -> api",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				mkdirAll(t, filepath.Join(dir, "handler"))
				writeGoMod(t, dir, nil)
			},
			wantType:   "api",
			wantSource: "inference:handler-dir",
		},
		{
			name: "both cmd and api present -> api wins, cli discarded",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				mkdirAll(t, filepath.Join(dir, "cmd", "myapp"))
				mkdirAll(t, filepath.Join(dir, "api"))
				writeGoMod(t, dir, nil)
			},
			wantType:   "api",
			wantSource: "inference:api-dir",
		},
		{
			name: "both cmd and handler present -> api wins",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				mkdirAll(t, filepath.Join(dir, "cmd", "myapp"))
				mkdirAll(t, filepath.Join(dir, "handler"))
				writeGoMod(t, dir, nil)
			},
			wantType:   "api",
			wantSource: "inference:handler-dir",
		},
		{
			name: "cmd file (not directory) -> ignored",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writeTestFile(t, dir, "cmd", "#!/bin/sh\necho hello")
				writeGoMod(t, dir, nil)
			},
			wantType:   "",
			wantSource: "",
		},
		{
			name: "no matching directories -> empty",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				mkdirAll(t, filepath.Join(dir, "internal"))
				mkdirAll(t, filepath.Join(dir, "pkg"))
				writeGoMod(t, dir, nil)
			},
			wantType:   "",
			wantSource: "",
		},
		{
			name: "no go.mod -> empty",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				mkdirAll(t, filepath.Join(dir, "cmd", "myapp"))
			},
			wantType:   "",
			wantSource: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			tt.setupDir(t, dir)
			gotType, gotSource := inferGoSurface(dir)
			if gotType != tt.wantType {
				t.Errorf("inferGoSurface() type = %q, want %q", gotType, tt.wantType)
			}
			if gotSource != tt.wantSource {
				t.Errorf("inferGoSurface() source = %q, want %q", gotSource, tt.wantSource)
			}
		})
	}
}

func TestInferNodeSurface(t *testing.T) {
	tests := []struct {
		name       string
		setupDir   func(t *testing.T, dir string)
		wantType   string
		wantSource string
	}{
		{
			name: "bin field (string form) -> cli",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writePackageJSONWithBin(t, dir, `"./bin/cli.js"`)
			},
			wantType:   "cli",
			wantSource: "inference:bin-field",
		},
		{
			name: "bin field (object form) -> cli",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writePackageJSONWithBin(t, dir, `{"myapp": "./bin/cli.js"}`)
			},
			wantType:   "cli",
			wantSource: "inference:bin-field",
		},
		{
			name: "index.html at root -> web",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writeTestFile(t, dir, "package.json", `{}`)
				writeTestFile(t, dir, "index.html", "<html></html>")
			},
			wantType:   "web",
			wantSource: "inference:index-html",
		},
		{
			name: "both bin and index.html -> web wins (higher priority)",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writePackageJSONWithBin(t, dir, `"./bin/cli.js"`)
				writeTestFile(t, dir, "index.html", "<html></html>")
			},
			wantType:   "web",
			wantSource: "inference:index-html",
		},
		{
			name: "index.html in subdir (not root) -> NOT web",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writeTestFile(t, dir, "package.json", `{}`)
				subDir := filepath.Join(dir, "public")
				mkdirAll(t, subDir)
				writeTestFile(t, subDir, "index.html", "<html></html>")
			},
			wantType:   "",
			wantSource: "",
		},
		{
			name: "no package.json -> empty",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writeTestFile(t, dir, "index.html", "<html></html>")
			},
			wantType:   "",
			wantSource: "",
		},
		{
			name: "package.json without bin or index.html -> empty",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writeTestFile(t, dir, "package.json", `{"name": "test"}`)
			},
			wantType:   "",
			wantSource: "",
		},
		{
			name: "malformed package.json -> empty, no crash",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writeTestFile(t, dir, "package.json", `{invalid json`)
			},
			wantType:   "",
			wantSource: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			tt.setupDir(t, dir)
			gotType, gotSource := inferNodeSurface(dir)
			if gotType != tt.wantType {
				t.Errorf("inferNodeSurface() type = %q, want %q", gotType, tt.wantType)
			}
			if gotSource != tt.wantSource {
				t.Errorf("inferNodeSurface() source = %q, want %q", gotSource, tt.wantSource)
			}
		})
	}
}

func TestInferPythonSurface(t *testing.T) {
	tests := []struct {
		name       string
		setupDir   func(t *testing.T, dir string)
		wantType   string
		wantSource string
	}{
		{
			name: "pyproject.toml with [project.scripts] -> cli",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writeTestFile(t, dir, "pyproject.toml", "[project]\nname = \"myapp\"\n\n[project.scripts]\nmyapp = \"myapp.cli:main\"\n")
			},
			wantType:   "cli",
			wantSource: "inference:py-scripts",
		},
		{
			name: "setup.py with entry_points -> cli",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writeTestFile(t, dir, "setup.py", "from setuptools import setup\nsetup(\n    name='myapp',\n    entry_points={\n        'console_scripts': ['myapp=myapp.cli:main'],\n    },\n)\n")
			},
			wantType:   "cli",
			wantSource: "inference:py-scripts",
		},
		{
			name: "app.py at root (no setup.py, no library markers) -> cli",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writeTestFile(t, dir, "app.py", "print('hello')")
			},
			wantType:   "cli",
			wantSource: "inference:py-main",
		},
		{
			name: "main.py at root (no setup.py, no library markers) -> cli",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writeTestFile(t, dir, "main.py", "print('hello')")
			},
			wantType:   "cli",
			wantSource: "inference:py-main",
		},
		{
			name: "app.py with setup.py having matching name -> NOT cli",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				dirName := filepath.Base(dir)
				writeTestFile(t, dir, "app.py", "print('hello')")
				writeTestFile(t, dir, "setup.py", "from setuptools import setup\nsetup(name='"+dirName+"')\n")
			},
			wantType:   "",
			wantSource: "",
		},
		{
			name: "app.py with [project.packages] in pyproject.toml -> NOT cli",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writeTestFile(t, dir, "app.py", "print('hello')")
				writeTestFile(t, dir, "pyproject.toml", "[project]\nname = \"mylib\"\n\n[project.packages]\nfind = {}\n")
			},
			wantType:   "",
			wantSource: "",
		},
		{
			name: "app.py with [tool.setuptools.packages.find] -> NOT cli",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writeTestFile(t, dir, "app.py", "print('hello')")
				writeTestFile(t, dir, "pyproject.toml", "[project]\nname = \"mylib\"\n\n[tool.setuptools.packages.find]\nwhere = [\"src\"]\n")
			},
			wantType:   "",
			wantSource: "",
		},
		{
			name: "app.py with setup.py but name does NOT match -> cli",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writeTestFile(t, dir, "app.py", "print('hello')")
				writeTestFile(t, dir, "setup.py", "from setuptools import setup\nsetup(name='different-name')\n")
			},
			wantType:   "cli",
			wantSource: "inference:py-main",
		},
		{
			name: "no Python markers -> empty",
			setupDir: func(t *testing.T, _ string) {
				t.Helper()
				// Empty dir, no manifests, no python files
			},
			wantType:   "",
			wantSource: "",
		},
		{
			name: "malformed pyproject.toml -> empty, no crash",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writeTestFile(t, dir, "pyproject.toml", "[this is [broken toml\n")
			},
			wantType:   "",
			wantSource: "",
		},
		{
			name: "malformed setup.py -> empty, no crash",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writeTestFile(t, dir, "setup.py", "this is not valid python {{{{")
			},
			wantType:   "",
			wantSource: "",
		},
		{
			name: "scripts in pyproject.toml takes priority over app.py",
			setupDir: func(t *testing.T, dir string) {
				t.Helper()
				writeTestFile(t, dir, "app.py", "print('hello')")
				writeTestFile(t, dir, "pyproject.toml", "[project]\nname = \"myapp\"\n\n[project.scripts]\nmyapp = \"myapp.cli:main\"\n")
			},
			wantType:   "cli",
			wantSource: "inference:py-scripts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			tt.setupDir(t, dir)
			gotType, gotSource := inferPythonSurface(dir)
			if gotType != tt.wantType {
				t.Errorf("inferPythonSurface() type = %q, want %q", gotType, tt.wantType)
			}
			if gotSource != tt.wantSource {
				t.Errorf("inferPythonSurface() source = %q, want %q", gotSource, tt.wantSource)
			}
		})
	}
}

// --- Priority chain: inference only when dependency signals empty ---

func TestInferencePriorityChain(t *testing.T) {
	t.Run("dependency signals present -> inference NOT called", func(t *testing.T) {
		dir := t.TempDir()
		// cobra dependency -> cli (dependency signal)
		writeGoMod(t, dir, []string{"github.com/spf13/cobra v1.7.0"})
		// Structural: cmd/ subdir would also infer cli
		mkdirAll(t, filepath.Join(dir, "cmd", "myapp"))

		result, err := DetectSurfacesWithConflicts(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Surfaces["."] != "cli" {
			t.Errorf("expected cli, got %q", result.Surfaces["."])
		}
		// Source should be dependency, not inference
		if result.Sources["."] != "dependency:cobra" {
			t.Errorf("expected source 'dependency:cobra', got %q", result.Sources["."])
		}
	})

	t.Run("dependency signals empty -> inference called", func(t *testing.T) {
		dir := t.TempDir()
		// go.mod with no known frameworks
		writeGoMod(t, dir, nil)
		// cmd/ subdir -> structural inference
		mkdirAll(t, filepath.Join(dir, "cmd", "myapp"))

		result, err := DetectSurfacesWithConflicts(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Surfaces["."] != "cli" {
			t.Errorf("expected cli, got %q", result.Surfaces["."])
		}
		if result.Sources["."] != "inference:cmd-dir" {
			t.Errorf("expected source 'inference:cmd-dir', got %q", result.Sources["."])
		}
	})

	t.Run("no signals no inference -> empty result", func(t *testing.T) {
		dir := t.TempDir()
		// Empty dir, no manifests
		result, err := DetectSurfacesWithConflicts(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Surfaces) != 0 {
			t.Errorf("expected empty surfaces, got %v", result.Surfaces)
		}
		if len(result.Sources) != 0 {
			t.Errorf("expected empty sources, got %v", result.Sources)
		}
	})
}

// --- Sources map population tests ---

func TestSourcesMapPopulation(t *testing.T) {
	t.Run("dependency detection populates Sources", func(t *testing.T) {
		dir := t.TempDir()
		writePackageJSON(t, dir, map[string]string{"react": "^18.0.0"}, nil)

		result, err := DetectSurfacesWithConflicts(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Sources["."] != "dependency:react" {
			t.Errorf("expected source 'dependency:react', got %q", result.Sources["."])
		}
	})

	t.Run("Go dependency detection populates Sources", func(t *testing.T) {
		dir := t.TempDir()
		writeGoMod(t, dir, []string{"github.com/spf13/cobra v1.7.0"})

		result, err := DetectSurfacesWithConflicts(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Sources["."] != "dependency:cobra" {
			t.Errorf("expected source 'dependency:cobra', got %q", result.Sources["."])
		}
	})

	t.Run("inference populates Sources", func(t *testing.T) {
		dir := t.TempDir()
		writeGoMod(t, dir, nil)
		mkdirAll(t, filepath.Join(dir, "cmd", "myapp"))

		result, err := DetectSurfacesWithConflicts(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Sources["."] != "inference:cmd-dir" {
			t.Errorf("expected source 'inference:cmd-dir', got %q", result.Sources["."])
		}
	})

	t.Run("Node.js inference populates Sources for bin field", func(t *testing.T) {
		dir := t.TempDir()
		writePackageJSONWithBin(t, dir, `"./bin/cli.js"`)

		// Empty deps so dependency signals are empty, inference fires
		// But wait: package.json with bin also needs no dependency signals
		result, err := DetectSurfacesWithConflicts(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Surfaces["."] != "cli" {
			t.Errorf("expected cli, got %q", result.Surfaces["."])
		}
		if result.Sources["."] != "inference:bin-field" {
			t.Errorf("expected source 'inference:bin-field', got %q", result.Sources["."])
		}
	})
}

// --- Filesystem error resilience ---

func TestInferenceFilesystemResilience(t *testing.T) {
	t.Run("unreadable directory returns empty without panic", func(t *testing.T) {
		// Test inferGoSurface with a non-existent directory
		gotType, gotSource := inferGoSurface("/nonexistent/path/that/does/not/exist")
		if gotType != "" || gotSource != "" {
			t.Errorf("expected empty result for non-existent dir, got type=%q source=%q", gotType, gotSource)
		}
	})

	t.Run("inferNodeSurface on non-existent dir returns empty", func(t *testing.T) {
		gotType, gotSource := inferNodeSurface("/nonexistent/path/that/does/not/exist")
		if gotType != "" || gotSource != "" {
			t.Errorf("expected empty result for non-existent dir, got type=%q source=%q", gotType, gotSource)
		}
	})

	t.Run("inferPythonSurface on non-existent dir returns empty", func(t *testing.T) {
		gotType, gotSource := inferPythonSurface("/nonexistent/path/that/does/not/exist")
		if gotType != "" || gotSource != "" {
			t.Errorf("expected empty result for non-existent dir, got type=%q source=%q", gotType, gotSource)
		}
	})
}

// --- Backward compatibility: DetectResult zero-value ---

func TestDetectResultZeroValue(t *testing.T) {
	var r DetectResult
	if r.Sources != nil {
		t.Error("zero-value DetectResult.Sources should be nil for backward compatibility")
	}
	if r.Surfaces != nil {
		t.Error("zero-value DetectResult.Surfaces should be nil")
	}
	if r.Conflicts != nil {
		t.Error("zero-value DetectResult.Conflicts should be nil")
	}
}

// --- Helper functions ---

func writePackageJSON(t *testing.T, dir string, deps, devDeps map[string]string) {
	t.Helper()
	type pkg struct {
		Dependencies    map[string]string `json:"dependencies,omitempty"`
		DevDependencies map[string]string `json:"devDependencies,omitempty"`
	}
	p := pkg{
		Dependencies:    deps,
		DevDependencies: devDeps,
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), append(data, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
}

func writePackageJSONWithWorkspaces(t *testing.T, dir string, workspaces []string) {
	t.Helper()
	type pkg struct {
		Workspaces []string `json:"workspaces"`
	}
	p := pkg{Workspaces: workspaces}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), append(data, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeGoMod(t *testing.T, dir string, require []string) {
	t.Helper()
	content := "module example.com/test\n\ngo 1.21\n"
	if len(require) > 0 {
		content += "\nrequire (\n"
		for _, r := range require {
			content += "\t" + r + "\n"
		}
		content += ")\n"
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeCargoToml(t *testing.T, dir string, deps []string) {
	t.Helper()
	content := "[package]\nname = \"test\"\nversion = \"0.1.0\"\n\n[dependencies]\n"
	for _, d := range deps {
		content += d + "\n"
	}
	if err := os.WriteFile(filepath.Join(dir, "Cargo.toml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeAndroidManifest(t *testing.T, dir string) {
	t.Helper()
	androidDir := filepath.Join(dir, "app", "src", "main")
	mkdirAll(t, androidDir)
	content := "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<manifest xmlns:android=\"http://schemas.android.com/apk/res/android\"\n    package=\"com.example.app\">\n</manifest>"
	if err := os.WriteFile(filepath.Join(androidDir, "AndroidManifest.xml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeTestFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func writePackageJSONWithBin(t *testing.T, dir string, binValue string) {
	t.Helper()
	content := "{\"name\": \"test\", \"bin\": " + binValue + "}\n"
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func mkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal(err)
	}
}
