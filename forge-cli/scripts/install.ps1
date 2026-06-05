# install.ps1 — Download and install forge CLI from GitHub Releases (Windows)
# Usage: irm https://github.com/bigfaner/forge/releases/latest/download/install.ps1 | iex
#
# Reuses platform detection, PATH management, and atomic replace patterns
# from install-local.ps1. This script differs by downloading a pre-compiled
# binary from GitHub Releases instead of building locally.

param(
    [switch]$Force
)

$ErrorActionPreference = "Stop"

# Configuration
$AppName = "forge.exe"
$InstallDir = "$env:USERPROFILE\.forge\bin"
$GitHubRepo = "bigfaner/forge"

function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] " -ForegroundColor Green -NoNewline
    Write-Host $Message
}

function Write-Warn {
    param([string]$Message)
    Write-Host "[WARN] " -ForegroundColor Yellow -NoNewline
    Write-Host $Message
}

function Write-ErrorMsg {
    param([string]$Message)
    Write-Host "[ERROR] " -ForegroundColor Red -NoNewline
    Write-Host $Message
}

# Detect Architecture
function Get-Platform {
    $Arch = $env:PROCESSOR_ARCHITECTURE

    switch ($Arch) {
        "AMD64" { $GoArch = "amd64" }
        "ARM64" { $GoArch = "arm64" }
        default {
            Write-ErrorMsg "Unsupported architecture: $Arch"
            exit 1
        }
    }

    Write-Info "Detected platform: windows/$GoArch"
    return $GoArch
}

# Fetch latest version from GitHub Release API
function Get-LatestVersion {
    Write-Info "Fetching latest version from GitHub..."

    $ApiUrl = "https://api.github.com/repos/$GitHubRepo/releases/latest"
    $Response = Invoke-RestMethod -Uri $ApiUrl -Method Get

    $Tag = $Response.tag_name
    # Tag format: forge-cli/v5.17.0 → extract version "5.17.0"
    if ($Tag -match "forge-cli/v(.+)") {
        $script:Version = $Matches[1]
    } else {
        Write-ErrorMsg "Unexpected tag format: $Tag"
        exit 1
    }

    Write-Info "Latest version: $script:Version"
}

# Download and install the binary
function Install-FromRelease {
    param([string]$Arch)

    $BinaryName = "forge-$script:Version-windows-$Arch.exe"
    $DownloadUrl = "https://github.com/$GitHubRepo/releases/download/forge-cli/v$script:Version/$BinaryName"

    Write-Info "Downloading $BinaryName..."

    # Create installation directory
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir | Out-Null
    }

    $TempPath = Join-Path $InstallDir "$AppName.new"
    $DestPath = Join-Path $InstallDir $AppName

    # Download to temp file
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempPath -UseBasicParsing

    if (-not (Test-Path $TempPath)) {
        Write-ErrorMsg "Download failed: $DownloadUrl"
        exit 1
    }

    # Atomic replacement: copy to temp file then rename
    # MoveTo is atomic on NTFS, avoids race with hooks reading the binary
    Move-Item -Path $TempPath -Destination $DestPath -Force

    Write-Info "Installed forge v$script:Version to $DestPath"
}

# Add to PATH
function Add-ToPath {
    # Check if already in current session PATH
    $PathParts = $env:PATH -split ";"
    if ($PathParts -contains $InstallDir) {
        Write-Info "$InstallDir is already in current PATH"
        return
    }

    # Check if already in user environment variables
    $UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    $UserPathParts = $UserPath -split ";"

    if ($UserPathParts -contains $InstallDir) {
        Write-Info "$InstallDir is already in user PATH environment variable"
        # Add to current session
        $env:PATH = "$env:PATH;$InstallDir"
        return
    }

    Write-Info "Adding $InstallDir to user PATH..."

    # Add to user environment variable (persistent)
    if ([string]::IsNullOrEmpty($UserPath)) {
        [Environment]::SetEnvironmentVariable("PATH", $InstallDir, "User")
    } else {
        [Environment]::SetEnvironmentVariable("PATH", "$UserPath;$InstallDir", "User")
    }

    # Add to current session
    $env:PATH = "$env:PATH;$InstallDir"

    Write-Warn "PATH has been updated. You may need to restart your terminal for changes to take effect in new sessions."
}

# Print verification instructions
function Write-VerifyInstructions {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "  forge v$script:Version installed successfully!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "To verify the installation:"
    Write-Host ""
    Write-Host "  forge --version"
    Write-Host ""
    Write-Host "Next steps:"
    Write-Host ""
    Write-Host "  forge upgrade    # Install or update the forge Plugin"
    Write-Host "  cd my-project; forge init  # Initialize forge in a project"
    Write-Host ""
}

# Main
function Main {
    Write-Info "Starting forge CLI installation..."

    $Arch = Get-Platform
    Get-LatestVersion
    Install-FromRelease -Arch $Arch
    Add-ToPath
    Write-VerifyInstructions
}

Main
