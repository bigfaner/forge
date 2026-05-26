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

	t.Run("bug: uppercase surface type normalized to lowercase (scalar)", func(t *testing.T) {
		dir := setupConfig(t, "surfaces: CLI\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Surfaces["."] != "cli" {
			t.Errorf("expected 'cli' (lowercased), got %q", cfg.Surfaces["."])
		}
		warnings := ValidateSurfaceTypes(cfg.Surfaces)
		if len(warnings) != 0 {
			t.Errorf("expected no warnings after normalization, got %v", warnings)
		}
	})

	t.Run("bug: mixed-case surface type normalized to lowercase (map)", func(t *testing.T) {
		dir := setupConfig(t, "surfaces:\n  frontend: Web\n  backend: API\n")
		cfg, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Surfaces["frontend"] != "web" {
			t.Errorf("expected 'web' (lowercased), got %q", cfg.Surfaces["frontend"])
		}
		if cfg.Surfaces["backend"] != "api" {
			t.Errorf("expected 'api' (lowercased), got %q", cfg.Surfaces["backend"])
		}
		warnings := ValidateSurfaceTypes(cfg.Surfaces)
		if len(warnings) != 0 {
			t.Errorf("expected no warnings after normalization, got %v", warnings)
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

	t.Run("unknown types are filtered out", func(t *testing.T) {
		surfaces := map[string]string{
			"frontend": "web",
			"backend":  "api",
			"legacy":   "web-ui",
			"mobile":   "mobile-ui",
			"custom":   "unknown-type",
		}
		warnings := ValidateSurfaceTypes(surfaces)
		result := SurfaceTypes(surfaces)

		// Only known types should appear
		if len(result) != 2 {
			t.Fatalf("expected 2 known types, got %d: %v", len(result), result)
		}
		seen := map[string]bool{}
		for _, typ := range result {
			seen[typ] = true
		}
		if !seen["web"] || !seen["api"] {
			t.Errorf("expected web and api, got %v", result)
		}

		// Unknown types should produce warnings
		if len(warnings) != 3 {
			t.Fatalf("expected 3 warnings for unknown types, got %d: %v", len(warnings), warnings)
		}
		for _, w := range warnings {
			if !strings.Contains(w, "unknown surface type") {
				t.Errorf("warning should mention 'unknown surface type', got: %s", w)
			}
		}
	})

	t.Run("all known types pass without warnings", func(t *testing.T) {
		surfaces := map[string]string{
			"a": "web",
			"b": "mobile",
			"c": "api",
			"d": "cli",
			"e": "tui",
		}
		warnings := ValidateSurfaceTypes(surfaces)
		result := SurfaceTypes(surfaces)

		if len(warnings) != 0 {
			t.Errorf("expected 0 warnings for all known types, got %d: %v", len(warnings), warnings)
		}
		if len(result) != 5 {
			t.Errorf("expected 5 types, got %d: %v", len(result), result)
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

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		// AC1: strip leading ./, trailing /, convert \ to /
		{name: "strip leading ./", input: "./frontend", want: "frontend"},
		{name: "strip trailing /", input: "frontend/", want: "frontend"},
		{name: "strip both ./ and trailing /", input: "./frontend/", want: "frontend"},
		{name: "convert backslash to forward slash", input: `frontend\src`, want: "frontend/src"},
		{name: "mixed backslashes and dots", input: `.\frontend\src\`, want: "frontend/src"},
		{name: "already clean path", input: "frontend/api/routes", want: "frontend/api/routes"},

		// AC2: paths containing .. return error
		{name: "double dot returns error", input: "../etc", wantErr: true},
		{name: "embedded double dot returns error", input: "foo/../bar", wantErr: true},
		{name: "trailing double dot returns error", input: "foo/..", wantErr: true},
		{name: "double dot with backslash returns error", input: `foo\..`, wantErr: true},

		// Edge cases
		{name: "empty string returns empty", input: "", want: ""},
		{name: "single dot is preserved", input: ".", want: "."},
		{name: "multiple slashes collapsed path", input: "frontend//api", want: "frontend//api"},
		{name: "only dots in segments are rejected", input: "a/.../b", want: "a/.../b"},
		{name: "double dot exact segment match", input: "a/..b", want: "a/..b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizePath(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NormalizePath(%q) expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("NormalizePath(%q) unexpected error: %v", tt.input, err)
				return
			}
			if got != tt.want {
				t.Errorf("NormalizePath(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMatchSurface(t *testing.T) {
	// AC4: scalar form -- any path returns value directly
	t.Run("scalar form returns value for any path", func(t *testing.T) {
		surfaces := map[string]string{".": "api"}
		got, err := MatchSurface(surfaces, "src/main.go")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Type != "api" {
			t.Errorf("expected Type='api', got %q", got.Type)
		}
		if got.Key != "." {
			t.Errorf("expected Key='.', got %q", got.Key)
		}
	})

	t.Run("scalar form returns value for empty path", func(t *testing.T) {
		surfaces := map[string]string{".": "web"}
		got, err := MatchSurface(surfaces, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Type != "web" {
			t.Errorf("expected Type='web', got %q", got.Type)
		}
		if got.Key != "." {
			t.Errorf("expected Key='.', got %q", got.Key)
		}
	})

	// AC5: map form -- segment prefix matching, longest wins
	t.Run("segment prefix matching longest wins", func(t *testing.T) {
		surfaces := map[string]string{
			"frontend":     "web",
			"frontend/api": "api",
		}
		got, err := MatchSurface(surfaces, "frontend/api/routes")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Type != "api" {
			t.Errorf("expected Type='api' (2 segments beat 1), got %q", got.Type)
		}
		if got.Key != "frontend/api" {
			t.Errorf("expected Key='frontend/api', got %q", got.Key)
		}
	})

	t.Run("exact match returns value and key", func(t *testing.T) {
		surfaces := map[string]string{
			"frontend": "web",
			"backend":  "api",
		}
		got, err := MatchSurface(surfaces, "frontend")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Type != "web" {
			t.Errorf("expected Type='web', got %q", got.Type)
		}
		if got.Key != "frontend" {
			t.Errorf("expected Key='frontend', got %q", got.Key)
		}
	})

	t.Run("partial segment match returns value and key", func(t *testing.T) {
		surfaces := map[string]string{
			"frontend": "web",
		}
		got, err := MatchSurface(surfaces, "frontend/src/components")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Type != "web" {
			t.Errorf("expected Type='web', got %q", got.Type)
		}
		if got.Key != "frontend" {
			t.Errorf("expected Key='frontend', got %q", got.Key)
		}
	})

	// AC6: no partial match -- frontend-new does NOT match frontend
	t.Run("frontend-new does not match frontend", func(t *testing.T) {
		surfaces := map[string]string{
			"frontend": "web",
		}
		_, err := MatchSurface(surfaces, "frontend-new")
		if err == nil {
			t.Error("expected error for unmatched path 'frontend-new'")
		}
	})

	// AC7: unmatched path returns error with manual config hint
	t.Run("unmatched path returns error", func(t *testing.T) {
		surfaces := map[string]string{
			"frontend": "web",
		}
		_, err := MatchSurface(surfaces, "unknown-dir")
		if err == nil {
			t.Error("expected error for unmatched path")
		}
		// Verify error message contains config hint
		errMsg := err.Error()
		if !strings.Contains(errMsg, "forge init") {
			t.Errorf("error should mention 'forge init', got: %s", errMsg)
		}
	})

	// AC3: symlinks NOT resolved -- literal path matching only
	// This is verified by the string-based matching; no filesystem access

	t.Run("nil surfaces returns error", func(t *testing.T) {
		_, err := MatchSurface(nil, "frontend")
		if err == nil {
			t.Error("expected error for nil surfaces")
		}
	})

	t.Run("empty surfaces returns error", func(t *testing.T) {
		_, err := MatchSurface(map[string]string{}, "frontend")
		if err == nil {
			t.Error("expected error for empty surfaces")
		}
	})

	t.Run("path normalization applied before matching", func(t *testing.T) {
		surfaces := map[string]string{
			"frontend/api": "api",
		}
		got, err := MatchSurface(surfaces, "./frontend/api/routes/")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Type != "api" {
			t.Errorf("expected Type='api', got %q", got.Type)
		}
		if got.Key != "frontend/api" {
			t.Errorf("expected Key='frontend/api', got %q", got.Key)
		}
	})

	t.Run("backslash path matches forward slash config", func(t *testing.T) {
		surfaces := map[string]string{
			"frontend/api": "api",
		}
		got, err := MatchSurface(surfaces, `frontend\api\routes`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Type != "api" {
			t.Errorf("expected Type='api', got %q", got.Type)
		}
		if got.Key != "frontend/api" {
			t.Errorf("expected Key='frontend/api', got %q", got.Key)
		}
	})

	t.Run("dotdot path returns error in matching", func(t *testing.T) {
		surfaces := map[string]string{
			"frontend": "web",
		}
		_, err := MatchSurface(surfaces, "../etc/passwd")
		if err == nil {
			t.Error("expected error for path with '..'")
		}
	})

	t.Run("map form returns correct key for admin-panel", func(t *testing.T) {
		surfaces := map[string]string{
			"admin-panel":     "web",
			"payment-service": "api",
		}
		got, err := MatchSurface(surfaces, "admin-panel/src/components/App.tsx")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Key != "admin-panel" {
			t.Errorf("expected Key='admin-panel', got %q", got.Key)
		}
		if got.Type != "web" {
			t.Errorf("expected Type='web', got %q", got.Type)
		}
	})
}
