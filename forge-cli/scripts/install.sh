#!/bin/bash
# install.sh — Download and install forge CLI from GitHub Releases
# Usage: curl -fsSL https://github.com/bigfaner/forge/releases/latest/download/install.sh | bash
#
# Reuses platform detection, PATH management, and atomic replace patterns
# from install-local.sh. This script differs by downloading a pre-compiled
# binary from GitHub Releases instead of building locally.

set -e

# Configuration
APP_NAME="forge"
INSTALL_DIR="${HOME}/.forge/bin"
GITHUB_REPO="bigfaner/forge"

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

# Fetch latest version from GitHub Release API
get_latest_version() {
    log_info "Fetching latest version from GitHub..."

    VERSION=$(curl -fsSL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" \
        | grep '"tag_name"' \
        | head -1 \
        | sed -E 's/.*"forge-cli\/v([^"]+)".*/\1/')

    if [ -z "$VERSION" ]; then
        log_error "Failed to determine latest version from GitHub API"
        exit 1
    fi

    log_info "Latest version: ${VERSION}"
}

# Download and install the binary
download_and_install() {
    local BINARY_NAME="forge-${VERSION}-${OS}-${ARCH}"
    local DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/forge-cli/v${VERSION}/${BINARY_NAME}"
    local TEMP_FILE="${INSTALL_DIR}/${APP_NAME}.new"

    log_info "Downloading ${BINARY_NAME}..."

    mkdir -p "${INSTALL_DIR}"

    curl -fsSL --progress-bar -o "${TEMP_FILE}" "${DOWNLOAD_URL}"

    if [ ! -f "${TEMP_FILE}" ]; then
        log_error "Download failed: ${DOWNLOAD_URL}"
        exit 1
    fi

    # Atomic replacement: chmod +x then mv (atomic on same filesystem)
    chmod +x "${TEMP_FILE}"
    mv -f "${TEMP_FILE}" "${INSTALL_DIR}/${APP_NAME}"

    log_info "Installed forge v${VERSION} to ${INSTALL_DIR}/${APP_NAME}"
}

# Add to PATH in shell RC files
add_to_path() {
    local path_export="export PATH=\"\${PATH}:${INSTALL_DIR}\""
    local added=false

    # Check if already in PATH
    if [[ ":$PATH:" == *":${INSTALL_DIR}:"* ]]; then
        log_info "${INSTALL_DIR} is already in PATH"
        return 0
    fi

    # Update shell RC files (.bashrc, .zshrc, .profile)
    for rc_file in "${HOME}/.bashrc" "${HOME}/.zshrc" "${HOME}/.profile"; do
        # Skip files that already have the path configured
        if [ -f "$rc_file" ] && grep -q "${INSTALL_DIR}" "$rc_file" 2>/dev/null; then
            log_info "${INSTALL_DIR} is already configured in ${rc_file}"
            continue
        fi

        # Only add to files that exist, or to .profile as fallback
        if [ -f "$rc_file" ] || [ "$rc_file" = "${HOME}/.profile" ]; then
            log_info "Adding ${INSTALL_DIR} to PATH in ${rc_file}..."
            echo "" >> "$rc_file"
            echo "# Added by forge-cli installer" >> "$rc_file"
            echo "$path_export" >> "$rc_file"
            added=true
        fi
    done

    if [ "$added" = true ]; then
        log_warn "Please run 'source ~/.bashrc' (or ~/.zshrc / ~/.profile) or restart your terminal to update PATH"
    fi
}

# Print verification instructions
print_verify_instructions() {
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  forge v${VERSION} installed successfully!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo "To verify the installation:"
    echo ""
    echo "  source ~/.bashrc  # or ~/.zshrc / ~/.profile"
    echo "  forge version"
    echo ""
    echo "Next steps:"
    echo ""
    echo "  forge upgrade    # Install or update the forge Plugin"
    echo "  cd my-project && forge init  # Initialize forge in a project"
    echo ""
}

# Main
main() {
    log_info "Starting forge CLI installation..."

    detect_platform
    get_latest_version
    download_and_install
    add_to_path
    print_verify_instructions
}

main "$@"
