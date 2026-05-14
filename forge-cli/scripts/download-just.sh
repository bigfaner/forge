#!/usr/bin/env bash
# download-just.sh — Download just v1.40.0 binaries for all 6 target platforms.
#
# Usage:
#   bash scripts/download-just.sh
#
# Binaries are saved to internal/embedded/just/binaries/ with the naming
# convention expected by the build-tagged embed files:
#   just-linux-amd64, just-linux-arm64,
#   just-darwin-amd64, just-darwin-arm64,
#   just-windows-amd64.exe, just-windows-arm64.exe

set -euo pipefail

VERSION="1.40.0"
BASE_URL="https://github.com/casey/just/releases/download/${VERSION}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARIES_DIR="${SCRIPT_DIR}/../internal/embedded/just/binaries"

# Mapping: platform_file -> github_asset_name
declare -A PLATFORMS=(
  ["just-linux-amd64"]="just-${VERSION}-x86_64-unknown-linux-musl.tar.gz"
  ["just-linux-arm64"]="just-${VERSION}-aarch64-unknown-linux-musl.tar.gz"
  ["just-darwin-amd64"]="just-${VERSION}-x86_64-apple-darwin.tar.gz"
  ["just-darwin-arm64"]="just-${VERSION}-aarch64-apple-darwin.tar.gz"
  ["just-windows-amd64.exe"]="just-${VERSION}-x86_64-pc-windows-msvc.zip"
  ["just-windows-arm64.exe"]="just-${VERSION}-aarch64-pc-windows-msvc.zip"
)

mkdir -p "${BINARIES_DIR}"

for platform_file in "${!PLATFORMS[@]}"; do
  asset="${PLATFORMS[${platform_file}]}"
  url="${BASE_URL}/${asset}"
  dest="${BINARIES_DIR}/${platform_file}"
  tmpdir=""

  echo "Downloading ${asset}..."

  # Download archive
  archive_path="${BINARIES_DIR}/${asset}"
  curl -fSL -o "${archive_path}" "${url}"

  # Extract binary from archive
  tmpdir=$(mktemp -d)
  case "${asset}" in
    *.tar.gz)
      tar xzf "${archive_path}" -C "${tmpdir}" just
      mv "${tmpdir}/just" "${dest}"
      ;;
    *.zip)
      unzip -o -j "${archive_path}" "just.exe" -d "${tmpdir}"
      mv "${tmpdir}/just.exe" "${dest}"
      ;;
  esac

  rm -rf "${tmpdir}"
  chmod +x "${dest}"

  size=$(wc -c < "${dest}")
  echo "  -> ${dest} (${size} bytes)"
done

# Download and verify checksums
echo ""
echo "Downloading SHA256SUMS for verification..."
checksums_path="${BINARIES_DIR}/SHA256SUMS"
curl -fSL -o "${checksums_path}" "${BASE_URL}/SHA256SUMS"

echo ""
echo "Verifying checksums..."
failed=0
for platform_file in "${!PLATFORMS[@]}"; do
  asset="${PLATFORMS[${platform_file}]}"
  dest="${BINARIES_DIR}/${platform_file}"

  # Get expected hash from SHA256SUMS file for the just binary inside the archive
  # The SHA256SUMS file lists hashes for files inside the archives
  expected_hash=$(grep -E "$(basename "${asset}")" "${checksums_path}" | awk '{print $1}' || true)

  if [ -z "${expected_hash}" ]; then
    echo "  WARNING: no checksum found for ${asset}"
    continue
  fi

  # Verify archive hash
  archive_path="${BINARIES_DIR}/${asset}"
  actual_hash=$(sha256sum "${archive_path}" | awk '{print $1}')

  if [ "${actual_hash}" = "${expected_hash}" ]; then
    echo "  OK: ${asset} checksum verified"
  else
    echo "  FAIL: ${asset} checksum mismatch"
    echo "    expected: ${expected_hash}"
    echo "    actual:   ${actual_hash}"
    failed=1
  fi
done

# Clean up downloaded archives
echo ""
echo "Cleaning up archives..."
for platform_file in "${!PLATFORMS[@]}"; do
  asset="${PLATFORMS[${platform_file}]}"
  rm -f "${BINARIES_DIR}/${asset}"
done
rm -f "${checksums_path}"

if [ "${failed}" -ne 0 ]; then
  echo ""
  echo "ERROR: checksum verification failed for one or more binaries"
  exit 1
fi

echo ""
echo "All ${#PLATFORMS[@]} platform binaries downloaded successfully."
