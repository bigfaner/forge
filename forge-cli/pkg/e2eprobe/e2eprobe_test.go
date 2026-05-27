package e2eprobe

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProbeEndpoint(t *testing.T) {
	t.Run("unreachable returns false", func(t *testing.T) {
		if ProbeEndpoint("http://127.0.0.1:1", 100*time.Millisecond) {
			t.Error("expected false for unreachable endpoint")
		}
	})

	t.Run("success returns true", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		if !ProbeEndpoint(ts.URL, 2*time.Second) {
			t.Error("expected true for healthy endpoint")
		}
	})

	t.Run("server error returns false", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		if ProbeEndpoint(ts.URL, 2*time.Second) {
			t.Error("expected false for 500 status")
		}
	})
}

func TestProbeServers(t *testing.T) {
	t.Run("no config.yaml returns true", func(t *testing.T) {
		dir := t.TempDir()
		if !ProbeServers(dir, "") {
			t.Error("expected true when no config.yaml exists")
		}
	})

	t.Run("empty config returns true", func(t *testing.T) {
		dir := t.TempDir()
		_ = os.MkdirAll(filepath.Join(dir, "tests"), 0755)
		_ = os.WriteFile(filepath.Join(dir, "tests", "config.yaml"), []byte(""), 0644)
		if !ProbeServers(dir, "") {
			t.Error("expected true for empty config")
		}
	})

	t.Run("unreachable endpoint returns false", func(t *testing.T) {
		dir := t.TempDir()
		_ = os.MkdirAll(filepath.Join(dir, "tests"), 0755)
		_ = os.WriteFile(filepath.Join(dir, "tests", "config.yaml"), []byte("baseUrl: http://127.0.0.1:1\n"), 0644)
		if ProbeServers(dir, "") {
			t.Error("expected false for unreachable endpoint")
		}
	})

	t.Run("reachable baseUrl returns true", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		dir := t.TempDir()
		_ = os.MkdirAll(filepath.Join(dir, "tests"), 0755)
		config := "baseUrl: " + ts.URL + "\n"
		_ = os.WriteFile(filepath.Join(dir, "tests", "config.yaml"), []byte(config), 0644)

		if !ProbeServers(dir, "") {
			t.Error("expected true for reachable baseUrl")
		}
	})

	t.Run("reachable apiBaseUrl returns true", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		dir := t.TempDir()
		_ = os.MkdirAll(filepath.Join(dir, "tests"), 0755)
		config := "apiBaseUrl: " + ts.URL + "\n"
		_ = os.WriteFile(filepath.Join(dir, "tests", "config.yaml"), []byte(config), 0644)

		if !ProbeServers(dir, "") {
			t.Error("expected true for reachable apiBaseUrl")
		}
	})

	t.Run("both baseUrl and apiBaseUrl reachable returns true", func(t *testing.T) {
		ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts1.Close()

		ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts2.Close()

		dir := t.TempDir()
		_ = os.MkdirAll(filepath.Join(dir, "tests"), 0755)
		config := "baseUrl: " + ts1.URL + "\napiBaseUrl: " + ts2.URL + "\n"
		_ = os.WriteFile(filepath.Join(dir, "tests", "config.yaml"), []byte(config), 0644)

		if !ProbeServers(dir, "") {
			t.Error("expected true when both endpoints are reachable")
		}
	})

	t.Run("custom path is appended to URL", func(t *testing.T) {
		var requestedPath string
		ts := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			requestedPath = r.URL.Path
		}))
		defer ts.Close()

		dir := t.TempDir()
		_ = os.MkdirAll(filepath.Join(dir, "tests"), 0755)
		config := "baseUrl: " + ts.URL + "\n"
		_ = os.WriteFile(filepath.Join(dir, "tests", "config.yaml"), []byte(config), 0644)

		if !ProbeServers(dir, "/api/health") {
			t.Error("expected true for reachable baseUrl with custom path")
		}
		if requestedPath != "/api/health" {
			t.Errorf("requested path = %q, want %q", requestedPath, "/api/health")
		}
	})

	t.Run("default path is /health", func(t *testing.T) {
		var requestedPath string
		ts := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			requestedPath = r.URL.Path
		}))
		defer ts.Close()

		dir := t.TempDir()
		_ = os.MkdirAll(filepath.Join(dir, "tests"), 0755)
		config := "baseUrl: " + ts.URL + "\n"
		_ = os.WriteFile(filepath.Join(dir, "tests", "config.yaml"), []byte(config), 0644)

		if !ProbeServers(dir, "") {
			t.Error("expected true for reachable baseUrl with default path")
		}
		if requestedPath != "/health" {
			t.Errorf("requested path = %q, want %q", requestedPath, "/health")
		}
	})
}

func TestExtractYAMLStringField(t *testing.T) {
	tests := []struct {
		name  string
		data  string
		field string
		want  string
	}{
		{"found", "baseUrl: http://localhost:3000\napiBaseUrl: http://localhost:8080\n", "baseUrl", "http://localhost:3000"},
		{"quoted", "baseUrl: 'http://localhost:3000'\n", "baseUrl", "http://localhost:3000"},
		{"double quoted", "baseUrl: \"http://localhost:3000\"\n", "baseUrl", "http://localhost:3000"},
		{"missing field", "other: value\n", "baseUrl", ""},
		{"empty data", "", "baseUrl", ""},
		{"field with spaces", "apiBaseUrl:   http://api:8080  \n", "apiBaseUrl", "http://api:8080"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ExtractYAMLStringField([]byte(tc.data), tc.field)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
