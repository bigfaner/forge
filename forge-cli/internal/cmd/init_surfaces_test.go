package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/forgeconfig"
)

// TestDetectSurfacesWithConflicts tests the enhanced detection that returns conflict metadata.
func TestDetectSurfacesWithConflicts(t *testing.T) {
	t.Run("single type returns scalar with IsScalar true", func(t *testing.T) {
		dir := t.TempDir()
		writePackageJSONDeps(t, dir, map[string]string{"react": "^18.0.0"})

		result, err := forgeconfig.DetectSurfacesWithConflicts(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.IsScalar {
			t.Error("expected IsScalar=true for single type")
		}
		if len(result.Surfaces) != 1 || result.Surfaces["."] != "web" {
			t.Errorf("expected {'.':'web'}, got %v", result.Surfaces)
		}
	})

	t.Run("conflict is recorded for web + api signals", func(t *testing.T) {
		dir := t.TempDir()
		writePackageJSONDeps(t, dir, map[string]string{
			"react":   "^18.0.0",
			"express": "^4.18.0",
		})

		result, err := forgeconfig.DetectSurfacesWithConflicts(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Conflicts) == 0 {
			t.Error("expected conflict for web + api signals")
		}
		c := result.Conflicts[0]
		if c.Resolved != "web" {
			t.Errorf("expected resolved=web, got %s", c.Resolved)
		}
		if len(c.Conflicting) < 2 {
			t.Errorf("expected at least 2 conflicting types, got %v", c.Conflicting)
		}
	})

	t.Run("no conflict for single signal", func(t *testing.T) {
		dir := t.TempDir()
		writePackageJSONDeps(t, dir, map[string]string{"react": "^18.0.0"})

		result, err := forgeconfig.DetectSurfacesWithConflicts(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Conflicts) != 0 {
			t.Errorf("expected no conflicts for single signal, got %d", len(result.Conflicts))
		}
	})

	t.Run("map form returns IsScalar false", func(t *testing.T) {
		dir := t.TempDir()
		writePackageJSONWorkspaces(t, dir, []string{"apps/*"})

		webDir := filepath.Join(dir, "apps", "frontend")
		mkdirAllTest(t, webDir)
		writePackageJSONDeps(t, webDir, map[string]string{"react": "^18.0.0"})

		apiDir := filepath.Join(dir, "apps", "backend")
		mkdirAllTest(t, apiDir)
		writePackageJSONDeps(t, apiDir, map[string]string{"express": "^4.18.0"})

		result, err := forgeconfig.DetectSurfacesWithConflicts(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.IsScalar {
			t.Error("expected IsScalar=false for map form")
		}
		if len(result.Surfaces) != 2 {
			t.Errorf("expected 2 surfaces, got %d", len(result.Surfaces))
		}
	})
}

// TestFormatConflictAnnotation tests the conflict annotation format.
func TestFormatConflictAnnotation(t *testing.T) {
	c := &forgeconfig.PathConflict{
		Path:        ".",
		Resolved:    "web",
		Conflicting: []string{"web", "api"},
	}
	annotation := formatConflictAnnotation(c)
	if !strings.Contains(annotation, "web + api") {
		t.Errorf("expected annotation to contain 'web + api', got %q", annotation)
	}
	if !strings.Contains(annotation, "web") {
		t.Errorf("expected annotation to mention resolved type 'web', got %q", annotation)
	}
}

// TestBuildDisplayLines tests the TUI display line builder.
func TestBuildDisplayLines(t *testing.T) {
	t.Run("scalar form shows only type", func(t *testing.T) {
		surfaces := forgeconfig.SurfacesMap{".": "api"}
		lines := buildDisplayLines(surfaces, nil)

		found := false
		for _, line := range lines {
			if strings.Contains(line, "api") {
				found = true
			}
		}
		if !found {
			t.Error("expected to find surface type in display lines")
		}
	})

	t.Run("map form shows path=surface", func(t *testing.T) {
		surfaces := forgeconfig.SurfacesMap{
			"frontend": "web",
			"backend":  "api",
		}
		lines := buildDisplayLines(surfaces, nil)

		joined := strings.Join(lines, "\n")
		if !strings.Contains(joined, "frontend:") {
			t.Errorf("expected 'frontend:' in display, got: %s", joined)
		}
		if !strings.Contains(joined, "backend:") {
			t.Errorf("expected 'backend:' in display, got: %s", joined)
		}
	})

	t.Run("conflict annotation displayed", func(t *testing.T) {
		surfaces := forgeconfig.SurfacesMap{"frontend": "web"}
		conflicts := map[string]*forgeconfig.PathConflict{
			"frontend": {
				Path:        "frontend",
				Resolved:    "web",
				Conflicting: []string{"web", "api"},
			},
		}
		lines := buildDisplayLines(surfaces, conflicts)
		joined := strings.Join(lines, "\n")
		if !strings.Contains(joined, "web + api") {
			t.Errorf("expected conflict annotation in display, got: %s", joined)
		}
	})
}

// TestRunSurfaceConfigIntegration tests the surface config step in init.
func TestRunSurfaceConfigIntegration(t *testing.T) {
	t.Run("skipped when no config file", func(t *testing.T) {
		dir := t.TempDir()
		action := runSurfaceConfig(dir)
		if action.status != "SKIPPED" {
			t.Errorf("expected SKIPPED, got %s", action.status)
		}
	})

	t.Run("test override for surface config step", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}

		// Create config file so surfaceConfigFunc is invoked
		configContent := "auto:\n  gitPush: false\n"
		if err := os.WriteFile(filepath.Join(forgeDir, feature.ForgeConfigFileName), []byte(configContent), 0o644); err != nil {
			t.Fatal(err)
		}

		// Override surfaceConfigFunc to avoid TTY requirement
		orig := surfaceConfigFunc
		surfaceConfigFunc = func(projectRoot string) initAction {
			cfg, err := forgeconfig.ReadConfig(projectRoot)
			if err != nil || cfg == nil {
				return initAction{status: "FAILED", target: "surfaces", detail: "read config"}
			}
			cfg.Surfaces = forgeconfig.SurfacesMap{".": "api"}
			configFile := filepath.Join(projectRoot, feature.ForgeDir, feature.ForgeConfigFileName)
			if err := writeConfigFile(configFile, cfg); err != nil {
				return initAction{status: "FAILED", target: "surfaces", detail: err.Error()}
			}
			return initAction{status: "CREATED", target: "surfaces", detail: "api"}
		}
		defer func() { surfaceConfigFunc = orig }()

		action := surfaceConfigFunc(dir)
		if action.status != "CREATED" {
			t.Errorf("expected CREATED, got %s: %s", action.status, action.detail)
		}
		if action.detail != "api" {
			t.Errorf("expected detail 'api', got %q", action.detail)
		}

		// Verify config file was updated
		cfg, err := forgeconfig.ReadConfig(dir)
		if err != nil {
			t.Fatalf("read config: %v", err)
		}
		if cfg.Surfaces["."] != "api" {
			t.Errorf("expected surfaces['.']=api, got %v", cfg.Surfaces)
		}
	})
}

// TestSurfaceConfigFuncOverride tests that the surfaceConfigFunc variable can be overridden.
func TestSurfaceConfigFuncOverride(t *testing.T) {
	t.Run("test override works for init integration", func(t *testing.T) {
		orig := surfaceConfigFunc
		surfaceConfigFunc = func(_ string) initAction {
			return initAction{status: "CREATED", target: "surfaces", detail: "test override"}
		}
		defer func() { surfaceConfigFunc = orig }()

		// Create config file to make the step work
		env := newInitTestEnv(t)
		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := env.stdout.String()
		if !strings.Contains(output, "surfaces") {
			t.Errorf("expected surfaces in output, got %q", output)
		}
	})
}

// TestSurfaceConfigWithTestOverride tests the full init flow with surface config.
func TestSurfaceConfigWithTestOverride(t *testing.T) {
	t.Run("surface step creates config with surfaces", func(t *testing.T) {
		origConfig := configInitFunc
		origSurface := surfaceConfigFunc

		configInitFunc = func(projectRoot string) initAction {
			configFile := filepath.Join(projectRoot, feature.ForgeDir, feature.ForgeConfigFileName)
			auto := autoConfigDefaults()
			cfg := forgeconfig.Config{Auto: auto}
			if err := writeConfigFile(configFile, &cfg); err != nil {
				return initAction{status: "FAILED", target: ".forge/config.yaml", detail: err.Error()}
			}
			return initAction{status: "CREATED", target: ".forge/config.yaml", detail: "test"}
		}

		surfaceConfigFunc = func(projectRoot string) initAction {
			configFile := filepath.Join(projectRoot, feature.ForgeDir, feature.ForgeConfigFileName)
			cfg, err := forgeconfig.ReadConfig(projectRoot)
			if err != nil || cfg == nil {
				return initAction{status: "FAILED", target: "surfaces", detail: "read config failed"}
			}
			cfg.Surfaces = forgeconfig.SurfacesMap{".": "api"}
			if err := writeConfigFile(configFile, cfg); err != nil {
				return initAction{status: "FAILED", target: "surfaces", detail: err.Error()}
			}
			return initAction{status: "CREATED", target: "surfaces", detail: "api"}
		}

		defer func() {
			configInitFunc = origConfig
			surfaceConfigFunc = origSurface
		}()

		env := &initTestEnv{dir: t.TempDir()}
		err := env.run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify surfaces written to config
		configFile := env.path(feature.ForgeDir, feature.ForgeConfigFileName)
		data, err := os.ReadFile(configFile)
		if err != nil {
			t.Fatalf("config not found: %v", err)
		}
		content := string(data)
		if !strings.Contains(content, "surfaces: api") {
			t.Errorf("expected 'surfaces: api' in config, got:\n%s", content)
		}
	})
}

// TestSortedPaths tests the sortedPaths helper.
func TestSortedPaths(t *testing.T) {
	surfaces := forgeconfig.SurfacesMap{
		"backend":  "api",
		"frontend": "web",
		"cli":      "cli",
	}
	paths := sortedPaths(surfaces)
	if len(paths) != 3 {
		t.Fatalf("expected 3 paths, got %d", len(paths))
	}
	if paths[0] != "backend" || paths[1] != "cli" || paths[2] != "frontend" {
		t.Errorf("expected sorted order [backend, cli, frontend], got %v", paths)
	}
}

// --- Test helpers ---

func writePackageJSONDeps(t *testing.T, dir string, deps map[string]string) {
	t.Helper()
	type pkg struct {
		Dependencies map[string]string `json:"dependencies,omitempty"`
	}
	p := pkg{Dependencies: deps}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), append(data, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
}

func writePackageJSONWorkspaces(t *testing.T, dir string, workspaces []string) {
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

func mkdirAllTest(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal(err)
	}
}
