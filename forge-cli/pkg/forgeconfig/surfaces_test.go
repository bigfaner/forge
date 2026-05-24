package forgeconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSurfacesMap_UnmarshalYAML(t *testing.T) {
	t.Run("scalar form converts to dot key map", func(t *testing.T) {
		dir := setupConfig(t, "surfaces: api\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cfg.Surfaces) != 1 {
			t.Fatalf("expected 1 surface, got %d", len(cfg.Surfaces))
		}
		if cfg.Surfaces["."] != "api" {
			t.Errorf("expected Surfaces['.']='api', got %q", cfg.Surfaces["."])
		}
	})

	t.Run("map form used as-is", func(t *testing.T) {
		dir := setupConfig(t, "surfaces:\n  frontend: web\n  backend: api\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cfg.Surfaces) != 2 {
			t.Fatalf("expected 2 surfaces, got %d", len(cfg.Surfaces))
		}
		if cfg.Surfaces["frontend"] != "web" {
			t.Errorf("expected Surfaces['frontend']='web', got %q", cfg.Surfaces["frontend"])
		}
		if cfg.Surfaces["backend"] != "api" {
			t.Errorf("expected Surfaces['backend']='api', got %q", cfg.Surfaces["backend"])
		}
	})

	t.Run("empty config gives nil surfaces", func(t *testing.T) {
		dir := setupConfig(t, "{}\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Surfaces != nil {
			t.Errorf("expected Surfaces nil for empty config, got %v", cfg.Surfaces)
		}
	})

	t.Run("absent field is nil", func(t *testing.T) {
		dir := setupConfig(t, "test-framework: pytest\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Surfaces != nil {
			t.Errorf("expected Surfaces nil when absent, got %v", cfg.Surfaces)
		}
	})
}

func TestSurfacesMap_MarshalYAML(t *testing.T) {
	t.Run("single dot key serializes as scalar", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			Surfaces: SurfacesMap{".": "api"},
		}
		if err := writeConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := os.ReadFile(filepath.Join(dir, ".forge", "config.yaml"))
		if err != nil {
			t.Fatalf("read file: %v", err)
		}
		content := string(data)
		if !strings.Contains(content, "surfaces: api") {
			t.Errorf("expected scalar form 'surfaces: api', got:\n%s", content)
		}
	})

	t.Run("map form serializes as map", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			Surfaces: SurfacesMap{"frontend": "web", "backend": "api"},
		}
		if err := writeConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		readback, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if readback.Surfaces["frontend"] != "web" {
			t.Errorf("expected frontend=web, got %q", readback.Surfaces["frontend"])
		}
		if readback.Surfaces["backend"] != "api" {
			t.Errorf("expected backend=api, got %q", readback.Surfaces["backend"])
		}
	})

	t.Run("nil surfaces serializes as empty map", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			Surfaces: nil,
		}
		if err := writeConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := os.ReadFile(filepath.Join(dir, ".forge", "config.yaml"))
		if err != nil {
			t.Fatalf("read file: %v", err)
		}
		content := string(data)
		if !strings.Contains(content, "surfaces: {}") {
			t.Errorf("expected 'surfaces: {}' for nil, got:\n%s", content)
		}
	})

	t.Run("round trip preserves scalar form", func(t *testing.T) {
		dir := setupConfig(t, "surfaces: api\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := writeConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cfg2, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cfg2.Surfaces) != 1 || cfg2.Surfaces["."] != "api" {
			t.Errorf("round trip failed: got %v", cfg2.Surfaces)
		}
	})
}

func TestSurfaceTypes(t *testing.T) {
	t.Run("nil returns nil", func(t *testing.T) {
		result := SurfaceTypes(nil)
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})

	t.Run("empty map returns nil", func(t *testing.T) {
		result := SurfaceTypes(map[string]string{})
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})

	t.Run("deduplicates values", func(t *testing.T) {
		surfaces := map[string]string{
			"frontend": "web",
			"backend":  "api",
			"shared":   "web",
		}
		result := SurfaceTypes(surfaces)
		if len(result) != 2 {
			t.Fatalf("expected 2 types, got %d: %v", len(result), result)
		}
		seen := map[string]bool{}
		for _, typ := range result {
			seen[typ] = true
		}
		if !seen["web"] || !seen["api"] {
			t.Errorf("expected web and api, got %v", result)
		}
	})

	t.Run("single type returns single value", func(t *testing.T) {
		surfaces := map[string]string{".": "api"}
		result := SurfaceTypes(surfaces)
		if len(result) != 1 || result[0] != "api" {
			t.Errorf("expected [api], got %v", result)
		}
	})
}

func TestReadSurfaces(t *testing.T) {
	t.Run("reads scalar surfaces", func(t *testing.T) {
		dir := setupConfig(t, "surfaces: api\n")
		surfaces, err := ReadSurfaces(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(surfaces) != 1 || surfaces["."] != "api" {
			t.Errorf("expected {'.':'api'}, got %v", surfaces)
		}
	})

	t.Run("reads map surfaces", func(t *testing.T) {
		dir := setupConfig(t, "surfaces:\n  frontend: web\n")
		surfaces, err := ReadSurfaces(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if surfaces["frontend"] != "web" {
			t.Errorf("expected frontend=web, got %v", surfaces)
		}
	})

	t.Run("missing config returns nil", func(t *testing.T) {
		dir := t.TempDir()
		surfaces, err := ReadSurfaces(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if surfaces != nil {
			t.Errorf("expected nil, got %v", surfaces)
		}
	})
}

func TestMigrateInterfacesToSurfaces(t *testing.T) {
	t.Run("no interfaces field is no-op", func(t *testing.T) {
		dir := setupConfig(t, "surfaces: api\n")
		err := MigrateInterfacesToSurfaces(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("single interface auto-migrates to scalar", func(t *testing.T) {
		dir := setupConfig(t, "interfaces:\n  - api\n")
		err := MigrateInterfacesToSurfaces(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cfg.Surfaces) != 1 || cfg.Surfaces["."] != "api" {
			t.Errorf("expected migrated surfaces {'.':'api'}, got %v", cfg.Surfaces)
		}
	})

	t.Run("multi-interface returns error", func(t *testing.T) {
		dir := setupConfig(t, "interfaces:\n  - web\n  - api\n")
		err := MigrateInterfacesToSurfaces(dir)
		if err == nil {
			t.Fatal("expected error for multi-interface migration")
		}
		multiErr, ok := err.(*ErrMultiInterfaceMigration)
		if !ok {
			t.Fatalf("expected ErrMultiInterfaceMigration, got %T: %v", err, err)
		}
		if len(multiErr.Interfaces) != 2 {
			t.Errorf("expected 2 interfaces in error, got %d", len(multiErr.Interfaces))
		}
	})

	t.Run("surfaces already configured skips migration", func(t *testing.T) {
		dir := setupConfig(t, "interfaces:\n  - api\nsurfaces:\n  backend: api\n")
		err := MigrateInterfacesToSurfaces(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Surfaces["."] != "" {
			t.Errorf("expected no dot key (original map preserved), got %v", cfg.Surfaces)
		}
		if cfg.Surfaces["backend"] != "api" {
			t.Errorf("expected backend=api, got %v", cfg.Surfaces)
		}
	})

	t.Run("missing config is no-op", func(t *testing.T) {
		dir := t.TempDir()
		err := MigrateInterfacesToSurfaces(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
