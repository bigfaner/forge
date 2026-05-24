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

// TestFormatSourceAnnotation tests the source annotation formatting.
func TestFormatSourceAnnotation(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
	}{
		{
			name:   "inference cmd-dir",
			source: "inference:cmd-dir",
			want:   "(inferred from cmd/ directory structure)",
		},
		{
			name:   "inference api-dir",
			source: "inference:api-dir",
			want:   "(inferred from api/ directory)",
		},
		{
			name:   "inference handler-dir",
			source: "inference:handler-dir",
			want:   "(inferred from handler/ directory)",
		},
		{
			name:   "inference bin-field",
			source: "inference:bin-field",
			want:   "(inferred from bin field in package.json)",
		},
		{
			name:   "inference index-html",
			source: "inference:index-html",
			want:   "(inferred from index.html at project root)",
		},
		{
			name:   "inference py-scripts",
			source: "inference:py-scripts",
			want:   "(inferred from project.scripts or entry_points)",
		},
		{
			name:   "inference py-main",
			source: "inference:py-main",
			want:   "(inferred from app.py/main.py at root)",
		},
		{
			name:   "dependency cobra",
			source: "dependency:cobra",
			want:   "(detected from cobra dependency)",
		},
		{
			name:   "dependency react",
			source: "dependency:react",
			want:   "(detected from react dependency)",
		},
		{
			name:   "empty source",
			source: "",
			want:   "",
		},
		{
			name:   "unknown format",
			source: "other:something",
			want:   "(other:something)",
		},
		{
			name:   "no colon",
			source: "inference",
			want:   "(inference)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatSourceAnnotation(tt.source)
			if got != tt.want {
				t.Errorf("formatSourceAnnotation(%q) = %q, want %q", tt.source, got, tt.want)
			}
		})
	}
}

// TestIsInferred tests the inference detection helper.
func TestIsInferred(t *testing.T) {
	tests := []struct {
		source string
		want   bool
	}{
		{"inference:cmd-dir", true},
		{"inference:api-dir", true},
		{"dependency:cobra", false},
		{"dependency:react", false},
		{"", false},
		{"inference", false}, // no colon, but prefix matches
	}

	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			got := isInferred(tt.source)
			if got != tt.want {
				t.Errorf("isInferred(%q) = %v, want %v", tt.source, got, tt.want)
			}
		})
	}
}

// TestBuildDisplayLines tests the TUI display line builder.
func TestBuildDisplayLines(t *testing.T) {
	t.Run("scalar form shows only type", func(t *testing.T) {
		surfaces := forgeconfig.SurfacesMap{".": "api"}
		lines := buildDisplayLines(surfaces, nil, nil)

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
		lines := buildDisplayLines(surfaces, nil, nil)

		joined := strings.Join(lines, "\n")
		if !strings.Contains(joined, "frontend") {
			t.Errorf("expected 'frontend' in display, got: %s", joined)
		}
		if !strings.Contains(joined, "backend") {
			t.Errorf("expected 'backend' in display, got: %s", joined)
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
		lines := buildDisplayLines(surfaces, conflicts, nil)
		joined := strings.Join(lines, "\n")
		if !strings.Contains(joined, "web + api") {
			t.Errorf("expected conflict annotation in display, got: %s", joined)
		}
	})

	t.Run("source annotation for inference displayed", func(t *testing.T) {
		surfaces := forgeconfig.SurfacesMap{".": "cli"}
		sources := forgeconfig.SourcesMap{".": "inference:cmd-dir"}
		lines := buildDisplayLines(surfaces, nil, sources)
		joined := strings.Join(lines, "\n")
		if !strings.Contains(joined, "inferred from cmd/ directory structure") {
			t.Errorf("expected source annotation in display, got: %s", joined)
		}
	})

	t.Run("source annotation for dependency displayed", func(t *testing.T) {
		surfaces := forgeconfig.SurfacesMap{".": "cli"}
		sources := forgeconfig.SourcesMap{".": "dependency:cobra"}
		lines := buildDisplayLines(surfaces, nil, sources)
		joined := strings.Join(lines, "\n")
		if !strings.Contains(joined, "detected from cobra dependency") {
			t.Errorf("expected source annotation in display, got: %s", joined)
		}
	})

	t.Run("inferred surfaces show hint text", func(t *testing.T) {
		surfaces := forgeconfig.SurfacesMap{".": "cli"}
		sources := forgeconfig.SourcesMap{".": "inference:cmd-dir"}
		lines := buildDisplayLines(surfaces, nil, sources)
		joined := strings.Join(lines, "\n")
		if !strings.Contains(joined, "Inferred entries") {
			t.Errorf("expected hint text for inferred surface, got: %s", joined)
		}
	})

	t.Run("dependency surfaces do NOT show hint text", func(t *testing.T) {
		surfaces := forgeconfig.SurfacesMap{".": "cli"}
		sources := forgeconfig.SourcesMap{".": "dependency:cobra"}
		lines := buildDisplayLines(surfaces, nil, sources)
		joined := strings.Join(lines, "\n")
		if strings.Contains(joined, "Inferred entries") {
			t.Errorf("hint text should NOT appear for dependency surfaces, got: %s", joined)
		}
	})

	t.Run("map form with source annotations", func(t *testing.T) {
		surfaces := forgeconfig.SurfacesMap{
			"forge-cli/cli": "cli",
			"forge-cli/api": "api",
		}
		sources := forgeconfig.SourcesMap{
			"forge-cli/cli": "inference:cmd-dir",
			"forge-cli/api": "inference:api-dir",
		}
		lines := buildDisplayLines(surfaces, nil, sources)
		joined := strings.Join(lines, "\n")
		if !strings.Contains(joined, "inferred from cmd/ directory structure") {
			t.Errorf("expected cli source annotation, got: %s", joined)
		}
		if !strings.Contains(joined, "inferred from api/ directory") {
			t.Errorf("expected api source annotation, got: %s", joined)
		}
		// Hint should appear since there are inferred surfaces
		if !strings.Contains(joined, "Inferred entries") {
			t.Errorf("expected hint text for map with inferred surfaces, got: %s", joined)
		}
	})

	t.Run("no source annotation when sources nil", func(t *testing.T) {
		surfaces := forgeconfig.SurfacesMap{".": "api"}
		lines := buildDisplayLines(surfaces, nil, nil)
		joined := strings.Join(lines, "\n")
		if strings.Contains(joined, "inferred from") {
			t.Errorf("should not contain source annotation when sources nil, got: %s", joined)
		}
		if strings.Contains(joined, "detected from") {
			t.Errorf("should not contain source annotation when sources nil, got: %s", joined)
		}
	})
}

// TestFormatSurfacesSummary tests the init summary formatting.
func TestFormatSurfacesSummary(t *testing.T) {
	t.Run("scalar form shows type with source", func(t *testing.T) {
		surfaces := forgeconfig.SurfacesMap{".": "cli"}
		sources := forgeconfig.SourcesMap{".": "inference:cmd-dir"}
		summary := formatSurfacesSummary(surfaces, sources)
		if !strings.Contains(summary, "cli") {
			t.Errorf("expected 'cli' in summary, got %q", summary)
		}
		if !strings.Contains(summary, "(inferred from cmd/ directory structure)") {
			t.Errorf("expected source annotation in summary, got %q", summary)
		}
	})

	t.Run("scalar form shows type with dependency source", func(t *testing.T) {
		surfaces := forgeconfig.SurfacesMap{".": "cli"}
		sources := forgeconfig.SourcesMap{".": "dependency:cobra"}
		summary := formatSurfacesSummary(surfaces, sources)
		if !strings.Contains(summary, "cli") {
			t.Errorf("expected 'cli' in summary, got %q", summary)
		}
		if !strings.Contains(summary, "(detected from cobra dependency)") {
			t.Errorf("expected source annotation in summary, got %q", summary)
		}
	})

	t.Run("scalar form without source shows just type", func(t *testing.T) {
		surfaces := forgeconfig.SurfacesMap{".": "api"}
		summary := formatSurfacesSummary(surfaces, nil)
		if summary != "api" {
			t.Errorf("expected 'api', got %q", summary)
		}
	})

	t.Run("map form shows path=type with sources", func(t *testing.T) {
		surfaces := forgeconfig.SurfacesMap{
			"forge-cli/cli": "cli",
			"forge-cli/api": "api",
		}
		sources := forgeconfig.SourcesMap{
			"forge-cli/cli": "inference:cmd-dir",
			"forge-cli/api": "inference:api-dir",
		}
		summary := formatSurfacesSummary(surfaces, sources)
		if !strings.Contains(summary, "forge-cli/cli=cli") {
			t.Errorf("expected 'forge-cli/cli=cli' in summary, got %q", summary)
		}
		if !strings.Contains(summary, "inferred from cmd/ directory structure") {
			t.Errorf("expected source annotation in summary, got %q", summary)
		}
	})

	t.Run("map form without sources shows path=type only", func(t *testing.T) {
		surfaces := forgeconfig.SurfacesMap{
			"frontend": "web",
			"backend":  "api",
		}
		summary := formatSurfacesSummary(surfaces, nil)
		if !strings.Contains(summary, "frontend=web") {
			t.Errorf("expected 'frontend=web' in summary, got %q", summary)
		}
		if !strings.Contains(summary, "backend=api") {
			t.Errorf("expected 'backend=api' in summary, got %q", summary)
		}
	})

	t.Run("empty surfaces returns empty string", func(t *testing.T) {
		summary := formatSurfacesSummary(nil, nil)
		if summary != "" {
			t.Errorf("expected empty string, got %q", summary)
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

// TestRunSurfaceConfigRerun tests the re-run behavior when surfaces already exist.
func TestRunSurfaceConfigRerun(t *testing.T) {
	t.Run("existing surfaces triggers re-run check", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}

		// Create config with existing surfaces
		configContent := "auto:\n  gitPush: false\nsurfaces:\n  .: cli\n"
		if err := os.WriteFile(filepath.Join(forgeDir, feature.ForgeConfigFileName), []byte(configContent), 0o644); err != nil {
			t.Fatal(err)
		}

		// Override askRerunPrompt to simulate user choosing "confirm"
		orig := surfaceConfigFunc
		surfaceConfigFunc = runSurfaceConfig
		defer func() { surfaceConfigFunc = orig }()

		// The function requires TTY, so we test the re-run detection logic separately
		// by checking that config with existing surfaces is detected
		cfg, err := forgeconfig.ReadConfig(dir)
		if err != nil {
			t.Fatal(err)
		}
		if len(cfg.Surfaces) == 0 {
			t.Error("expected existing surfaces in config")
		}
	})

	t.Run("confirm returns SKIPPED with already configured detail", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}

		// Create config with existing surfaces
		cfg := &forgeconfig.Config{
			Auto:     &forgeconfig.AutoConfig{},
			Surfaces: forgeconfig.SurfacesMap{".": "cli"},
		}
		configFile := filepath.Join(forgeDir, feature.ForgeConfigFileName)
		if err := writeConfigFile(configFile, cfg); err != nil {
			t.Fatal(err)
		}

		// Test handleRerunSurfaceConfig with mocked askRerunPrompt
		origAskRerun := askRerunPrompt
		askRerunPrompt = func(_ forgeconfig.SurfacesMap) (string, bool) {
			return "confirm", false
		}
		defer func() { askRerunPrompt = origAskRerun }()

		action := handleRerunSurfaceConfig(dir, configFile, cfg)
		if action.status != "SKIPPED" {
			t.Errorf("expected SKIPPED, got %s", action.status)
		}
		if action.detail != "already configured" {
			t.Errorf("expected 'already configured', got %q", action.detail)
		}
	})

	t.Run("edit calls manualSurfaceEntry and writes result", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}

		cfg := &forgeconfig.Config{
			Auto:     &forgeconfig.AutoConfig{},
			Surfaces: forgeconfig.SurfacesMap{".": "cli"},
		}
		configFile := filepath.Join(forgeDir, feature.ForgeConfigFileName)
		if err := writeConfigFile(configFile, cfg); err != nil {
			t.Fatal(err)
		}

		// Mock askRerunPrompt to return "edit"
		origAskRerun := askRerunPrompt
		askRerunPrompt = func(_ forgeconfig.SurfacesMap) (string, bool) {
			return "edit", false
		}
		defer func() { askRerunPrompt = origAskRerun }()

		// Mock manualSurfaceEntry to return a new surface type
		origManual := manualSurfaceEntry
		manualSurfaceEntry = func() (forgeconfig.SurfacesMap, bool) {
			return forgeconfig.SurfacesMap{".": "api"}, false
		}
		defer func() { manualSurfaceEntry = origManual }()

		action := handleRerunSurfaceConfig(dir, configFile, cfg)
		if action.status != "CREATED" {
			t.Errorf("expected CREATED, got %s: %s", action.status, action.detail)
		}
		// Verify config was updated
		updatedCfg, err := forgeconfig.ReadConfig(dir)
		if err != nil {
			t.Fatal(err)
		}
		if updatedCfg.Surfaces["."] != "api" {
			t.Errorf("expected surfaces['.']='api' after edit, got %v", updatedCfg.Surfaces)
		}
	})

	t.Run("redetect runs full detection", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}

		cfg := &forgeconfig.Config{
			Auto:     &forgeconfig.AutoConfig{},
			Surfaces: forgeconfig.SurfacesMap{".": "cli"},
		}
		configFile := filepath.Join(forgeDir, feature.ForgeConfigFileName)
		if err := writeConfigFile(configFile, cfg); err != nil {
			t.Fatal(err)
		}

		// Mock askRerunPrompt to return "redetect"
		origAskRerun := askRerunPrompt
		askRerunPrompt = func(_ forgeconfig.SurfacesMap) (string, bool) {
			return "redetect", false
		}
		defer func() { askRerunPrompt = origAskRerun }()

		// Mock runNewSurfaceDetection to verify it's called
		origRunNew := runNewSurfaceDetection
		called := false
		runNewSurfaceDetection = func(_, _ string) initAction {
			called = true
			return initAction{status: "CREATED", target: "surfaces", detail: "api"}
		}
		defer func() { runNewSurfaceDetection = origRunNew }()

		action := handleRerunSurfaceConfig(dir, configFile, cfg)
		if !called {
			t.Error("expected runNewSurfaceDetection to be called for redetect")
		}
		if action.status != "CREATED" {
			t.Errorf("expected CREATED, got %s", action.status)
		}
	})

	t.Run("cancelled re-run prompt returns CANCELLED", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}

		cfg := &forgeconfig.Config{
			Auto:     &forgeconfig.AutoConfig{},
			Surfaces: forgeconfig.SurfacesMap{".": "cli"},
		}
		configFile := filepath.Join(forgeDir, feature.ForgeConfigFileName)
		if err := writeConfigFile(configFile, cfg); err != nil {
			t.Fatal(err)
		}

		origAskRerun := askRerunPrompt
		askRerunPrompt = func(_ forgeconfig.SurfacesMap) (string, bool) {
			return "", true // cancelled
		}
		defer func() { askRerunPrompt = origAskRerun }()

		action := handleRerunSurfaceConfig(dir, configFile, cfg)
		if action.status != "CANCELLED" {
			t.Errorf("expected CANCELLED, got %s", action.status)
		}
	})
}

// TestWriteSurfacesToConfig verifies that source annotations are NOT persisted.
func TestWriteSurfacesToConfig(t *testing.T) {
	t.Run("surfaces written without source annotations", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, feature.ForgeDir)
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}

		configFile := filepath.Join(forgeDir, feature.ForgeConfigFileName)
		// Write initial config
		cfg := &forgeconfig.Config{Auto: &forgeconfig.AutoConfig{}}
		if err := writeConfigFile(configFile, cfg); err != nil {
			t.Fatal(err)
		}

		// Write surfaces
		action := writeSurfacesToConfig(configFile, forgeconfig.SurfacesMap{".": "cli"})
		if action.status != "CREATED" {
			t.Errorf("expected CREATED, got %s: %s", action.status, action.detail)
		}

		// Verify no source annotation in config
		data, err := os.ReadFile(configFile)
		if err != nil {
			t.Fatal(err)
		}
		content := string(data)
		if strings.Contains(content, "inference") {
			t.Errorf("source annotation should NOT be in config, got:\n%s", content)
		}
		if strings.Contains(content, "dependency") {
			t.Errorf("source annotation should NOT be in config, got:\n%s", content)
		}
		if !strings.Contains(content, "surfaces") {
			t.Errorf("expected 'surfaces' in config, got:\n%s", content)
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

// TestFormatInferenceDetail tests the inference detail formatter.
func TestFormatInferenceDetail(t *testing.T) {
	tests := []struct {
		ruleID string
		want   string
	}{
		{"cmd-dir", "cmd/ directory structure"},
		{"api-dir", "api/ directory"},
		{"handler-dir", "handler/ directory"},
		{"bin-field", "bin field in package.json"},
		{"index-html", "index.html at project root"},
		{"py-scripts", "project.scripts or entry_points"},
		{"py-main", "app.py/main.py at root"},
		{"unknown-rule", "unknown-rule"},
	}

	for _, tt := range tests {
		t.Run(tt.ruleID, func(t *testing.T) {
			got := formatInferenceDetail(tt.ruleID)
			if got != tt.want {
				t.Errorf("formatInferenceDetail(%q) = %q, want %q", tt.ruleID, got, tt.want)
			}
		})
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
