package scripts_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// scriptDir is the directory containing the install scripts.
func scriptDir(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func readScript(t *testing.T, name string) string {
	t.Helper()
	path := filepath.Join(scriptDir(t), name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", name, err)
	}
	return string(data)
}

// --- install.sh tests ---

func TestInstallSh_DetectsOS(t *testing.T) {
	content := readScript(t, "install.sh")
	// Must detect darwin and linux
	if !strings.Contains(content, "darwin") {
		t.Error("install.sh does not handle darwin OS")
	}
	if !strings.Contains(content, "linux") {
		t.Error("install.sh does not handle linux OS")
	}
}

func TestInstallSh_DetectsArchitecture(t *testing.T) {
	content := readScript(t, "install.sh")
	// Must detect amd64 and arm64
	if !strings.Contains(content, "amd64") {
		t.Error("install.sh does not handle amd64 architecture")
	}
	if !strings.Contains(content, "arm64") {
		t.Error("install.sh does not handle arm64 architecture")
	}
	// x86_64 must map to amd64
	if !strings.Contains(content, "x86_64|amd64) ARCH=\"amd64\"") {
		t.Error("install.sh does not map x86_64 to amd64")
	}
	// aarch64 must map to arm64
	if !strings.Contains(content, "arm64|aarch64) ARCH=\"arm64\"") {
		t.Error("install.sh does not map aarch64 to arm64")
	}
}

func TestInstallSh_FetchesLatestVersionFromGitHubAPI(t *testing.T) {
	content := readScript(t, "install.sh")
	// Must use GitHub API to fetch latest version
	if !strings.Contains(content, "api.github.com/repos/") {
		t.Error("install.sh does not use GitHub API to fetch latest version")
	}
	if !strings.Contains(content, "releases/latest") {
		t.Error("install.sh does not fetch latest release")
	}
	if !strings.Contains(content, `"tag_name"`) {
		t.Error("install.sh does not parse tag_name from API response")
	}
}

func TestInstallSh_ConstructsCorrectDownloadURL(t *testing.T) {
	content := readScript(t, "install.sh")
	// Tag must use v prefix: forge-cli/v{VERSION}
	if !strings.Contains(content, "forge-cli/v${VERSION}") {
		t.Error("install.sh tag does not use v prefix in download URL")
	}
	// Binary name must NOT use v prefix: forge-{VERSION}-{OS}-{ARCH}
	if !strings.Contains(content, "forge-${VERSION}-${OS}-${ARCH}") {
		t.Error("install.sh binary name format is incorrect")
	}
}

func TestInstallSh_HardRule_TagVPrefix_BinaryNoVPrefix(t *testing.T) {
	content := readScript(t, "install.sh")
	// Hard Rule: tag uses v prefix, binary filename does NOT use v prefix
	// Verify: the sed/regex extraction strips v prefix from tag
	if !regexp.MustCompile(`sed.*v`).MatchString(content) {
		t.Error("install.sh does not strip v prefix from tag to get version number")
	}
	// Verify: BINARY_NAME does NOT contain v prefix before version
	if strings.Contains(content, "forge-v${VERSION}") {
		t.Error("install.sh binary name incorrectly uses v prefix (Hard Rule violation)")
	}
}

func TestInstallSh_DownloadsToForgeNew(t *testing.T) {
	content := readScript(t, "install.sh")
	// Must download to .new temp file (uses ${APP_NAME}.new where APP_NAME="forge")
	if !strings.Contains(content, "${APP_NAME}.new") {
		t.Error("install.sh does not download to ${APP_NAME}.new temp file")
	}
}

func TestInstallSh_AtomicReplace(t *testing.T) {
	content := readScript(t, "install.sh")
	// Must use mv -f for atomic replacement
	if !strings.Contains(content, "mv -f") {
		t.Error("install.sh does not use atomic mv -f for replacement")
	}
}

func TestInstallSh_AddsToPathInRCFiles(t *testing.T) {
	content := readScript(t, "install.sh")
	// Must mention .bashrc, .zshrc, .profile
	for _, rc := range []string{".bashrc", ".zshrc", ".profile"} {
		if !strings.Contains(content, rc) {
			t.Errorf("install.sh does not handle %s", rc)
		}
	}
	// Must export PATH
	if !strings.Contains(content, "export PATH") {
		t.Error("install.sh does not export PATH")
	}
}

func TestInstallSh_PrintsVerificationInstructions(t *testing.T) {
	content := readScript(t, "install.sh")
	// Must print verification instructions
	if !strings.Contains(content, "forge --version") {
		t.Error("install.sh does not print forge --version verification instruction")
	}
}

func TestInstallSh_CurlDownload(t *testing.T) {
	content := readScript(t, "install.sh")
	// Must use curl to download
	if !strings.Contains(content, "curl -fsSL") {
		t.Error("install.sh does not use curl -fsSL for downloading")
	}
}

func TestInstallSh_ChmodBeforeMove(t *testing.T) {
	content := readScript(t, "install.sh")
	// Must chmod +x before atomic move
	if !strings.Contains(content, "chmod +x") {
		t.Error("install.sh does not chmod +x the downloaded binary")
	}
}

// --- install.ps1 tests ---

func TestInstallPs1_DetectsWindowsArchitecture(t *testing.T) {
	content := readScript(t, "install.ps1")
	// Must handle AMD64 and ARM64
	if !strings.Contains(content, `"AMD64"`) {
		t.Error("install.ps1 does not handle AMD64 architecture")
	}
	if !strings.Contains(content, `"ARM64"`) {
		t.Error("install.ps1 does not handle ARM64 architecture")
	}
}

func TestInstallPs1_FetchesLatestVersionFromGitHubAPI(t *testing.T) {
	content := readScript(t, "install.ps1")
	// Must use GitHub API
	if !strings.Contains(content, "api.github.com/repos/") {
		t.Error("install.ps1 does not use GitHub API")
	}
	if !strings.Contains(content, "Invoke-RestMethod") {
		t.Error("install.ps1 does not use Invoke-RestMethod for API call")
	}
	if !strings.Contains(content, "tag_name") {
		t.Error("install.ps1 does not parse tag_name from API response")
	}
}

func TestInstallPs1_ConstructsCorrectDownloadURL(t *testing.T) {
	content := readScript(t, "install.ps1")
	// Tag uses v prefix: forge-cli/v
	if !strings.Contains(content, "forge-cli/v") {
		t.Error("install.ps1 tag does not use v prefix in download URL")
	}
	// Binary name does NOT use v prefix: forge-{version}-windows-{arch}.exe
	if !strings.Contains(content, "forge-$") {
		t.Error("install.ps1 binary name format may be incorrect")
	}
}

func TestInstallPs1_HardRule_TagVPrefix_BinaryNoVPrefix(t *testing.T) {
	content := readScript(t, "install.ps1")
	// Hard Rule: tag uses v prefix, binary filename does NOT
	if !strings.Contains(content, "forge-cli/v$") || !strings.Contains(content, `forge-cli/v`) {
		t.Error("install.ps1 does not use v prefix in tag")
	}
	// Binary name should NOT have v prefix
	if strings.Contains(content, "forge-v$") {
		t.Error("install.ps1 binary name incorrectly uses v prefix (Hard Rule violation)")
	}
}

func TestInstallPs1_DownloadsToUserProfileForgeBin(t *testing.T) {
	content := readScript(t, "install.ps1")
	// Must install to %USERPROFILE%\.forge\bin
	if !strings.Contains(content, "$env:USERPROFILE") {
		t.Error("install.ps1 does not use USERPROFILE for install directory")
	}
	if !strings.Contains(content, ".forge\\bin") {
		t.Error("install.ps1 does not install to .forge\\bin")
	}
}

func TestInstallPs1_UpdatesUserPATH(t *testing.T) {
	content := readScript(t, "install.ps1")
	// Must use [Environment]::SetEnvironmentVariable for persistent PATH
	if !strings.Contains(content, "[Environment]::SetEnvironmentVariable") {
		t.Error("install.ps1 does not use [Environment]::SetEnvironmentVariable for PATH")
	}
	if !strings.Contains(content, `"User"`) {
		t.Error("install.ps1 does not set User-level PATH")
	}
}

func TestInstallPs1_AtomicReplace(t *testing.T) {
	content := readScript(t, "install.ps1")
	// Must use Move-Item -Force for atomic replacement
	if !strings.Contains(content, "Move-Item") {
		t.Error("install.ps1 does not use Move-Item for atomic replacement")
	}
	// Must download to .new temp file first
	if !strings.Contains(content, ".new") {
		t.Error("install.ps1 does not use temp .new file for atomic replacement")
	}
}

func TestInstallPs1_NoRunningExeRenameDance(t *testing.T) {
	content := readScript(t, "install.ps1")
	// Per task spec: install.ps1 does NOT need running exe rename dance
	if strings.Contains(content, ".old") {
		t.Error("install.ps1 should NOT contain .old rename dance (only needed in forge upgrade)")
	}
}

func TestInstallPs1_PrintsVerificationInstructions(t *testing.T) {
	content := readScript(t, "install.ps1")
	if !strings.Contains(content, "forge --version") {
		t.Error("install.ps1 does not print forge --version verification instruction")
	}
}

// --- Cross-cutting tests ---

func TestBothScripts_UseSameGitHubRepo(t *testing.T) {
	sh := readScript(t, "install.sh")
	ps1 := readScript(t, "install.ps1")
	// Both must reference the same GitHub repo
	if !strings.Contains(sh, "bigfaner/forge") {
		t.Error("install.sh does not reference bigfaner/forge")
	}
	if !strings.Contains(ps1, "bigfaner/forge") {
		t.Error("install.ps1 does not reference bigfaner/forge")
	}
}

func TestBothScripts_InstallToForgeBin(t *testing.T) {
	sh := readScript(t, "install.sh")
	ps1 := readScript(t, "install.ps1")
	if !strings.Contains(sh, ".forge/bin") {
		t.Error("install.sh does not install to .forge/bin")
	}
	if !strings.Contains(ps1, ".forge") {
		t.Error("install.ps1 does not install to .forge")
	}
}
