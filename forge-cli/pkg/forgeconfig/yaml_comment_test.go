package forgeconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteConfigWithSources_ScalarForm(t *testing.T) {
	t.Run("appends source comment to scalar surfaces value", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			Surfaces: SurfacesMap{".": "cli"},
		}
		sources := SourcesMap{".": "inference:cmd-dir"}

		if err := WriteConfigWithSources(dir, cfg, sources); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := os.ReadFile(filepath.Join(dir, ".forge", "config.yaml"))
		if err != nil {
			t.Fatalf("unexpected error reading config: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "# source: inference:cmd-dir") {
			t.Errorf("expected source comment in config, got:\n%s", content)
		}
	})
}

func TestWriteConfigWithSources_MapForm(t *testing.T) {
	t.Run("appends source comment to map-form surfaces value", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			Surfaces: SurfacesMap{"frontend": "web", "backend": "api"},
		}
		sources := SourcesMap{"frontend": "inference:index-html", "backend": "inference:api-dir"}

		if err := WriteConfigWithSources(dir, cfg, sources); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := os.ReadFile(filepath.Join(dir, ".forge", "config.yaml"))
		if err != nil {
			t.Fatalf("unexpected error reading config: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "# source: inference:index-html") {
			t.Errorf("expected frontend source comment in config, got:\n%s", content)
		}
		if !strings.Contains(content, "# source: inference:api-dir") {
			t.Errorf("expected backend source comment in config, got:\n%s", content)
		}
	})
}

func TestWriteConfigWithSources_NoSources(t *testing.T) {
	t.Run("no sources produces no comment", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			Surfaces: SurfacesMap{".": "cli"},
		}

		if err := WriteConfigWithSources(dir, cfg, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := os.ReadFile(filepath.Join(dir, ".forge", "config.yaml"))
		if err != nil {
			t.Fatalf("unexpected error reading config: %v", err)
		}

		content := string(data)
		if strings.Contains(content, "# source:") {
			t.Errorf("expected no source comment when sources is nil, got:\n%s", content)
		}
	})

	t.Run("empty sources map produces no comment", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			Surfaces: SurfacesMap{".": "cli"},
		}
		sources := SourcesMap{}

		if err := WriteConfigWithSources(dir, cfg, sources); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := os.ReadFile(filepath.Join(dir, ".forge", "config.yaml"))
		if err != nil {
			t.Fatalf("unexpected error reading config: %v", err)
		}

		content := string(data)
		if strings.Contains(content, "# source:") {
			t.Errorf("expected no source comment when sources is empty, got:\n%s", content)
		}
	})
}

func TestWriteConfigWithSources_YAMLStillParses(t *testing.T) {
	t.Run("config with comment round-trips correctly", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			Surfaces: SurfacesMap{".": "cli"},
		}
		sources := SourcesMap{".": "inference:cmd-dir"}

		if err := WriteConfigWithSources(dir, cfg, sources); err != nil {
			t.Fatalf("unexpected error writing: %v", err)
		}

		readback, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error reading: %v", err)
		}
		if readback.Surfaces == nil {
			t.Fatal("expected Surfaces non-nil")
		}
		if v, ok := readback.Surfaces["."]; !ok || v != "cli" {
			t.Errorf("expected surfaces '.': 'cli', got %v", readback.Surfaces)
		}
	})
}

func TestWriteConfigWithSources_CommentStrippedStillWorks(t *testing.T) {
	t.Run("removing comment does not affect config parsing", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			Surfaces: SurfacesMap{".": "api"},
		}
		sources := SourcesMap{".": "inference:api-dir"}

		if err := WriteConfigWithSources(dir, cfg, sources); err != nil {
			t.Fatalf("unexpected error writing: %v", err)
		}

		// Read the file and strip the comment (simulating external editor)
		configPath := filepath.Join(dir, ".forge", "config.yaml")
		data, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("unexpected error reading config: %v", err)
		}
		stripped := strings.ReplaceAll(string(data), "# source: inference:api-dir", "")
		if err := os.WriteFile(configPath, []byte(stripped), 0o644); err != nil {
			t.Fatalf("unexpected error writing stripped config: %v", err)
		}

		// Config should still parse correctly
		readback, err := ReadConfig(dir)
		if err != nil {
			t.Fatalf("unexpected error reading stripped config: %v", err)
		}
		if v, ok := readback.Surfaces["."]; !ok || v != "api" {
			t.Errorf("expected surfaces '.': 'api' after comment strip, got %v", readback.Surfaces)
		}
	})
}

func TestReadSurfaceComment(t *testing.T) {
	t.Run("extracts comment from surfaces node", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			Surfaces: SurfacesMap{".": "cli"},
		}
		sources := SourcesMap{".": "inference:cmd-dir"}

		if err := WriteConfigWithSources(dir, cfg, sources); err != nil {
			t.Fatalf("unexpected error writing: %v", err)
		}

		comment, err := ReadSurfaceComment(dir)
		if err != nil {
			t.Fatalf("unexpected error reading comment: %v", err)
		}
		if comment != "source: inference:cmd-dir" {
			t.Errorf("expected 'source: inference:cmd-dir', got %q", comment)
		}
	})

	t.Run("returns empty when no comment present", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			Surfaces: SurfacesMap{".": "cli"},
		}

		if err := writeConfig(dir, cfg); err != nil {
			t.Fatalf("unexpected error writing: %v", err)
		}

		comment, err := ReadSurfaceComment(dir)
		if err != nil {
			t.Fatalf("unexpected error reading comment: %v", err)
		}
		if comment != "" {
			t.Errorf("expected empty comment when none present, got %q", comment)
		}
	})

	t.Run("returns empty when no config file exists", func(t *testing.T) {
		dir := t.TempDir()

		comment, err := ReadSurfaceComment(dir)
		if err != nil {
			t.Fatalf("unexpected error reading comment: %v", err)
		}
		if comment != "" {
			t.Errorf("expected empty comment for missing config, got %q", comment)
		}
	})
}

func TestWriteConfigWithSources_DependencySource(t *testing.T) {
	t.Run("appends dependency source comment", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &Config{
			Surfaces: SurfacesMap{".": "cli"},
		}
		sources := SourcesMap{".": "dependency:cobra"}

		if err := WriteConfigWithSources(dir, cfg, sources); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := os.ReadFile(filepath.Join(dir, ".forge", "config.yaml"))
		if err != nil {
			t.Fatalf("unexpected error reading config: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "# source: dependency:cobra") {
			t.Errorf("expected dependency source comment in config, got:\n%s", content)
		}
	})
}
