#!/usr/bin/env bash
# download-just.sh — Download just binaries for all 6 target platforms.
#
# Usage:
#   bash scripts/download-just.sh [VERSION]
#
# VERSION defaults to 1.51.0. Binaries are saved to
# internal/embedded/just/binaries/ with the naming convention expected by
# the build-tagged embed files:
#   just-linux-amd64, just-linux-arm64,
#   just-darwin-amd64, just-darwin-arm64,
#   just-windows-amd64.exe, just-windows-arm64.exe
#
# After downloading, this script verifies SHA256 checksums against the
# upstream SHA256SUMS file.

set -euo pipefail

VERSION="${1:-1.51.0}"
BASE_URL="https://github.com/casey/just/releases/download/${VERSION}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARIES_DIR="${SCRIPT_DIR}/../internal/embedded/just/binaries"

# Cross-platform sha256 (macOS has shasum, Linux has sha256sum)
if command -v sha256sum >/dev/null 2>&1; then
  sha256_of() { sha256sum "$1" | awk '{print $1}'; }
else
  sha256_of() { shasum -a 256 "$1" | awk '{print $1}'; }
fi

# Platform definitions: local_name|github_asset
PLATFORMS="
just-darwin-arm64|just-${VERSION}-aarch64-apple-darwin.tar.gz
just-darwin-amd64|just-${VERSION}-x86_64-apple-darwin.tar.gz
just-linux-arm64|just-${VERSION}-aarch64-unknown-linux-musl.tar.gz
just-linux-amd64|just-${VERSION}-x86_64-unknown-linux-musl.tar.gz
just-windows-arm64.exe|just-${VERSION}-aarch64-pc-windows-msvc.zip
just-windows-amd64.exe|just-${VERSION}-x86_64-pc-windows-msvc.zip
"

mkdir -p "${BINARIES_DIR}"

# ---- Download SHA256SUMS first ----
echo "Downloading SHA256SUMS for just v${VERSION}..."
checksums_path="${BINARIES_DIR}/SHA256SUMS"
curl -fSL -o "${checksums_path}" "${BASE_URL}/SHA256SUMS"

# ---- Download, verify, and extract each platform ----
failed=0
echo "${PLATFORMS}" | while IFS='|' read -r local_name asset; do
  # Skip blank lines
  [ -z "${local_name}" ] && continue

  url="${BASE_URL}/${asset}"
  dest="${BINARIES_DIR}/${local_name}"
  archive_path="${BINARIES_DIR}/${asset}"

  echo ""
  echo "Downloading ${asset}..."
  curl -fSL -o "${archive_path}" "${url}"

  # Verify archive checksum before extraction
  expected_hash=$(grep -F "$(basename "${asset}")" "${checksums_path}" | awk '{print $1}' || true)
  if [ -n "${expected_hash}" ]; then
    actual_hash=$(sha256_of "${archive_path}")
    if [ "${actual_hash}" != "${expected_hash}" ]; then
      echo "  FAIL: checksum mismatch for ${asset}"
      echo "    expected: ${expected_hash}"
      echo "    actual:   ${actual_hash}"
      rm -f "${archive_path}"
      # Signal failure via exit code (subshell)
      exit 1
    fi
    echo "  OK: checksum verified"
  else
    echo "  WARNING: no checksum found for ${asset}, skipping verification"
  fi

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

  size=$(wc -c < "${dest}" | tr -d ' ')
  echo "  -> ${dest} (${size} bytes)"

  # Clean up archive
  rm -f "${archive_path}"
done || failed=1

rm -f "${checksums_path}"

if [ "${failed}" -ne 0 ]; then
  echo ""
  echo "ERROR: download or checksum verification failed"
  exit 1
fi

echo ""
echo "All 6 platform binaries (just v${VERSION}) downloaded successfully."
