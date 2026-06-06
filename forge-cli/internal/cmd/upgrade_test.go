package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"forge-cli/pkg/types"
)

// --- Test helpers ---

// upgradeTestEnv provides a controlled environment for upgrade tests.
type upgradeTestEnv struct {
	stdout              bytes.Buffer
	stderr              bytes.Buffer
	origLookPath        func(string) (string, error)
	origFetch           func(string) ([]byte, error)
	origHTTPGet         func(url string) (*http.Response, error)
	origRunClaude       func([]string) error
	origPluginInstalled func() bool
	origBinDir          func() string
	origRename          func(string, string) error
	origVersion         string
}

func newUpgradeTestEnv(t *testing.T) *upgradeTestEnv {
	t.Helper()
	env := &upgradeTestEnv{
		origLookPath:        lookPathForUpgrade,
		origFetch:           fetchLatestRelease,
		origHTTPGet:         httpGet,
		origRunClaude:       runClaudeCommand,
		origPluginInstalled: pluginInstalledCheck,
		origBinDir:          forgeBinaryDir,
		origRename:          osRename,
		origVersion:         types.Version,
	}

	t.Cleanup(func() {
		lookPathForUpgrade = env.origLookPath
		fetchLatestRelease = env.origFetch
		httpGet = env.origHTTPGet
		runClaudeCommand = env.origRunClaude
		pluginInstalledCheck = env.origPluginInstalled
		forgeBinaryDir = env.origBinDir
		osRename = env.origRename
		types.Version = env.origVersion
	})

	return env
}

func (e *upgradeTestEnv) run() error {
	rootCmd.SetOut(&e.stdout)
	rootCmd.SetErr(&e.stderr)
	rootCmd.SetArgs([]string{"upgrade"})
	return rootCmd.Execute()
}

// --- Version comparison tests ---

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want int
	}{
		{"equal versions", "5.16.0", "5.16.0", 0},
		{"patch higher", "5.16.1", "5.16.0", 1},
		{"patch lower", "5.16.0", "5.16.1", -1},
		{"minor higher", "5.17.0", "5.16.0", 1},
		{"minor lower", "5.16.0", "5.17.0", -1},
		{"major higher", "6.0.0", "5.16.0", 1},
		{"major lower", "5.16.0", "6.0.0", -1},
		{"same major diff minor and patch", "5.15.9", "5.16.0", -1},
		{"large version numbers", "10.20.30", "9.99.99", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compareVersions(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("compareVersions(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestParseSemver(t *testing.T) {
	tests := []struct {
		input string
		want  [3]int
	}{
		{"5.16.0", [3]int{5, 16, 0}},
		{"1.2.3", [3]int{1, 2, 3}},
		{"0.0.0", [3]int{0, 0, 0}},
		{"10.20.30", [3]int{10, 20, 30}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseSemver(tt.input)
			if got != tt.want {
				t.Errorf("parseSemver(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// --- Tag parsing tests ---

func TestParseVersionFromTag(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		want    string
		wantErr bool
	}{
		{"forge-cli prefix", "forge-cli/v5.17.0", "5.17.0", false},
		{"forge-cli without v", "forge-cli/5.17.0", "5.17.0", false},
		{"v prefix only", "v5.17.0", "5.17.0", false},
		{"no prefix", "5.17.0", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseVersionFromTag(tt.tag)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseVersionFromTag(%q) expected error, got nil", tt.tag)
				}
				return
			}
			if err != nil {
				t.Errorf("parseVersionFromTag(%q) unexpected error: %v", tt.tag, err)
				return
			}
			if got != tt.want {
				t.Errorf("parseVersionFromTag(%q) = %q, want %q", tt.tag, got, tt.want)
			}
		})
	}
}

// --- Download URL construction tests ---

func TestBuildDownloadURL(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantOS  string
		wantExt string
	}{
		{
			name:    "non-windows",
			version: "5.17.0",
			wantOS:  runtime.GOOS,
			wantExt: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := buildDownloadURL(tt.version)
			wantSuffix := fmt.Sprintf("/forge-cli/v%s/forge-%s-%s-%s%s",
				tt.version, tt.version, runtime.GOOS, runtime.GOARCH, tt.wantExt)
			if !strings.Contains(url, wantSuffix) {
				t.Errorf("buildDownloadURL(%q) = %q, want to contain %q", tt.version, url, wantSuffix)
			}
			if !strings.HasPrefix(url, "https://github.com/bigfaner/forge/releases/download/") {
				t.Errorf("URL should start with github base, got %q", url)
			}
		})
	}
}

func TestBuildDownloadURL_Windows(t *testing.T) {
	// Test that Windows gets .exe suffix by testing the URL pattern
	if runtime.GOOS == "windows" {
		url := buildDownloadURL("5.17.0")
		if !strings.HasSuffix(url, ".exe") {
			t.Errorf("Windows download URL should end with .exe, got %q", url)
		}
	} else {
		url := buildDownloadURL("5.17.0")
		if strings.HasSuffix(url, ".exe") {
			t.Errorf("Non-Windows download URL should not end with .exe, got %q", url)
		}
	}
}

// --- CLI binary upgrade tests ---

func TestUpgradeCLIBinary_DevVersion(t *testing.T) {
	env := newUpgradeTestEnv(t)
	types.Version = "dev"

	action := upgradeCLIBinary(&env.stdout)

	if action.status != "SKIPPED" {
		t.Errorf("expected SKIPPED for dev version, got %q", action.status)
	}
	if !strings.Contains(action.detail, "development build") {
		t.Errorf("expected 'development build' in detail, got %q", action.detail)
	}
}

func TestUpgradeCLIBinary_AlreadyUpToDate(t *testing.T) {
	env := newUpgradeTestEnv(t)
	types.Version = "5.17.0"

	fetchLatestRelease = func(_ string) ([]byte, error) {
		release := githubRelease{TagName: "forge-cli/v5.17.0"}
		return json.Marshal(release)
	}

	action := upgradeCLIBinary(&env.stdout)

	if action.status != "SKIPPED" {
		t.Errorf("expected SKIPPED for same version, got %q", action.status)
	}
	if !strings.Contains(action.detail, "up to date") {
		t.Errorf("expected 'up to date' in detail, got %q", action.detail)
	}
}

func TestUpgradeCLIBinary_NewerAvailableButDownloadFails(t *testing.T) {
	env := newUpgradeTestEnv(t)
	types.Version = "5.16.0"

	fetchLatestRelease = func(_ string) ([]byte, error) {
		release := githubRelease{TagName: "forge-cli/v5.17.0"}
		return json.Marshal(release)
	}

	// Use a server that returns 404 for downloads
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	httpGet = func(_ string) (*http.Response, error) {
		return http.Get(server.URL)
	}

	action := upgradeCLIBinary(&env.stdout)

	if action.status != "FAILED" {
		t.Errorf("expected FAILED for download failure, got %q", action.status)
	}
}

func TestUpgradeCLIBinary_APIFailure(t *testing.T) {
	env := newUpgradeTestEnv(t)
	types.Version = "5.16.0"

	fetchLatestRelease = func(_ string) ([]byte, error) {
		return nil, fmt.Errorf("network error")
	}

	action := upgradeCLIBinary(&env.stdout)

	if action.status != "FAILED" {
		t.Errorf("expected FAILED for API failure, got %q", action.status)
	}
}

func TestUpgradeCLIBinary_InvalidTagFormat(t *testing.T) {
	env := newUpgradeTestEnv(t)
	types.Version = "5.16.0"

	fetchLatestRelease = func(_ string) ([]byte, error) {
		release := githubRelease{TagName: "invalid-tag"}
		return json.Marshal(release)
	}

	action := upgradeCLIBinary(&env.stdout)

	if action.status != "FAILED" {
		t.Errorf("expected FAILED for invalid tag, got %q", action.status)
	}
}

// --- Binary replacement tests ---

func TestReplaceUnixBinary(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-only test")
	}

	dir := t.TempDir()
	forgePath := filepath.Join(dir, "forge")

	// Create initial binary
	if err := os.WriteFile(forgePath, []byte("old-binary"), 0o755); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err := replaceUnixBinary(&buf, forgePath, []byte("new-binary-content"))
	if err != nil {
		t.Fatalf("replaceUnixBinary failed: %v", err)
	}

	// Verify content was replaced
	data, err := os.ReadFile(forgePath)
	if err != nil {
		t.Fatalf("failed to read replaced binary: %v", err)
	}
	if string(data) != "new-binary-content" {
		t.Errorf("expected 'new-binary-content', got %q", string(data))
	}

	// Verify .new file is cleaned up
	if _, err := os.Stat(forgePath + ".new"); !os.IsNotExist(err) {
		t.Error(".new file should be cleaned up after rename")
	}
}

func TestReplaceUnixBinary_RenameFailsThenRetrySucceeds(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-only test")
	}

	origRename := osRename
	defer func() { osRename = origRename }()

	callCount := 0
	osRename = func(oldpath, newpath string) error {
		callCount++
		if callCount == 1 {
			return os.ErrNotExist
		}
		return origRename(oldpath, newpath)
	}

	dir := t.TempDir()
	forgePath := filepath.Join(dir, "forge")
	if err := os.WriteFile(forgePath, []byte("old"), 0o755); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err := replaceUnixBinary(&buf, forgePath, []byte("new-content"))
	if err != nil {
		t.Fatalf("replaceUnixBinary should succeed via fallback: %v", err)
	}

	data, err := os.ReadFile(forgePath)
	if err != nil {
		t.Fatalf("failed to read binary: %v", err)
	}
	if string(data) != "new-content" {
		t.Errorf("expected 'new-content', got %q", string(data))
	}

	if callCount < 2 {
		t.Errorf("expected at least 2 rename attempts, got %d", callCount)
	}
}

func TestReplaceUnixBinary_RenameAlwaysFailsDirectWriteSucceeds(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-only test")
	}

	origRename := osRename
	defer func() { osRename = origRename }()

	osRename = func(_, _ string) error {
		return os.ErrNotExist
	}

	dir := t.TempDir()
	forgePath := filepath.Join(dir, "forge")
	if err := os.WriteFile(forgePath, []byte("old"), 0o755); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err := replaceUnixBinary(&buf, forgePath, []byte("direct-write-content"))
	if err != nil {
		t.Fatalf("replaceUnixBinary should succeed via direct write: %v", err)
	}

	data, err := os.ReadFile(forgePath)
	if err != nil {
		t.Fatalf("failed to read binary: %v", err)
	}
	if string(data) != "direct-write-content" {
		t.Errorf("expected 'direct-write-content', got %q", string(data))
	}

	// .new file should be cleaned up
	if _, err := os.Stat(forgePath + ".new"); !os.IsNotExist(err) {
		t.Error(".new file should be cleaned up after direct write")
	}
}

func TestReplaceUnixBinary_CleansUpStaleNewFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-only test")
	}

	dir := t.TempDir()
	forgePath := filepath.Join(dir, "forge")

	// Create a stale .new file from a previous failed attempt
	if err := os.WriteFile(forgePath+".new", []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err := replaceUnixBinary(&buf, forgePath, []byte("fresh-content"))
	if err != nil {
		t.Fatalf("replaceUnixBinary failed: %v", err)
	}

	data, err := os.ReadFile(forgePath)
	if err != nil {
		t.Fatalf("failed to read binary: %v", err)
	}
	if string(data) != "fresh-content" {
		t.Errorf("expected 'fresh-content', got %q", string(data))
	}
}

func TestDownloadAndReplace_NetworkError(t *testing.T) {
	origHTTPGet := httpGet
	defer func() { httpGet = origHTTPGet }()

	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	forgeBinaryDir = func() string { return binDir }

	httpGet = func(_ string) (*http.Response, error) {
		return nil, errors.New("connection refused")
	}

	var buf bytes.Buffer
	err := downloadAndReplace(&buf, "5.17.0")
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "download failed") {
		t.Errorf("expected 'download failed' error, got: %v", err)
	}
}

func TestDownloadAndReplace_ServerError(t *testing.T) {
	origHTTPGet := httpGet
	defer func() { httpGet = origHTTPGet }()

	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	forgeBinaryDir = func() string { return binDir }

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	httpGet = func(_ string) (*http.Response, error) {
		return http.Get(server.URL)
	}

	var buf bytes.Buffer
	err := downloadAndReplace(&buf, "5.17.0")
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
	if !strings.Contains(err.Error(), "HTTP 500") {
		t.Errorf("expected HTTP 500 error, got: %v", err)
	}
}

func TestDownloadAndReplace_IncompleteDownload(t *testing.T) {
	origHTTPGet := httpGet
	defer func() { httpGet = origHTTPGet }()

	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	forgeBinaryDir = func() string { return binDir }

	// Server sends Content-Length: 100 but only 50 bytes
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Length", "100")
		_, _ = w.Write([]byte("short-data-that-is-only-29-bytes"))
	}))
	defer server.Close()

	httpGet = func(_ string) (*http.Response, error) {
		return http.Get(server.URL)
	}

	var buf bytes.Buffer
	err := downloadAndReplace(&buf, "5.17.0")
	if err == nil {
		t.Fatal("expected error for incomplete download")
	}
	if !strings.Contains(err.Error(), "download incomplete") {
		t.Errorf("expected 'download incomplete' error, got: %v", err)
	}
}

func TestReplaceWindowsBinary(t *testing.T) {
	dir := t.TempDir()
	forgePath := filepath.Join(dir, "forge.exe")

	// Create initial binary
	if err := os.WriteFile(forgePath, []byte("old-binary"), 0o755); err != nil {
		t.Fatal(err)
	}

	content := []byte("new-binary-content")
	var buf bytes.Buffer
	err := replaceWindowsBinary(&buf, forgePath, content)
	if err != nil {
		t.Fatalf("replaceWindowsBinary failed: %v", err)
	}

	// Verify content was replaced
	data, err := os.ReadFile(forgePath)
	if err != nil {
		t.Fatalf("failed to read replaced binary: %v", err)
	}
	if string(data) != "new-binary-content" {
		t.Errorf("expected 'new-binary-content', got %q", string(data))
	}

	// Verify .old file is cleaned up
	if _, err := os.Stat(forgePath + ".old"); !os.IsNotExist(err) {
		t.Error("old .old file should be cleaned up after replacement")
	}
}

func TestReplaceWindowsBinary_FreshInstall(t *testing.T) {
	dir := t.TempDir()
	forgePath := filepath.Join(dir, "forge.exe")

	// No existing binary — fresh install scenario
	content := []byte("fresh-binary")
	var buf bytes.Buffer
	err := replaceWindowsBinary(&buf, forgePath, content)
	if err != nil {
		t.Fatalf("replaceWindowsBinary failed on fresh install: %v", err)
	}

	data, err := os.ReadFile(forgePath)
	if err != nil {
		t.Fatalf("failed to read fresh binary: %v", err)
	}
	if string(data) != "fresh-binary" {
		t.Errorf("expected 'fresh-binary', got %q", string(data))
	}
}

func TestReplaceWindowsBinary_StaleOldFile(t *testing.T) {
	dir := t.TempDir()
	forgePath := filepath.Join(dir, "forge.exe")
	oldPath := forgePath + ".old"

	// Create both current and stale .old
	if err := os.WriteFile(forgePath, []byte("current"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(oldPath, []byte("stale"), 0o755); err != nil {
		t.Fatal(err)
	}

	content := []byte("new")
	var buf bytes.Buffer
	err := replaceWindowsBinary(&buf, forgePath, content)
	if err != nil {
		t.Fatalf("replaceWindowsBinary failed with stale .old: %v", err)
	}

	data, _ := os.ReadFile(forgePath)
	if string(data) != "new" {
		t.Errorf("expected 'new', got %q", string(data))
	}

	// .old should be removed
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Error("stale .old file should be cleaned up")
	}
}

// --- Integration-level: runUpgrade end-to-end tests ---

func TestUpgradeCommand_PrerequisiteCheck(t *testing.T) {
	env := newUpgradeTestEnv(t)

	// Simulate claude not in PATH
	lookPathForUpgrade = func(_ string) (string, error) {
		return "", exec.ErrNotFound
	}

	err := env.run()
	if err == nil {
		t.Fatal("expected error when claude not in PATH")
	}
	if !strings.Contains(err.Error(), "prerequisite") {
		t.Errorf("expected prerequisite error, got: %v", err)
	}
	if !strings.Contains(env.stderr.String(), "claude CLI not found") {
		t.Errorf("expected 'claude CLI not found' in stderr, got: %q", env.stderr.String())
	}
}

func TestUpgradeCommand_CLISameVersion(t *testing.T) {
	env := newUpgradeTestEnv(t)
	types.Version = "5.17.0"

	lookPathForUpgrade = func(_ string) (string, error) {
		return "/usr/bin/claude", nil
	}
	fetchLatestRelease = func(_ string) ([]byte, error) {
		release := githubRelease{TagName: "forge-cli/v5.17.0"}
		return json.Marshal(release)
	}
	runClaudeCommand = func(_ []string) error {
		return nil
	}

	err := env.run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := env.stdout.String()
	if !strings.Contains(output, "SKIPPED") || !strings.Contains(output, "CLI binary") {
		t.Errorf("expected SKIPPED for CLI binary, got: %q", output)
	}
	if !strings.Contains(output, ">>>") || !strings.Contains(output, "<<<") {
		t.Errorf("expected summary block markers, got: %q", output)
	}
}

func TestUpgradeCommand_CLINewerVersion(t *testing.T) {
	env := newUpgradeTestEnv(t)
	types.Version = "5.16.0"
	tmpDir := t.TempDir()
	binDir := filepath.Join(tmpDir, ".forge", "bin")

	lookPathForUpgrade = func(_ string) (string, error) {
		return "/usr/bin/claude", nil
	}
	fetchLatestRelease = func(_ string) ([]byte, error) {
		release := githubRelease{TagName: "forge-cli/v5.17.0"}
		return json.Marshal(release)
	}
	forgeBinaryDir = func() string { return binDir }

	// Set up a test server that serves the download
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("fake-binary-content"))
	}))
	defer server.Close()

	httpGet = func(_ string) (*http.Response, error) {
		return http.Get(server.URL)
	}

	runClaudeCommand = func(_ []string) error {
		return nil
	}

	err := env.run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := env.stdout.String()
	if !strings.Contains(output, "UPGRADED") || !strings.Contains(output, "CLI binary") {
		t.Errorf("expected UPGRADED for CLI binary, got: %q", output)
	}
}

func TestUpgradeCommand_FullSummary(t *testing.T) {
	env := newUpgradeTestEnv(t)
	types.Version = "5.17.0"

	lookPathForUpgrade = func(_ string) (string, error) {
		return "/usr/bin/claude", nil
	}
	fetchLatestRelease = func(_ string) ([]byte, error) {
		release := githubRelease{TagName: "forge-cli/v5.17.0"}
		return json.Marshal(release)
	}
	runClaudeCommand = func(_ []string) error {
		return nil
	}

	err := env.run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := env.stdout.String()
	// Should have both CLI and Plugin results
	if !strings.Contains(output, "CLI binary") {
		t.Errorf("expected 'CLI binary' in output, got: %q", output)
	}
	if !strings.Contains(output, "Plugin") {
		t.Errorf("expected 'Plugin' in output, got: %q", output)
	}
	// Should have block markers
	if !strings.Contains(output, ">>>") || !strings.Contains(output, "<<<") {
		t.Errorf("expected summary block markers, got: %q", output)
	}
}

// --- Fetch latest release tests ---

func TestFetchLatestReleaseImpl_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		release := githubRelease{TagName: "forge-cli/v5.17.0"}
		data, _ := json.Marshal(release)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))
	defer server.Close()

	body, err := fetchLatestReleaseImpl(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var release githubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if release.TagName != "forge-cli/v5.17.0" {
		t.Errorf("expected tag 'forge-cli/v5.17.0', got %q", release.TagName)
	}
}

func TestFetchLatestReleaseImpl_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := fetchLatestReleaseImpl(server.URL)
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
	if !strings.Contains(err.Error(), "status 500") {
		t.Errorf("expected status 500 error, got: %v", err)
	}
}

// --- Download and replace integration test ---

func TestDownloadAndReplace_Integration(t *testing.T) {
	origHTTPGet := httpGet
	defer func() { httpGet = origHTTPGet }()

	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	forgeBinaryDir = func() string { return binDir }

	// Set up a test server that serves binary content
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "test-binary-data")
	}))
	defer server.Close()

	httpGet = func(_ string) (*http.Response, error) {
		return http.Get(server.URL)
	}

	var buf bytes.Buffer
	err := downloadAndReplace(&buf, "5.17.0")
	if err != nil {
		t.Fatalf("downloadAndReplace failed: %v", err)
	}

	// Verify the binary was written
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	binaryPath := filepath.Join(binDir, "forge"+ext)
	data, err := os.ReadFile(binaryPath)
	if err != nil {
		t.Fatalf("failed to read binary: %v", err)
	}
	if string(data) != "test-binary-data" {
		t.Errorf("expected 'test-binary-data', got %q", string(data))
	}
}

func TestDownloadAndReplace_404Error(t *testing.T) {
	origHTTPGet := httpGet
	defer func() { httpGet = origHTTPGet }()

	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	forgeBinaryDir = func() string { return binDir }

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	httpGet = func(_ string) (*http.Response, error) {
		return http.Get(server.URL)
	}

	var buf bytes.Buffer
	err := downloadAndReplace(&buf, "5.17.0")
	if err == nil {
		t.Fatal("expected error for HTTP 404")
	}
	if !strings.Contains(err.Error(), "HTTP 404") {
		t.Errorf("expected HTTP 404 error, got: %v", err)
	}
}

// --- Summary output test ---

func TestPrintUpgradeSummary(t *testing.T) {
	var buf bytes.Buffer
	actions := []upgradeAction{
		{status: "UPGRADED", target: "CLI binary", detail: "v5.16.0 -> v5.17.0"},
		{status: "INSTALLED", target: "Plugin", detail: "installed latest version"},
	}

	printUpgradeSummary(&buf, actions)

	output := buf.String()
	if !strings.Contains(output, ">>>") || !strings.Contains(output, "<<<") {
		t.Errorf("expected block markers, got: %q", output)
	}
	if !strings.Contains(output, "UPGRADED") {
		t.Errorf("expected UPGRADED status, got: %q", output)
	}
	if !strings.Contains(output, "INSTALLED") {
		t.Errorf("expected INSTALLED status, got: %q", output)
	}
	if !strings.Contains(output, "CLI binary") || !strings.Contains(output, "Plugin") {
		t.Errorf("expected component names, got: %q", output)
	}
}

func TestPrintUpgradeSummary_SkipAndFail(t *testing.T) {
	var buf bytes.Buffer
	actions := []upgradeAction{
		{status: "SKIPPED", target: "CLI binary", detail: "already up to date (v5.17.0)"},
		{status: "FAILED", target: "Plugin", detail: "install failed: timeout"},
	}

	printUpgradeSummary(&buf, actions)

	output := buf.String()
	if !strings.Contains(output, "SKIPPED") {
		t.Errorf("expected SKIPPED, got: %q", output)
	}
	if !strings.Contains(output, "FAILED") {
		t.Errorf("expected FAILED, got: %q", output)
	}
}

// --- DefaultForgeBinaryDir test ---

func TestDefaultForgeBinaryDir(t *testing.T) {
	dir := defaultForgeBinaryDir()
	if dir == "" {
		t.Fatal("expected non-empty directory")
	}
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".forge", "bin")
	if dir != expected {
		t.Errorf("expected %q, got %q", expected, dir)
	}
}

// --- Plugin update tests ---

func TestUpgradePlugin_UsesQualifiedName(t *testing.T) {
	env := newUpgradeTestEnv(t)
	types.Version = "5.17.0"

	lookPathForUpgrade = func(_ string) (string, error) {
		return "/usr/bin/claude", nil
	}
	fetchLatestRelease = func(_ string) ([]byte, error) {
		release := githubRelease{TagName: "forge-cli/v5.17.0"}
		return json.Marshal(release)
	}

	pluginInstalledCheck = func() bool { return true }

	var capturedArgs []string
	runClaudeCommand = func(args []string) error {
		capturedArgs = args
		return nil
	}

	action := upgradePlugin(&env.stdout)

	if action.status != "UPGRADED" {
		t.Errorf("expected UPGRADED, got %q", action.status)
	}
	// Verify fully qualified name "forge@forge" is used, not bare "forge"
	expectedArgs := []string{"plugin", "update", "forge@forge"}
	if len(capturedArgs) != len(expectedArgs) {
		t.Fatalf("expected %d args, got %d: %v", len(expectedArgs), len(capturedArgs), capturedArgs)
	}
	for i, got := range capturedArgs {
		if got != expectedArgs[i] {
			t.Errorf("arg[%d]: expected %q, got %q", i, expectedArgs[i], got)
		}
	}
}

func TestUpgradePlugin_InstallUsesQualifiedName(t *testing.T) {
	env := newUpgradeTestEnv(t)
	types.Version = "5.17.0"

	lookPathForUpgrade = func(_ string) (string, error) {
		return "/usr/bin/claude", nil
	}
	fetchLatestRelease = func(_ string) ([]byte, error) {
		release := githubRelease{TagName: "forge-cli/v5.17.0"}
		return json.Marshal(release)
	}

	pluginInstalledCheck = func() bool { return false }

	var capturedArgs []string
	runClaudeCommand = func(args []string) error {
		capturedArgs = args
		return nil
	}

	action := upgradePlugin(&env.stdout)

	if action.status != "INSTALLED" {
		t.Errorf("expected INSTALLED, got %q", action.status)
	}
	expectedArgs := []string{"plugin", "install", "forge@forge"}
	if len(capturedArgs) != len(expectedArgs) {
		t.Fatalf("expected %d args, got %d: %v", len(expectedArgs), len(capturedArgs), capturedArgs)
	}
	for i, got := range capturedArgs {
		if got != expectedArgs[i] {
			t.Errorf("arg[%d]: expected %q, got %q", i, expectedArgs[i], got)
		}
	}
}

func TestUpgradePlugin_UpdateFails(t *testing.T) {
	env := newUpgradeTestEnv(t)
	types.Version = "5.17.0"

	lookPathForUpgrade = func(_ string) (string, error) {
		return "/usr/bin/claude", nil
	}
	fetchLatestRelease = func(_ string) ([]byte, error) {
		release := githubRelease{TagName: "forge-cli/v5.17.0"}
		return json.Marshal(release)
	}

	pluginInstalledCheck = func() bool { return true }

	runClaudeCommand = func(_ []string) error {
		return fmt.Errorf("some error")
	}

	action := upgradePlugin(&env.stdout)

	if action.status != "FAILED" {
		t.Errorf("expected FAILED, got %q", action.status)
	}
	if !strings.Contains(action.detail, "update failed") {
		t.Errorf("expected 'update failed' in detail, got %q", action.detail)
	}
}
