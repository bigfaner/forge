package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/forgelog"
	"forge-cli/pkg/types"

	"github.com/spf13/cobra"
)

// upgradeAction records a single action taken during upgrade.
type upgradeAction struct {
	status string // UPGRADED, SKIPPED, INSTALLED, FAILED
	target string // component name
	detail string // extra info
}

// githubRelease represents the subset of the GitHub Release API response we need.
type githubRelease struct {
	TagName string `json:"tag_name"`
}

// Plugin marketplace URL.
const pluginMarketplaceURL = "https://github.com/bigfaner/forge.git"

// downloadTimeout is the maximum time allowed for downloading the CLI binary.
const downloadTimeout = 5 * time.Minute

// Variables for testability — overridden in tests.
var (
	// lookPathForUpgrade resolves a binary name to its full path.
	lookPathForUpgrade = exec.LookPath

	// fetchLatestRelease calls the GitHub Release API and returns the raw response body.
	fetchLatestRelease = fetchLatestReleaseImpl

	// runClaudeCommand executes claude with the given args.
	runClaudeCommand = base.RunClaude

	// pluginInstalledCheck reports whether the forge plugin is installed.
	pluginInstalledCheck = pluginInstalledImpl

	// forgeBinaryDir returns the directory containing the forge binary.
	// Defaults to ~/.forge/bin.
	forgeBinaryDir = defaultForgeBinaryDir

	// httpGet performs an HTTP GET request.
	// Uses a client with timeout to avoid hanging on weak networks.
	httpGet = (&http.Client{Timeout: downloadTimeout}).Get

	// osRename renames (moves) oldpath to newpath.
	osRename = os.Rename
)

func defaultForgeBinaryDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".forge", "bin")
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade forge CLI binary and Plugin",
	Long: `Upgrade forge CLI binary to the latest version from GitHub Releases
and install or update the forge Plugin via Claude Code.

Prerequisite: 'claude' CLI must be in PATH.`,
	Args: cobra.NoArgs,
	RunE: runUpgrade,
}

func runUpgrade(cmd *cobra.Command, _ []string) error {
	out := cmd.OutOrStdout()
	var actions []upgradeAction

	// Prerequisite check: claude must be in PATH
	if _, err := lookPathForUpgrade("claude"); err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "error: claude CLI not found in PATH\n")
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "hint: install Claude Code first, then run 'forge upgrade'\n")
		return fmt.Errorf("prerequisite: %w", err)
	}

	// Phase 1: CLI binary upgrade
	action := upgradeCLIBinary(out)
	actions = append(actions, action)

	// Phase 2: Plugin management
	action = upgradePlugin(out)
	actions = append(actions, action)

	// Print summary
	printUpgradeSummary(out, actions)

	return nil
}

// upgradeCLIBinary handles the CLI binary upgrade phase.
// Compares current version with latest from GitHub Release API.
// Downloads and replaces binary if newer version is available.
func upgradeCLIBinary(out io.Writer) upgradeAction {
	currentVersion := types.GetVersion()
	if currentVersion == "dev" {
		return upgradeAction{status: "SKIPPED", target: "CLI binary", detail: "development build (version=dev)"}
	}

	// Fetch latest release info
	body, err := fetchLatestRelease("https://api.github.com/repos/bigfaner/forge/releases/latest")
	if err != nil {
		forgelog.Error("error: failed to fetch latest release: %v\n", err)
		return upgradeAction{status: "FAILED", target: "CLI binary", detail: err.Error()}
	}

	var release githubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		forgelog.Error("error: failed to parse release info: %v\n", err)
		return upgradeAction{status: "FAILED", target: "CLI binary", detail: err.Error()}
	}

	// Extract version from tag: "forge-cli/v5.17.0" -> "5.17.0"
	latestVersion, err := parseVersionFromTag(release.TagName)
	if err != nil {
		forgelog.Error("error: %v\n", err)
		return upgradeAction{status: "FAILED", target: "CLI binary", detail: err.Error()}
	}

	// Compare versions
	cmp := compareVersions(currentVersion, latestVersion)
	if cmp >= 0 {
		return upgradeAction{
			status: "SKIPPED",
			target: "CLI binary",
			detail: fmt.Sprintf("already up to date (v%s)", currentVersion),
		}
	}

	// Download and replace binary
	if err := downloadAndReplace(out, latestVersion); err != nil {
		forgelog.Error("error: CLI upgrade failed: %v\n", err)
		return upgradeAction{status: "FAILED", target: "CLI binary", detail: err.Error()}
	}

	return upgradeAction{
		status: "UPGRADED",
		target: "CLI binary",
		detail: fmt.Sprintf("v%s -> v%s", currentVersion, latestVersion),
	}
}

// upgradePlugin handles the Plugin install/update phase.
func upgradePlugin(out io.Writer) upgradeAction {
	// Step 1: Ensure marketplace is added
	if err := ensureMarketplaceAdded(out); err != nil {
		forgelog.Error("error: marketplace setup failed: %v\n", err)
		return upgradeAction{status: "FAILED", target: "Plugin", detail: err.Error()}
	}

	// Step 2: Install or update plugin
	// Use fully qualified name "forge@forge" (plugin@marketplace) because
	// "claude plugin update forge" fails with "Plugin not found" when the
	// marketplace name differs from the plugin name or resolution is ambiguous.
	const qualifiedName = "forge@forge"
	if pluginInstalledCheck() {
		_, _ = fmt.Fprintf(out, "Updating forge plugin...\n")
		if err := runClaudeCommand([]string{"plugin", "update", qualifiedName}); err != nil {
			return upgradeAction{status: "FAILED", target: "Plugin", detail: fmt.Sprintf("update failed: %v", err)}
		}
		return upgradeAction{status: "UPGRADED", target: "Plugin", detail: "updated to latest version"}
	}

	_, _ = fmt.Fprintf(out, "Installing forge plugin...\n")
	if err := runClaudeCommand([]string{"plugin", "install", qualifiedName}); err != nil {
		return upgradeAction{status: "FAILED", target: "Plugin", detail: fmt.Sprintf("install failed: %v", err)}
	}
	return upgradeAction{status: "INSTALLED", target: "Plugin", detail: "installed latest version"}
}

// ensureMarketplaceAdded checks if the forge marketplace is registered.
// If not, adds it via `claude plugin marketplace add`.
func ensureMarketplaceAdded(out io.Writer) error {
	// Check if marketplace is already added by listing marketplaces
	cmd := exec.Command("claude", "plugin", "marketplace", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If list fails, try to add anyway
		_, _ = fmt.Fprintf(out, "Adding forge marketplace...\n")
		return runClaudeCommand([]string{
			"plugin", "marketplace", "add",
			pluginMarketplaceURL,
			"--sparse", ".claude-plugin", "plugins",
		})
	}

	// Check if forge marketplace is already in the list
	if strings.Contains(string(output), "forge") {
		return nil // Already added
	}

	// Not found, add it
	_, _ = fmt.Fprintf(out, "Adding forge marketplace...\n")
	return runClaudeCommand([]string{
		"plugin", "marketplace", "add",
		pluginMarketplaceURL,
		"--sparse", ".claude-plugin", "plugins",
	})
}

// pluginInstalledImpl checks if the forge plugin is already installed.
func pluginInstalledImpl() bool {
	cmd := exec.Command("claude", "plugin", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "forge")
}

// fetchLatestReleaseImpl performs the actual HTTP GET to the GitHub Release API.
func fetchLatestReleaseImpl(url string) ([]byte, error) {
	resp, err := httpGet(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return body, nil
}

// parseVersionFromTag extracts version number from a GitHub tag.
// Expected format: "forge-cli/v{version}" or just "v{version}".
func parseVersionFromTag(tag string) (string, error) {
	// Try "forge-cli/vX.Y.Z" format first
	if after, ok := strings.CutPrefix(tag, "forge-cli/"); ok {
		return strings.TrimPrefix(after, "v"), nil
	}
	// Try "vX.Y.Z" format
	if after, ok := strings.CutPrefix(tag, "v"); ok {
		return after, nil
	}
	return "", fmt.Errorf("unexpected tag format: %q (expected 'forge-cli/v{version}')", tag)
}

// compareVersions compares two semver version strings.
// Returns -1 if a < b, 0 if a == b, 1 if a > b.
func compareVersions(a, b string) int {
	aParts := parseSemver(a)
	bParts := parseSemver(b)

	for i := 0; i < 3; i++ {
		if aParts[i] < bParts[i] {
			return -1
		}
		if aParts[i] > bParts[i] {
			return 1
		}
	}
	return 0
}

// parseSemver parses a version string "X.Y.Z" into [3]int.
func parseSemver(v string) [3]int {
	var parts [3]int
	_, _ = fmt.Sscanf(v, "%d.%d.%d", &parts[0], &parts[1], &parts[2])
	return parts
}

// buildDownloadURL constructs the binary download URL for the given version and current platform.
func buildDownloadURL(version string) string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	ext := ""
	if goos == "windows" {
		ext = ".exe"
	}

	return fmt.Sprintf(
		"https://github.com/bigfaner/forge/releases/download/forge-cli/v%s/forge-%s-%s-%s%s",
		version, version, goos, goarch, ext,
	)
}

// downloadAndReplace downloads the new binary and performs atomic replacement.
func downloadAndReplace(out io.Writer, version string) error {
	binDir := forgeBinaryDir()
	if binDir == "" {
		return fmt.Errorf("cannot determine forge binary directory")
	}

	forgePath := filepath.Join(binDir, "forge")
	if runtime.GOOS == "windows" {
		forgePath = filepath.Join(binDir, "forge.exe")
	}

	// Ensure bin directory exists
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Download new binary
	downloadURL := buildDownloadURL(version)
	_, _ = fmt.Fprintf(out, "Downloading forge v%s...\n", version)

	resp, err := httpGet(downloadURL)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	// Read body into memory and verify completeness against Content-Length
	// to catch incomplete downloads on weak networks.
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("download incomplete: %w", err)
	}
	if resp.ContentLength > 0 && int64(len(content)) != resp.ContentLength {
		return fmt.Errorf("download incomplete: got %d bytes, expected %d", len(content), resp.ContentLength)
	}

	if runtime.GOOS == "windows" {
		return replaceWindowsBinary(out, forgePath, content)
	}

	return replaceUnixBinary(out, forgePath, content)
}

// replaceUnixBinary atomically replaces the forge binary on Unix systems.
// content is the complete binary data already verified by the caller.
// If atomic rename fails, falls back to removing the target then retrying,
// then to a direct write as a last resort.
func replaceUnixBinary(out io.Writer, forgePath string, content []byte) error {
	newPath := forgePath + ".new"

	// Clean up stale .new file from a previous failed attempt
	_ = os.Remove(newPath)

	// Write new binary to temp file
	if err := os.WriteFile(newPath, content, 0o755); err != nil {
		return fmt.Errorf("failed to write new binary: %w", err)
	}

	// Step 1: Atomic rename
	if err := osRename(newPath, forgePath); err == nil {
		_, _ = fmt.Fprintf(out, "Binary replaced successfully.\n")
		return nil
	}

	// Step 2: Remove target and retry (handles locked-running-binary edge case)
	_ = os.Remove(forgePath)
	if err := osRename(newPath, forgePath); err == nil {
		_, _ = fmt.Fprintf(out, "Binary replaced successfully.\n")
		return nil
	}

	// Step 3: Direct write as last resort
	_ = os.Remove(newPath)
	if err := os.WriteFile(forgePath, content, 0o755); err != nil {
		return fmt.Errorf("failed to replace binary: direct write failed: %w", err)
	}

	_, _ = fmt.Fprintf(out, "Binary replaced successfully.\n")
	return nil
}

// replaceWindowsBinary handles the Windows rename dance:
// forge.exe -> forge.old, write new forge.exe, delete forge.old.
func replaceWindowsBinary(out io.Writer, forgePath string, content []byte) error {
	oldPath := forgePath + ".old"

	// Step 1: Rename current binary to .old (best effort — may not exist on fresh install)
	if _, err := os.Stat(forgePath); err == nil {
		// Remove stale .old if it exists
		_ = os.Remove(oldPath)
		if err := os.Rename(forgePath, oldPath); err != nil {
			return fmt.Errorf("failed to rename current binary to .old: %w", err)
		}
	}

	// Step 2: Write new binary
	if err := os.WriteFile(forgePath, content, 0o755); err != nil {
		return fmt.Errorf("failed to write new binary: %w", err)
	}

	// Step 3: Delete .old (best effort — file may be locked)
	_ = os.Remove(oldPath)

	_, _ = fmt.Fprintf(out, "Binary replaced successfully.\n")
	return nil
}

// printUpgradeSummary prints the upgrade results in a structured block.
func printUpgradeSummary(out io.Writer, actions []upgradeAction) {
	_, _ = fmt.Fprintf(out, ">>>\n")
	for _, a := range actions {
		detail := ""
		if a.detail != "" {
			detail = fmt.Sprintf(" (%s)", a.detail)
		}
		_, _ = fmt.Fprintf(out, "%-10s %s%s\n", a.status, a.target, detail)
	}
	_, _ = fmt.Fprintf(out, "<<<\n")
}
