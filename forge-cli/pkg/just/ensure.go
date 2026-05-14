// Package just provides utilities for interacting with justfile recipes,
// including quality gate execution, scope resolution, and just binary management.
package just

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	embedded "forge-cli/internal/embedded/just"
)

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// EnsureStatus represents the outcome of an ensure-just operation.
type EnsureStatus string

const (
	// StatusInstalled indicates just was successfully installed.
	StatusInstalled EnsureStatus = "INSTALLED"
	// StatusSkipped indicates just was already available and meets requirements.
	StatusSkipped EnsureStatus = "SKIPPED"
	// StatusFailed indicates the ensure-just operation failed.
	StatusFailed EnsureStatus = "FAILED"
)

// EnsureResult carries the result of the ensureJust flow.
type EnsureResult struct {
	Status  EnsureStatus // INSTALLED, SKIPPED, or FAILED
	Version string       // installed/found version, or detail message
	Method  string       // brew, cargo, scoop, choco, embedded, or ""
	Detail  string       // human-readable explanation
}

// Minimum required version of just.
// Task spec says >= 1.40.0.
var minimumVersion = "1.40.0"

// ---------------------------------------------------------------------------
// Detect
// ---------------------------------------------------------------------------

// DetectJust locates just in PATH and parses its version.
// Returns (path, version, found). If not found, path and version are empty.
func DetectJust() (path string, version string, found bool) {
	p, err := exec.LookPath("just")
	if err != nil {
		return "", "", false
	}

	out, err := exec.Command(p, "--version").Output()
	if err != nil {
		return p, "", true // found binary but version unknown
	}

	v, parseErr := ParseJustVersion(string(out))
	if parseErr != nil {
		return p, "", true
	}

	return p, v, true
}

// ---------------------------------------------------------------------------
// ParseJustVersion
// ---------------------------------------------------------------------------

// versionRe matches lines like "just 1.40.0" or "just 1.40.0-beta.1".
var versionRe = regexp.MustCompile(`^just\s+(\S+)`)

// ParseJustVersion parses the output of `just --version` and returns the
// semver-like version string (e.g., "1.40.0").
func ParseJustVersion(output string) (string, error) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		matches := versionRe.FindStringSubmatch(line)
		if len(matches) >= 2 && matches[1] != "" {
			return matches[1], nil
		}
	}
	return "", fmt.Errorf("could not parse just version from output: %q", output)
}

// ---------------------------------------------------------------------------
// IsMinimumVersion
// ---------------------------------------------------------------------------

// IsMinimumVersion reports whether version >= minimum using semantic version
// comparison (major.minor.patch). Pre-release suffixes are ignored for the
// numeric comparison but considered less than the release (e.g., 1.40.0-pre < 1.40.0).
func IsMinimumVersion(version, minimum string) bool {
	v := parseSemver(version)
	m := parseSemver(minimum)

	if v.major != m.major {
		return v.major > m.major
	}
	if v.minor != m.minor {
		return v.minor > m.minor
	}
	if v.patch != m.patch {
		return v.patch > m.patch
	}
	// Equal numeric parts. Pre-release versions are less than release.
	return !v.prerelease
}

type semver struct {
	major, minor, patch int
	prerelease          bool // true if version has a pre-release suffix
}

// semverRe captures major, minor, patch, and optional pre-release suffix.
var semverRe = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(-.+)?$`)

func parseSemver(s string) semver {
	matches := semverRe.FindStringSubmatch(s)
	if len(matches) < 4 {
		return semver{}
	}
	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])
	prerelease := len(matches) > 4 && matches[4] != ""
	return semver{major: major, minor: minor, patch: patch, prerelease: prerelease}
}

// ---------------------------------------------------------------------------
// Package Manager Installation
// ---------------------------------------------------------------------------

// detectPackageManager determines the best package manager to use for
// installing just. Returns one of: "brew", "cargo", "scoop", "choco", or "".
// This is a function variable for testability.
var detectPackageManager = detectPackageManagerImpl

func detectPackageManagerImpl() string {
	switch runtime.GOOS {
	case "darwin":
		if _, err := exec.LookPath("brew"); err == nil {
			return "brew"
		}
		if _, err := exec.LookPath("cargo"); err == nil {
			return "cargo"
		}
	case "linux":
		if _, err := exec.LookPath("cargo"); err == nil {
			return "cargo"
		}
	case "windows":
		if _, err := exec.LookPath("scoop"); err == nil {
			return "scoop"
		}
		if _, err := exec.LookPath("choco"); err == nil {
			return "choco"
		}
		if _, err := exec.LookPath("cargo"); err == nil {
			return "cargo"
		}
	}
	return ""
}

// packageManagerCommands maps package manager name to install command.
var packageManagerCommands = map[string][]string{
	"brew":  {"brew", "install", "just"},
	"cargo": {"cargo", "install", "just"},
	"scoop": {"scoop", "install", "just"},
	"choco": {"choco", "install", "just", "-y"},
}

// InstallViaPackageManagerFunc is the function that attempts to install just
// via the detected package manager. Variable for testability.
var InstallViaPackageManagerFunc = installViaPackageManagerImpl

// installViaPackageManagerImpl attempts to install just via the detected package
// manager. Returns an EnsureResult indicating success or failure.
func installViaPackageManagerImpl(_ string) EnsureResult {
	pm := detectPackageManager()
	if pm == "" {
		return EnsureResult{
			Status: StatusFailed,
			Detail: "no supported package manager found",
		}
	}

	cmdParts, ok := packageManagerCommands[pm]
	if !ok {
		return EnsureResult{
			Status: StatusFailed,
			Detail: fmt.Sprintf("unknown package manager: %s", pm),
		}
	}

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return EnsureResult{
			Status: StatusFailed,
			Method: pm,
			Detail: fmt.Sprintf("%s install failed: %s: %s", pm, err, strings.TrimSpace(string(output))),
		}
	}

	// Verify installation succeeded by running detect again.
	_, v, found := DetectJust()
	if !found {
		return EnsureResult{
			Status: StatusFailed,
			Method: pm,
			Detail: "just installed via " + pm + " but not found in PATH after installation",
		}
	}

	return EnsureResult{
		Status:  StatusInstalled,
		Version: v,
		Method:  pm,
		Detail:  fmt.Sprintf("installed via %s", pm),
	}
}

// ---------------------------------------------------------------------------
// Embedded Binary Extraction
// ---------------------------------------------------------------------------

// EmbeddedBinaryFunc is the function that returns the embedded just binary.
// Variable for testability.
var EmbeddedBinaryFunc = embedded.Binary

// userHomeDir is the function that returns the user's home directory.
// Variable for testability.
var userHomeDir = os.UserHomeDir

// ExtractEmbeddedBinaryFunc is the function that extracts the embedded just binary.
// Variable for testability.
var ExtractEmbeddedBinaryFunc = extractEmbeddedBinaryImpl

// extractEmbeddedBinaryImpl extracts the embedded just binary to ~/.forge/bin/.
// Returns an EnsureResult with the extraction outcome.
func extractEmbeddedBinaryImpl() EnsureResult {
	homeDir, err := userHomeDir()
	if err != nil {
		return EnsureResult{
			Status: StatusFailed,
			Detail: fmt.Sprintf("cannot determine home directory: %s", err),
		}
	}

	binData := EmbeddedBinaryFunc()
	if len(binData) == 0 {
		return EnsureResult{
			Status: StatusFailed,
			Detail: "no embedded just binary available for this platform",
		}
	}

	binDir := filepath.Join(homeDir, ".forge", "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return EnsureResult{
			Status: StatusFailed,
			Detail: fmt.Sprintf("cannot create %s: %s", binDir, err),
		}
	}

	binName := "just"
	if runtime.GOOS == "windows" {
		binName = "just.exe"
	}
	binPath := filepath.Join(binDir, binName)

	if err := os.WriteFile(binPath, binData, 0o755); err != nil {
		return EnsureResult{
			Status: StatusFailed,
			Detail: fmt.Sprintf("cannot write %s: %s", binPath, err),
		}
	}

	return EnsureResult{
		Status:  StatusInstalled,
		Version: binPath, // Store the path so tests can verify the file
		Method:  "embedded",
		Detail:  fmt.Sprintf("extracted to %s — add %s to your PATH", binPath, binDir),
	}
}

// ---------------------------------------------------------------------------
// EnsureJust — Main Orchestrator
// ---------------------------------------------------------------------------

// DetectJustFunc is the detection function used by EnsureJust.
// Variable for testability.
var DetectJustFunc = DetectJust

// isTerminalFunc checks whether the given reader is connected to a terminal.
// Hard rule: user confirmation MUST check that stdin is a terminal.
// Variable for testability.
var isTerminalFunc = isTerminalImpl

func isTerminalImpl(in io.Reader) bool {
	// Check if the reader is *os.File and if it's a terminal.
	f, ok := in.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// writeStr writes a string to the writer, discarding any error.
// Used for informational output where errors are not actionable.
func writeStr(w io.Writer, s string) {
	_, _ = io.WriteString(w, s)
}

// EnsureJust orchestrates the full detect → confirm → install flow.
// in: stdin for user confirmation prompts.
// out: stdout for status messages.
func EnsureJust(in io.Reader, out io.Writer) EnsureResult {
	// Step 1: Detect
	path, version, found := DetectJustFunc()

	if found && version != "" {
		if IsMinimumVersion(version, minimumVersion) {
			return EnsureResult{
				Status:  StatusSkipped,
				Version: version,
				Detail:  fmt.Sprintf("just %s found at %s (meets minimum %s)", version, path, minimumVersion),
			}
		}

		// Outdated version — prompt for upgrade.
		writeStr(out, fmt.Sprintf("WARNING: just %s found at %s, but minimum required is %s\n", version, path, minimumVersion))
		writeStr(out, "Upgrade? [y/N]: ")

		if !isTerminalFunc(in) {
			return EnsureResult{
				Status:  StatusFailed,
				Version: version,
				Detail:  "outdated just and non-interactive stdin (piped); cannot prompt for upgrade",
			}
		}

		reader := bufio.NewReader(in)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input != "y" && input != "yes" {
			return EnsureResult{
				Status:  StatusSkipped,
				Version: version,
				Detail:  fmt.Sprintf("user declined upgrade; just %s < %s", version, minimumVersion),
			}
		}

		// User wants upgrade — fall through to installation.
	}

	// Step 2: Confirm installation (if not already installed or upgrade requested).
	if !found {
		writeStr(out, "just is not installed. ")
	} else {
		writeStr(out, "Upgrading just. ")
	}

	// Hard rule: check stdin is terminal before prompting.
	if !isTerminalFunc(in) {
		return EnsureResult{
			Status: StatusFailed,
			Detail: "non-interactive stdin (piped); cannot prompt for installation",
		}
	}

	writeStr(out, "Install just? [Y/n]: ")
	reader := bufio.NewReader(in)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "n" || input == "no" {
		return EnsureResult{
			Status: StatusSkipped,
			Detail: "user declined installation",
		}
	}

	// Step 3: Try package manager first.
	result := InstallViaPackageManagerFunc(minimumVersion)
	if result.Status == StatusInstalled {
		writeStr(out, fmt.Sprintf("just installed successfully via %s\n", result.Method))
		return result
	}

	// Step 4: Fallback to embedded binary.
	writeStr(out, fmt.Sprintf("Package manager installation failed (%s). Trying embedded binary...\n", result.Detail))
	result = ExtractEmbeddedBinaryFunc()
	if result.Status == StatusInstalled {
		writeStr(out, "just installed via embedded binary.\n")
		return result
	}

	return result
}
