#!/bin/bash
# Install script for Unix-like systems (Linux, macOS)
# Builds and installs the task CLI to ~/.zcode-task-cli/

set -e

# Configuration
APP_NAME="task"
INSTALL_DIR="${HOME}/.zcode-task-cli"
BIN_DIR="bin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS and Architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$OS" in
        darwin) OS="darwin" ;;
        linux)  OS="linux" ;;
        *)      log_error "Unsupported OS: $OS"; exit 1 ;;
    esac

    case "$ARCH" in
        x86_64|amd64) ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        *)             log_error "Unsupported architecture: $ARCH"; exit 1 ;;
    esac

    log_info "Detected platform: ${OS}/${ARCH}"
}

# Read version from version.txt
get_version() {
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    VERSION_FILE="${SCRIPT_DIR}/version.txt"
    if [ -f "$VERSION_FILE" ]; then
        VERSION=$(tr -d '[:space:]' < "$VERSION_FILE")
    else
        VERSION="dev"
    fi
    log_info "Version: ${VERSION}"
}

# Build the executable
build() {
    log_info "Building ${APP_NAME} for ${OS}/${ARCH}..."

    # Get script directory and project root
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

    cd "$PROJECT_ROOT"

    # Create bin directory
    mkdir -p "${BIN_DIR}"

    # Build with version injection
    OUTPUT="${BIN_DIR}/${APP_NAME}"
    if [ "$OS" = "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
    fi

    LDFLAGS="-s -w -X task-cli/pkg/version.Version=${VERSION}"
    CGO_ENABLED=0 GOOS="$OS" GOARCH="$ARCH" go build -ldflags="$LDFLAGS" -o "$OUTPUT" ./cmd/task

    log_info "Build complete: ${OUTPUT}"
}

# Install to user directory
install() {
    log_info "Installing to ${INSTALL_DIR}..."

    # Create installation directory
    mkdir -p "${INSTALL_DIR}"

    # Copy binary
    cp "${BIN_DIR}/${APP_NAME}" "${INSTALL_DIR}/${APP_NAME}"
    chmod +x "${INSTALL_DIR}/${APP_NAME}"

    log_info "Installation complete: ${INSTALL_DIR}/${APP_NAME}"
}

# Add to PATH
add_to_path() {
    local shell_rc=""
    local path_export="export PATH=\"\${PATH}:${INSTALL_DIR}\""

    # Detect shell configuration file
    if [ -n "$ZSH_VERSION" ]; then
        shell_rc="${HOME}/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        shell_rc="${HOME}/.bashrc"
    else
        shell_rc="${HOME}/.profile"
    fi

    # Check if already in PATH
    if [[ ":$PATH:" == *":${INSTALL_DIR}:"* ]]; then
        log_info "${INSTALL_DIR} is already in PATH"
        return 0
    fi

    # Check if already in shell config
    if [ -f "$shell_rc" ] && grep -q "${INSTALL_DIR}" "$shell_rc" 2>/dev/null; then
        log_info "${INSTALL_DIR} is already configured in ${shell_rc}"
        return 0
    fi

    log_info "Adding ${INSTALL_DIR} to PATH in ${shell_rc}..."
    echo "" >> "$shell_rc"
    echo "# Added by task-cli installer" >> "$shell_rc"
    echo "$path_export" >> "$shell_rc"

    log_warn "Please run 'source ${shell_rc}' or restart your terminal to update PATH"
}

# Main
main() {
    log_info "Starting local installation..."
    get_version
    detect_platform
    build
    install
    add_to_path
    log_info "Done! Run 'task --help' to verify installation."
}

main "$@"
