# Install script for Windows
# Builds and installs the task CLI to ~/.zcode-task-cli/

param(
    [switch]$Force
)

$ErrorActionPreference = "Stop"

# Configuration
$AppName = "task.exe"
$InstallDir = "$env:USERPROFILE\.zcode-task-cli"
$BinDir = "bin"

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

# Read version from version.txt
function Get-Version {
    $ScriptDir = Split-Path -Parent $MyInvocation.ScriptName
    if (-not $ScriptDir) {
        $ScriptDir = $PSScriptRoot
    }
    $VersionFile = Join-Path $ScriptDir "version.txt"

    if (Test-Path $VersionFile) {
        $script:Version = (Get-Content $VersionFile -Raw).Trim()
    } else {
        $script:Version = "dev"
    }
    Write-Info "Version: $script:Version"
}

# Build the executable
function Build-App {
    param([string]$Arch)

    Write-Info "Building $AppName for windows/$Arch..."

    $ScriptDir = Split-Path -Parent $MyInvocation.ScriptName
    if (-not $ScriptDir) {
        $ScriptDir = $PSScriptRoot
    }
    $ProjectRoot = Split-Path -Parent $ScriptDir

    Push-Location $ProjectRoot

    try {
        # Create bin directory
        if (-not (Test-Path $BinDir)) {
            New-Item -ItemType Directory -Path $BinDir | Out-Null
        }

        $Output = Join-Path $BinDir $AppName

        $env:CGO_ENABLED = "0"
        $env:GOOS = "windows"
        $env:GOARCH = $Arch

        $LdFlags = "-s -w -X task-cli/pkg/version.Version=$script:Version"
        go build -ldflags="$LdFlags" -o $Output ./cmd/task

        Write-Info "Build complete: $Output"
    }
    finally {
        Pop-Location
    }
}

# Install to user directory
function Install-App {
    Write-Info "Installing to $InstallDir..."

    # Create installation directory
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir | Out-Null
    }

    # Copy binary
    $SourcePath = Join-Path $BinDir $AppName
    $DestPath = Join-Path $InstallDir $AppName

    Copy-Item -Path $SourcePath -Destination $DestPath -Force

    Write-Info "Installation complete: $DestPath"
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

# Main
function Main {
    Write-Info "Starting local installation..."

    Get-Version
    $Arch = Get-Platform
    Build-App -Arch $Arch
    Install-App
    Add-ToPath

    Write-Info "Done! Run 'task --help' to verify installation."
}

Main
