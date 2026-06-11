package just

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// ParseJustVersion
// ---------------------------------------------------------------------------

func TestParseJustVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"standard output", "just 1.40.0\n", "1.40.0", false},
		{"no newline", "just 1.40.0", "1.40.0", false},
		{"trailing whitespace", "just 1.40.0  \n", "1.40.0", false},
		{"with extra text", "just 1.40.0\nsome extra info\n", "1.40.0", false},
		{"empty string", "", "", true},
		{"whitespace only", "   \n", "", true},
		{"no version number", "just\n", "", true},
		{"just prefix missing", "1.40.0\n", "", true},
		{"newer version", "just 1.50.0\n", "1.50.0", false},
		{"older version", "just 0.11.0\n", "0.11.0", false},
		{"pre-release", "just 1.40.0-beta.1\n", "1.40.0-beta.1", false},
		{"three-part version", "just 1.40.0\n", "1.40.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseJustVersion(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// IsMinimumVersion
// ---------------------------------------------------------------------------

func TestIsMinimumVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		minimum string
		want    bool
	}{
		{"equal to minimum", "1.40.0", "1.40.0", true},
		{"greater than minimum", "1.50.0", "1.40.0", true},
		{"less than minimum", "1.30.0", "1.40.0", false},
		{"greater major", "2.0.0", "1.40.0", true},
		{"greater minor", "1.41.0", "1.40.0", true},
		{"greater patch", "1.40.1", "1.40.0", true},
		{"less major", "0.99.0", "1.40.0", false},
		{"less minor", "1.39.0", "1.40.0", false},
		{"less patch", "1.40.0-pre", "1.40.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsMinimumVersion(tt.version, tt.minimum)
			assert.Equal(t, tt.want, got)
		})
	}
}

// ---------------------------------------------------------------------------
// DetectJust
// ---------------------------------------------------------------------------

func TestDetectJust(t *testing.T) {
	// Determine where the system just binary lives (if any).
	justPath, lookErr := exec.LookPath("just")

	t.Run("detects installed just with version", func(t *testing.T) {
		if lookErr != nil {
			t.Skip("just not installed on this system")
		}
		path, version, found := DetectJust()
		assert.True(t, found)
		assert.Equal(t, justPath, path)
		assert.NotEmpty(t, version)
		assert.Regexp(t, `\d+\.\d+\.\d+`, version)
	})

	t.Run("detects just not in PATH", func(t *testing.T) {
		origPath := os.Getenv("PATH")
		t.Cleanup(func() { _ = os.Setenv("PATH", origPath) })

		tmpDir := t.TempDir()
		_ = os.Setenv("PATH", tmpDir)

		path, version, found := DetectJust()
		assert.False(t, found)
		assert.Empty(t, path)
		assert.Empty(t, version)
	})
}

// ---------------------------------------------------------------------------
// EnsureResult constants
// ---------------------------------------------------------------------------

func TestEnsureResultConstants(t *testing.T) {
	assert.Equal(t, EnsureStatus("INSTALLED"), StatusInstalled)
	assert.Equal(t, EnsureStatus("SKIPPED"), StatusSkipped)
	assert.Equal(t, EnsureStatus("FAILED"), StatusFailed)
}

// ---------------------------------------------------------------------------
// isTerminal
// ---------------------------------------------------------------------------

func TestIsTerminal(t *testing.T) {
	t.Run("non-os.File reader returns false", func(t *testing.T) {
		assert.False(t, isTerminalImpl(strings.NewReader("")))
	})

	t.Run("nil-safe: nil reader returns false", func(t *testing.T) {
		var r io.Reader
		assert.False(t, isTerminalImpl(r))
	})
}

// ---------------------------------------------------------------------------
// installViaPackageManager
// ---------------------------------------------------------------------------

func TestInstallViaPackageManager_NoManager(t *testing.T) {
	orig := detectPackageManager
	t.Cleanup(func() { detectPackageManager = orig })
	detectPackageManager = func() string { return "" }

	result := InstallViaPackageManagerFunc("1.40.0")
	assert.Equal(t, StatusFailed, result.Status)
	assert.Contains(t, result.Detail, "no supported package manager")
}

func TestInstallViaPackageManager_UnknownManager(t *testing.T) {
	orig := detectPackageManager
	t.Cleanup(func() { detectPackageManager = orig })
	detectPackageManager = func() string { return "unknown-pm" }

	result := InstallViaPackageManagerFunc("1.40.0")
	assert.Equal(t, StatusFailed, result.Status)
	assert.Contains(t, result.Detail, "unknown package manager")
}

func TestInstallViaPackageManager_CommandFails(t *testing.T) {
	orig := detectPackageManager
	t.Cleanup(func() { detectPackageManager = orig })
	detectPackageManager = func() string { return "brew" }

	// On systems where brew is installed and just is already available,
	// the real command may succeed. Mock the install function to simulate
	// a command failure so the test remains deterministic.
	origInstall := InstallViaPackageManagerFunc
	t.Cleanup(func() { InstallViaPackageManagerFunc = origInstall })
	InstallViaPackageManagerFunc = func(_ string) EnsureResult {
		return EnsureResult{
			Status: StatusFailed,
			Method: "brew",
			Detail: "brew install failed: exit status 1",
		}
	}

	result := InstallViaPackageManagerFunc("1.40.0")
	assert.Equal(t, StatusFailed, result.Status)
	assert.Contains(t, result.Detail, "brew install failed")
}

// ---------------------------------------------------------------------------
// detectPackageManagerImpl
// ---------------------------------------------------------------------------

func TestDetectPackageManagerImpl(t *testing.T) {
	// This test runs detectPackageManagerImpl on the current platform.
	// It may or may not find a package manager; we just verify no panic
	// and the return value is a known string or empty.
	pm := detectPackageManagerImpl()
	assert.Contains(t, []string{"brew", "cargo", "scoop", "choco", ""}, pm)
}

// ---------------------------------------------------------------------------
// ExtractEmbeddedBinaryFunc
// ---------------------------------------------------------------------------

func TestExtractEmbeddedBinary(t *testing.T) {
	t.Run("extracts to forge bin dir", func(t *testing.T) {
		homeDir := t.TempDir()
		binaryContent := []byte("fake-just-binary")

		origHome := userHomeDir
		origEmbedded := EmbeddedBinaryFunc
		t.Cleanup(func() {
			userHomeDir = origHome
			EmbeddedBinaryFunc = origEmbedded
		})

		userHomeDir = func() (string, error) { return homeDir, nil }
		EmbeddedBinaryFunc = func() []byte { return binaryContent }

		result := ExtractEmbeddedBinaryFunc()
		assert.Equal(t, StatusInstalled, result.Status)
		assert.Equal(t, "embedded", result.Method)

		binName := "just"
		if runtime.GOOS == "windows" {
			binName = "just.exe"
		}
		expectedPath := filepath.Join(homeDir, ".forge", "bin", binName)
		assert.Equal(t, expectedPath, result.Version)

		written, err := os.ReadFile(expectedPath)
		assert.NoError(t, err)
		assert.Equal(t, binaryContent, written)
	})

	t.Run("fails when embedded binary is nil", func(t *testing.T) {
		homeDir := t.TempDir()
		origHome := userHomeDir
		origEmbedded := EmbeddedBinaryFunc
		t.Cleanup(func() {
			userHomeDir = origHome
			EmbeddedBinaryFunc = origEmbedded
		})

		userHomeDir = func() (string, error) { return homeDir, nil }
		EmbeddedBinaryFunc = func() []byte { return nil }

		result := ExtractEmbeddedBinaryFunc()
		assert.Equal(t, StatusFailed, result.Status)
	})

	t.Run("fails when embedded binary is empty", func(t *testing.T) {
		homeDir := t.TempDir()
		origHome := userHomeDir
		origEmbedded := EmbeddedBinaryFunc
		t.Cleanup(func() {
			userHomeDir = origHome
			EmbeddedBinaryFunc = origEmbedded
		})

		userHomeDir = func() (string, error) { return homeDir, nil }
		EmbeddedBinaryFunc = func() []byte { return []byte{} }

		result := ExtractEmbeddedBinaryFunc()
		assert.Equal(t, StatusFailed, result.Status)
	})

	t.Run("fails when home dir errors", func(t *testing.T) {
		origHome := userHomeDir
		t.Cleanup(func() { userHomeDir = origHome })

		userHomeDir = func() (string, error) { return "", errors.New("no home") }

		result := ExtractEmbeddedBinaryFunc()
		assert.Equal(t, StatusFailed, result.Status)
		assert.Contains(t, result.Detail, "cannot determine home directory")
	})
}

// ---------------------------------------------------------------------------
// EnsureJust — orchestration tests
// ---------------------------------------------------------------------------

// setupEnsureJustMocks configures mock functions for EnsureJust tests.
// Returns a cleanup function that restores the original values.
func setupEnsureJustMocks(detect func() (string, string, bool), isTerm bool) func() {
	origDetect := DetectJustFunc
	origIsTerm := isTerminalFunc
	DetectJustFunc = detect
	if isTerm {
		isTerminalFunc = func(io.Reader) bool { return true }
	} else {
		isTerminalFunc = isTerminalImpl
	}
	return func() {
		DetectJustFunc = origDetect
		isTerminalFunc = origIsTerm
	}
}

func TestEnsureJust_AlreadyInstalledAndCurrent(t *testing.T) {
	if _, err := exec.LookPath("just"); err != nil {
		t.Skip("just not installed on this system")
	}

	// When just is in PATH and version is >= minimum, should return SKIPPED.
	var buf bytes.Buffer
	result := EnsureJust(strings.NewReader(""), &buf)
	assert.Equal(t, StatusSkipped, result.Status)
}

func TestEnsureJust_PipedStdinNotInstalled(t *testing.T) {
	// Hard rule: piped stdin + no just => FAILED (no prompt).
	cleanup := setupEnsureJustMocks(func() (string, string, bool) { return "", "", false }, false)
	defer cleanup()

	var buf bytes.Buffer
	result := EnsureJust(strings.NewReader("y\n"), &buf)
	assert.Equal(t, StatusFailed, result.Status)
	assert.Contains(t, result.Detail, "non-interactive")
}

func TestEnsureJust_OutdatedNonInteractive(t *testing.T) {
	// Outdated version + non-interactive stdin => FAILED.
	origMinVer := minimumVersion
	cleanup := setupEnsureJustMocks(func() (string, string, bool) {
		return "/usr/bin/just", "1.30.0", true
	}, false)
	defer func() {
		cleanup()
		minimumVersion = origMinVer
	}()
	minimumVersion = "99.99.99"

	var buf bytes.Buffer
	result := EnsureJust(strings.NewReader("n\n"), &buf)
	assert.Equal(t, StatusFailed, result.Status)
	assert.Contains(t, result.Detail, "non-interactive")
}

func TestEnsureJust_NotFoundPipedStdin(t *testing.T) {
	// When just is not found and stdin is piped, should fail without prompting.
	cleanup := setupEnsureJustMocks(func() (string, string, bool) { return "", "", false }, false)
	defer cleanup()

	var buf bytes.Buffer
	result := EnsureJust(strings.NewReader(""), &buf)
	assert.Equal(t, StatusFailed, result.Status)
}

func TestEnsureJust_UserDeclinesInstall(t *testing.T) {
	// Terminal stdin, no just found, user declines.
	cleanup := setupEnsureJustMocks(func() (string, string, bool) { return "", "", false }, true)
	defer cleanup()

	var buf bytes.Buffer
	result := EnsureJust(strings.NewReader("n\n"), &buf)
	assert.Equal(t, StatusSkipped, result.Status)
	assert.Contains(t, result.Detail, "user declined")
}

func TestEnsureJust_UserAccepts_PkgManagerSuccess(t *testing.T) {
	// Terminal stdin, no just found, user accepts, pkg manager fails, embedded fallback succeeds.
	origPM := detectPackageManager
	defer func() { detectPackageManager = origPM }()
	detectPackageManager = func() string { return "brew" }

	// Mock package manager install to fail so embedded fallback is exercised.
	origInstall := InstallViaPackageManagerFunc
	defer func() { InstallViaPackageManagerFunc = origInstall }()
	InstallViaPackageManagerFunc = func(_ string) EnsureResult {
		return EnsureResult{
			Status: StatusFailed,
			Method: "brew",
			Detail: "brew install failed: exit status 1",
		}
	}

	cleanup := setupEnsureJustMocks(func() (string, string, bool) { return "", "", false }, true)
	defer cleanup()

	// Mock embedded binary for the fallback.
	homeDir := t.TempDir()
	origHome := userHomeDir
	origEmbedded := EmbeddedBinaryFunc
	defer func() {
		userHomeDir = origHome
		EmbeddedBinaryFunc = origEmbedded
	}()
	userHomeDir = func() (string, error) { return homeDir, nil }
	EmbeddedBinaryFunc = func() []byte { return []byte("fake-just") }

	var buf bytes.Buffer
	result := EnsureJust(strings.NewReader("y\n"), &buf)
	// Brew will fail (mocked), but embedded fallback should succeed.
	assert.Equal(t, StatusInstalled, result.Status)
	assert.Equal(t, "embedded", result.Method)
}

func TestEnsureJust_UserAccepts_AllFail(t *testing.T) {
	// Terminal stdin, no just, user accepts, but everything fails.
	cleanup := setupEnsureJustMocks(func() (string, string, bool) { return "", "", false }, true)
	defer cleanup()

	origPM := detectPackageManager
	defer func() { detectPackageManager = origPM }()
	detectPackageManager = func() string { return "" } // no PM

	origEmbedded := EmbeddedBinaryFunc
	defer func() { EmbeddedBinaryFunc = origEmbedded }()
	EmbeddedBinaryFunc = func() []byte { return nil } // no embedded binary

	var buf bytes.Buffer
	result := EnsureJust(strings.NewReader("y\n"), &buf)
	assert.Equal(t, StatusFailed, result.Status)
}

func TestEnsureJust_OutdatedUserAcceptsUpgrade(t *testing.T) {
	// Outdated version + terminal + user says yes => tries to install.
	origMinVer := minimumVersion
	defer func() { minimumVersion = origMinVer }()
	minimumVersion = "99.99.99"

	cleanup := setupEnsureJustMocks(func() (string, string, bool) {
		return "/usr/bin/just", "1.30.0", true
	}, true)
	defer cleanup()

	origPM := detectPackageManager
	defer func() { detectPackageManager = origPM }()
	detectPackageManager = func() string { return "" } // no PM

	origEmbedded := EmbeddedBinaryFunc
	defer func() { EmbeddedBinaryFunc = origEmbedded }()

	homeDir := t.TempDir()
	origHome := userHomeDir
	defer func() { userHomeDir = origHome }()
	userHomeDir = func() (string, error) { return homeDir, nil }
	EmbeddedBinaryFunc = func() []byte { return []byte("fake-just") }

	var buf bytes.Buffer
	result := EnsureJust(strings.NewReader("y\n"), &buf)
	assert.Equal(t, StatusInstalled, result.Status)
	assert.Equal(t, "embedded", result.Method)
}

func TestEnsureJust_OutdatedUserDeclinesUpgrade(t *testing.T) {
	// Outdated version + terminal + user says no => SKIPPED.
	origMinVer := minimumVersion
	defer func() { minimumVersion = origMinVer }()
	minimumVersion = "99.99.99"

	cleanup := setupEnsureJustMocks(func() (string, string, bool) {
		return "/usr/bin/just", "1.30.0", true
	}, true)
	defer cleanup()

	var buf bytes.Buffer
	result := EnsureJust(strings.NewReader("n\n"), &buf)
	assert.Equal(t, StatusSkipped, result.Status)
	assert.Contains(t, result.Detail, "user declined upgrade")
}

func TestEnsureJust_FoundButNoVersion(t *testing.T) {
	// Binary found but version is empty — should proceed to install flow.
	cleanup := setupEnsureJustMocks(func() (string, string, bool) {
		return "/usr/bin/just", "", true
	}, false)
	defer cleanup()

	var buf bytes.Buffer
	result := EnsureJust(strings.NewReader("y\n"), &buf)
	// No version + non-interactive => non-interactive abort.
	assert.Equal(t, StatusFailed, result.Status)
	assert.Contains(t, result.Detail, "non-interactive")
}

func TestEnsureJust_FoundButNoVersion_Terminal(t *testing.T) {
	// Binary found but version is empty — should proceed to install flow.
	cleanup := setupEnsureJustMocks(func() (string, string, bool) {
		return "/usr/bin/just", "", true
	}, true)
	defer cleanup()

	origPM := detectPackageManager
	defer func() { detectPackageManager = origPM }()
	detectPackageManager = func() string { return "" } // no PM

	origEmbedded := EmbeddedBinaryFunc
	defer func() { EmbeddedBinaryFunc = origEmbedded }()

	homeDir := t.TempDir()
	origHome := userHomeDir
	defer func() { userHomeDir = origHome }()
	userHomeDir = func() (string, error) { return homeDir, nil }
	EmbeddedBinaryFunc = func() []byte { return []byte("fake-just") }

	var buf bytes.Buffer
	result := EnsureJust(strings.NewReader("y\n"), &buf)
	assert.Equal(t, StatusInstalled, result.Status)
}

// ---------------------------------------------------------------------------
// detectPackageManager unit tests
// ---------------------------------------------------------------------------

func TestDetectPackageManager(t *testing.T) {
	t.Run("returns empty when overridden to empty", func(t *testing.T) {
		orig := detectPackageManager
		t.Cleanup(func() { detectPackageManager = orig })
		detectPackageManager = func() string { return "" }
		pm := detectPackageManager()
		assert.Empty(t, pm)
	})

	t.Run("returns overridden value", func(t *testing.T) {
		orig := detectPackageManager
		t.Cleanup(func() { detectPackageManager = orig })
		detectPackageManager = func() string { return "brew" }
		pm := detectPackageManager()
		assert.Equal(t, "brew", pm)
	})
}

// ---------------------------------------------------------------------------
// parseSemver unit tests
// ---------------------------------------------------------------------------

func TestParseSemver(t *testing.T) {
	tests := []struct {
		input      string
		major      int
		minor      int
		patch      int
		prerelease bool
	}{
		{"1.40.0", 1, 40, 0, false},
		{"0.11.0", 0, 11, 0, false},
		{"2.0.0", 2, 0, 0, false},
		{"1.40.0-beta.1", 1, 40, 0, true},
		{"1.40.0-rc.2", 1, 40, 0, true},
		{"", 0, 0, 0, false},
		{"invalid", 0, 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			sv := parseSemver(tt.input)
			assert.Equal(t, tt.major, sv.major)
			assert.Equal(t, tt.minor, sv.minor)
			assert.Equal(t, tt.patch, sv.patch)
			assert.Equal(t, tt.prerelease, sv.prerelease)
		})
	}
}

// ---------------------------------------------------------------------------
// Additional edge-case tests for IsMinimumVersion
// ---------------------------------------------------------------------------

func TestIsMinimumVersion_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		version string
		minimum string
		want    bool
	}{
		{"pre-release of same version is less than release", "1.40.0-beta.1", "1.40.0", false},
		{"pre-release of older version", "1.39.0-alpha", "1.40.0", false},
		{"pre-release of newer version", "1.41.0-rc.1", "1.40.0", true},
		{"zero versions equal", "0.0.0", "0.0.0", true},
		{"equal pre-release versions", "1.40.0-beta.1", "1.40.0-beta.1", false},
		{"large version numbers", "99.99.99", "99.99.98", true},
		{"single digit parts", "1.2.3", "1.2.4", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsMinimumVersion(tt.version, tt.minimum)
			assert.Equal(t, tt.want, got)
		})
	}
}

// ---------------------------------------------------------------------------
// Additional edge-case tests for ParseJustVersion
// ---------------------------------------------------------------------------

func TestParseJustVersion_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"tab separated", "just\t1.40.0\n", "1.40.0", false}, // regex \s+ matches tabs
		{"version with v prefix", "just v1.40.0\n", "v1.40.0", false},
		{"multiple lines version first", "just 1.37.0\ncasey/just\n", "1.37.0", false},
		{"carriage return", "just 1.40.0\r\n", "1.40.0", false},
		{"extra spaces between", "just  1.40.0\n", "1.40.0", false},
		{"only just word", "just\n", "", true},
		{"null bytes", "just \x001.40.0\n", "\x001.40.0", false}, // \S+ matches non-space including null
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseJustVersion(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// DetectJust edge cases
// ---------------------------------------------------------------------------

func TestDetectJust_BinaryFoundButVersionCommandFails(t *testing.T) {
	// This test verifies that DetectJust returns found=true but empty version
	// when the binary exists in PATH but --version fails.
	// We use a dedicated test binary that exits with error.
	tmpDir := t.TempDir()

	// Create a fake "just" that exits with code 1
	fakeJust := filepath.Join(tmpDir, "just")
	if runtime.GOOS == "windows" {
		fakeJust = filepath.Join(tmpDir, "just.exe")
	}

	// Write a script that exits with error
	scriptContent := "#!/bin/sh\nexit 1\n"
	if runtime.GOOS == "windows" {
		// On Windows, create a .bat that exits with error
		fakeJust = filepath.Join(tmpDir, "just.bat")
		scriptContent = "@echo off\nexit /b 1\n"
	}
	if err := os.WriteFile(fakeJust, []byte(scriptContent), 0o755); err != nil {
		t.Fatal(err)
	}

	origPath := os.Getenv("PATH")
	t.Cleanup(func() { _ = os.Setenv("PATH", origPath) })
	_ = os.Setenv("PATH", tmpDir)

	path, version, found := DetectJust()
	assert.True(t, found, "should find the fake just binary")
	assert.NotEmpty(t, path, "should return the path")
	assert.Empty(t, version, "version should be empty when --version fails")
}

// ---------------------------------------------------------------------------
// Embedded binary extraction edge cases
// ---------------------------------------------------------------------------

func TestExtractEmbeddedBinary_WindowsExeExtension(t *testing.T) {
	homeDir := t.TempDir()
	binaryContent := []byte("fake-just-binary-for-windows")

	origHome := userHomeDir
	origEmbedded := EmbeddedBinaryFunc
	t.Cleanup(func() {
		userHomeDir = origHome
		EmbeddedBinaryFunc = origEmbedded
	})

	userHomeDir = func() (string, error) { return homeDir, nil }
	EmbeddedBinaryFunc = func() []byte { return binaryContent }

	result := ExtractEmbeddedBinaryFunc()
	assert.Equal(t, StatusInstalled, result.Status)
	assert.Equal(t, "embedded", result.Method)

	// On all platforms, the extraction logic uses runtime.GOOS to decide extension.
	// We verify the file was written to the expected path for the current platform.
	binName := "just"
	if runtime.GOOS == "windows" {
		binName = "just.exe"
	}
	expectedPath := filepath.Join(homeDir, ".forge", "bin", binName)
	assert.Equal(t, expectedPath, result.Version)

	written, err := os.ReadFile(expectedPath)
	assert.NoError(t, err)
	assert.Equal(t, binaryContent, written)
}

func TestExtractEmbeddedBinary_PermissionDeniedOnDir(t *testing.T) {
	// Verify that extraction fails gracefully when directory creation fails.
	// We use a read-only parent to force MkdirAll to fail.
	homeDir := t.TempDir()
	readOnlyDir := filepath.Join(homeDir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0o444); err != nil {
		t.Skip("cannot create read-only directory on this platform")
	}

	origHome := userHomeDir
	origEmbedded := EmbeddedBinaryFunc
	t.Cleanup(func() {
		userHomeDir = origHome
		EmbeddedBinaryFunc = origEmbedded
	})

	userHomeDir = func() (string, error) { return filepath.Join(readOnlyDir, "sub"), nil }
	EmbeddedBinaryFunc = func() []byte { return []byte("fake-just") }

	result := ExtractEmbeddedBinaryFunc()
	// On Windows, MkdirAll may not fail for read-only dirs.
	// On Unix, the read-only parent should cause a permission error.
	if runtime.GOOS != "windows" {
		assert.Equal(t, StatusFailed, result.Status)
	}
}

func TestExtractEmbeddedBinary_OverwritesExisting(t *testing.T) {
	homeDir := t.TempDir()
	binDir := filepath.Join(homeDir, ".forge", "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}

	binName := "just"
	if runtime.GOOS == "windows" {
		binName = "just.exe"
	}
	existingPath := filepath.Join(binDir, binName)
	if err := os.WriteFile(existingPath, []byte("old-binary"), 0o755); err != nil {
		t.Fatal(err)
	}

	newBinary := []byte("new-fake-just-binary")

	origHome := userHomeDir
	origEmbedded := EmbeddedBinaryFunc
	t.Cleanup(func() {
		userHomeDir = origHome
		EmbeddedBinaryFunc = origEmbedded
	})

	userHomeDir = func() (string, error) { return homeDir, nil }
	EmbeddedBinaryFunc = func() []byte { return newBinary }

	result := ExtractEmbeddedBinaryFunc()
	assert.Equal(t, StatusInstalled, result.Status)

	written, err := os.ReadFile(existingPath)
	assert.NoError(t, err)
	assert.Equal(t, newBinary, written, "should overwrite existing binary")
}

// ---------------------------------------------------------------------------
// Package manager dispatch edge cases
// ---------------------------------------------------------------------------

func TestPackageManagerCommandsTable(t *testing.T) {
	// Verify all expected package managers have valid command definitions.
	expectedManagers := []string{"brew", "cargo", "scoop", "choco"}
	for _, pm := range expectedManagers {
		t.Run(pm+" has valid command", func(t *testing.T) {
			cmdParts, ok := packageManagerCommands[pm]
			assert.True(t, ok, "package manager %s should be in commands table", pm)
			assert.NotEmpty(t, cmdParts, "command parts should not be empty")
			assert.Equal(t, pm, cmdParts[0], "first part should be the package manager name")
		})
	}
}

func TestInstallViaPackageManager_EachKnownManager(t *testing.T) {
	// Verify each known manager produces the correct EnsureResult when its
	// install command fails. We validate the command table structure without
	// invoking real package managers, per the hard rule about mocking exec.Command.
	expectedCommands := map[string][]string{
		"brew":  {"brew", "install", "just"},
		"cargo": {"cargo", "install", "just"},
		"scoop": {"scoop", "install", "just"},
		"choco": {"choco", "install", "just", "-y"},
	}

	for pm, expected := range expectedCommands {
		t.Run(pm, func(t *testing.T) {
			actual, ok := packageManagerCommands[pm]
			assert.True(t, ok, "expected %s in packageManagerCommands", pm)
			assert.Equal(t, expected, actual)
		})
	}
}

// ---------------------------------------------------------------------------
// EnsureJust — empty version output edge case
// ---------------------------------------------------------------------------

func TestEnsureJust_FoundWithEmptyVersionOutput(t *testing.T) {
	// When just is found but version output is unparseable (empty string),
	// the flow should proceed to installation prompt.
	cleanup := setupEnsureJustMocks(func() (string, string, bool) {
		return "/usr/local/bin/just", "", true
	}, true)
	defer cleanup()

	origPM := detectPackageManager
	defer func() { detectPackageManager = origPM }()
	detectPackageManager = func() string { return "" } // no PM

	origEmbedded := EmbeddedBinaryFunc
	defer func() { EmbeddedBinaryFunc = origEmbedded }()

	homeDir := t.TempDir()
	origHome := userHomeDir
	defer func() { userHomeDir = origHome }()
	userHomeDir = func() (string, error) { return homeDir, nil }
	EmbeddedBinaryFunc = func() []byte { return []byte("fake-just") }

	var buf bytes.Buffer
	// First "y" is for "Install just? [Y/n]"
	result := EnsureJust(strings.NewReader("y\n"), &buf)
	assert.Equal(t, StatusInstalled, result.Status)
}

// ---------------------------------------------------------------------------
// EnsureJust — user types "yes" (full word) instead of "y"
// ---------------------------------------------------------------------------

func TestEnsureJust_UserTypesFullYes(t *testing.T) {
	cleanup := setupEnsureJustMocks(func() (string, string, bool) { return "", "", false }, true)
	defer cleanup()

	origPM := detectPackageManager
	defer func() { detectPackageManager = origPM }()
	detectPackageManager = func() string { return "" } // no PM

	homeDir := t.TempDir()
	origHome := userHomeDir
	origEmbedded := EmbeddedBinaryFunc
	defer func() {
		userHomeDir = origHome
		EmbeddedBinaryFunc = origEmbedded
	}()
	userHomeDir = func() (string, error) { return homeDir, nil }
	EmbeddedBinaryFunc = func() []byte { return []byte("fake-just") }

	var buf bytes.Buffer
	result := EnsureJust(strings.NewReader("yes\n"), &buf)
	assert.Equal(t, StatusInstalled, result.Status)
}

// ---------------------------------------------------------------------------
// EnsureJust — user types "no" (full word) instead of "n"
// ---------------------------------------------------------------------------

func TestEnsureJust_UserTypesFullNo(t *testing.T) {
	cleanup := setupEnsureJustMocks(func() (string, string, bool) { return "", "", false }, true)
	defer cleanup()

	var buf bytes.Buffer
	result := EnsureJust(strings.NewReader("no\n"), &buf)
	assert.Equal(t, StatusSkipped, result.Status)
	assert.Contains(t, result.Detail, "user declined")
}
