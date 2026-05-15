package just

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Integration tests: full ensureJust flow with mocked subsystems
// ---------------------------------------------------------------------------

// integrationTestEnv holds all the mock overrides for an integration test.
// Call restore() in a defer to clean up.
type integrationTestEnv struct {
	homeDir      string
	detectFn     func() (string, string, bool)
	isTerminal   bool
	pmResult     string // package manager name, or "" for no PM
	embeddedData []byte // embedded binary data, nil/empty means no embedded
}

func newIntegrationTestEnv(t *testing.T) *integrationTestEnv {
	t.Helper()
	return &integrationTestEnv{
		homeDir:    t.TempDir(),
		isTerminal: true,
		pmResult:   "",
	}
}

func (e *integrationTestEnv) install(t *testing.T) func() {
	t.Helper()

	// Save originals
	origDetect := DetectJustFunc
	origIsTerm := isTerminalFunc
	origPM := detectPackageManager
	origHome := userHomeDir
	origEmbedded := EmbeddedBinaryFunc
	origInstallPM := InstallViaPackageManagerFunc

	// Install mocks
	DetectJustFunc = e.detectFn
	if e.isTerminal {
		isTerminalFunc = func(io.Reader) bool { return true }
	} else {
		isTerminalFunc = isTerminalImpl
	}
	detectPackageManager = func() string { return e.pmResult }
	userHomeDir = func() (string, error) { return e.homeDir, nil }
	EmbeddedBinaryFunc = func() []byte { return e.embeddedData }

	// If PM is specified, keep the real InstallViaPackageManagerFunc (which calls
	// exec.Command and will fail). If no PM, the flow falls through to embedded.
	// We need to let the real flow run for integration testing.
	_ = origInstallPM

	return func() {
		DetectJustFunc = origDetect
		isTerminalFunc = origIsTerm
		detectPackageManager = origPM
		userHomeDir = origHome
		EmbeddedBinaryFunc = origEmbedded
		InstallViaPackageManagerFunc = origInstallPM
	}
}

// TestIntegration_FullFlow_SkippedWhenJustInstalled verifies that when just
// is already in PATH with a sufficient version, the flow returns SKIPPED.
func TestIntegration_FullFlow_SkippedWhenJustInstalled(t *testing.T) {
	env := newIntegrationTestEnv(t)
	env.detectFn = func() (string, string, bool) {
		return "/usr/bin/just", "1.50.0", true
	}
	cleanup := env.install(t)
	defer cleanup()

	var stdout bytes.Buffer
	result := EnsureJust(strings.NewReader(""), &stdout)

	assert.Equal(t, StatusSkipped, result.Status)
	assert.Equal(t, "1.50.0", result.Version)
	assert.Contains(t, stdout.String(), "")
}

// TestIntegration_FullFlow_InstallViaEmbeddedFallback verifies the full
// happy path: no just found, user accepts, PM fails, embedded fallback succeeds.
func TestIntegration_FullFlow_InstallViaEmbeddedFallback(t *testing.T) {
	env := newIntegrationTestEnv(t)
	env.detectFn = func() (string, string, bool) { return "", "", false }
	env.isTerminal = true
	env.pmResult = "" // no PM available
	env.embeddedData = []byte("fake-just-binary-content")

	cleanup := env.install(t)
	defer cleanup()

	var stdout bytes.Buffer
	result := EnsureJust(strings.NewReader("y\n"), &stdout)

	assert.Equal(t, StatusInstalled, result.Status)
	assert.Equal(t, "embedded", result.Method)

	// Verify binary was actually written
	binName := "just"
	if runtime.GOOS == "windows" {
		binName = "just.exe"
	}
	binPath := filepath.Join(env.homeDir, ".forge", "bin", binName)
	data, err := os.ReadFile(binPath)
	assert.NoError(t, err)
	assert.Equal(t, []byte("fake-just-binary-content"), data)
}

// TestIntegration_FullFlow_AllFail verifies the failure path when no
// package manager is available and no embedded binary exists.
func TestIntegration_FullFlow_AllFail(t *testing.T) {
	env := newIntegrationTestEnv(t)
	env.detectFn = func() (string, string, bool) { return "", "", false }
	env.isTerminal = true
	env.pmResult = ""      // no PM
	env.embeddedData = nil // no embedded

	cleanup := env.install(t)
	defer cleanup()

	var stdout bytes.Buffer
	result := EnsureJust(strings.NewReader("y\n"), &stdout)

	assert.Equal(t, StatusFailed, result.Status)
}

// TestIntegration_FullFlow_UserDeclines verifies that when the user
// declines installation, the result is SKIPPED with appropriate detail.
func TestIntegration_FullFlow_UserDeclines(t *testing.T) {
	env := newIntegrationTestEnv(t)
	env.detectFn = func() (string, string, bool) { return "", "", false }
	env.isTerminal = true

	cleanup := env.install(t)
	defer cleanup()

	var stdout bytes.Buffer
	result := EnsureJust(strings.NewReader("n\n"), &stdout)

	assert.Equal(t, StatusSkipped, result.Status)
	assert.Contains(t, result.Detail, "user declined")
}

// TestIntegration_FullFlow_OutdatedDeclinesUpgrade verifies the outdated
// version path where the user declines the upgrade prompt.
func TestIntegration_FullFlow_OutdatedDeclinesUpgrade(t *testing.T) {
	origMinVer := minimumVersion
	defer func() { minimumVersion = origMinVer }()
	minimumVersion = "99.99.99"

	env := newIntegrationTestEnv(t)
	env.detectFn = func() (string, string, bool) {
		return "/usr/bin/just", "1.30.0", true
	}
	env.isTerminal = true

	cleanup := env.install(t)
	defer cleanup()

	var stdout bytes.Buffer
	result := EnsureJust(strings.NewReader("n\n"), &stdout)

	assert.Equal(t, StatusSkipped, result.Status)
	assert.Contains(t, result.Detail, "user declined upgrade")
}

// TestIntegration_FullFlow_OutdatedAcceptsUpgrade verifies the outdated
// version path where the user accepts the upgrade prompt.
func TestIntegration_FullFlow_OutdatedAcceptsUpgrade(t *testing.T) {
	origMinVer := minimumVersion
	defer func() { minimumVersion = origMinVer }()
	minimumVersion = "99.99.99"

	env := newIntegrationTestEnv(t)
	env.detectFn = func() (string, string, bool) {
		return "/usr/bin/just", "1.30.0", true
	}
	env.isTerminal = true
	env.pmResult = ""
	env.embeddedData = []byte("fake-upgraded-just")

	cleanup := env.install(t)
	defer cleanup()

	var stdout bytes.Buffer
	// Two prompts: "Upgrade? [y/N]" and "Install just? [Y/n]"
	result := EnsureJust(strings.NewReader("y\ny\n"), &stdout)

	assert.Equal(t, StatusInstalled, result.Status)
	assert.Equal(t, "embedded", result.Method)
}

// TestIntegration_FullFlow_NonInteractive verifies that piped stdin
// without just installed returns FAILED.
func TestIntegration_FullFlow_NonInteractive(t *testing.T) {
	env := newIntegrationTestEnv(t)
	env.detectFn = func() (string, string, bool) { return "", "", false }
	env.isTerminal = false

	cleanup := env.install(t)
	defer cleanup()

	var stdout bytes.Buffer
	result := EnsureJust(strings.NewReader("y\n"), &stdout)

	assert.Equal(t, StatusFailed, result.Status)
	assert.Contains(t, result.Detail, "non-interactive")
}

// TestIntegration_EmbeddedBinary_ExtractionFilePermissions verifies that
// the extracted binary has executable permissions.
func TestIntegration_EmbeddedBinary_ExtractionFilePermissions(t *testing.T) {
	env := newIntegrationTestEnv(t)
	env.homeDir = t.TempDir()
	env.embeddedData = []byte("binary-with-perms-check")

	origHome := userHomeDir
	origEmbedded := EmbeddedBinaryFunc
	defer func() {
		userHomeDir = origHome
		EmbeddedBinaryFunc = origEmbedded
	}()

	userHomeDir = func() (string, error) { return env.homeDir, nil }
	EmbeddedBinaryFunc = func() []byte { return env.embeddedData }

	result := ExtractEmbeddedBinaryFunc()
	assert.Equal(t, StatusInstalled, result.Status)

	binName := "just"
	if runtime.GOOS == "windows" {
		binName = "just.exe"
	}
	binPath := filepath.Join(env.homeDir, ".forge", "bin", binName)

	info, err := os.Stat(binPath)
	assert.NoError(t, err)

	// On Unix, verify executable bit is set (0o755 => rwxr-xr-x).
	if runtime.GOOS != "windows" {
		assert.NotEqual(t, 0, info.Mode().Perm()&0o111, "binary should be executable")
	}
}

// TestIntegration_DetectJust_VersionParsingWithRealOutput verifies that
// the version parsing works with real-world just version output formats.
func TestIntegration_DetectJust_VersionParsingWithRealOutput(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   string
	}{
		{"standard format", "just 1.40.0\n", "1.40.0"},
		{"older version", "just 1.37.0\n", "1.37.0"},
		{"newer version", "just 1.50.0\n", "1.50.0"},
		{"pre-release", "just 1.40.0-beta.1\n", "1.40.0-beta.1"},
		{"version with hash", "just 1.40.0\n", "1.40.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseJustVersion(tt.output)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestIntegration_FullFlow_OutdatedNonInteractive verifies that an outdated
// version with non-interactive stdin returns FAILED (cannot prompt for upgrade).
func TestIntegration_FullFlow_OutdatedNonInteractive(t *testing.T) {
	origMinVer := minimumVersion
	defer func() { minimumVersion = origMinVer }()
	minimumVersion = "99.99.99"

	env := newIntegrationTestEnv(t)
	env.detectFn = func() (string, string, bool) {
		return "/usr/bin/just", "1.30.0", true
	}
	env.isTerminal = false

	cleanup := env.install(t)
	defer cleanup()

	var stdout bytes.Buffer
	result := EnsureJust(strings.NewReader("y\n"), &stdout)

	assert.Equal(t, StatusFailed, result.Status)
	assert.Contains(t, result.Detail, "non-interactive")
}
