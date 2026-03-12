#!/usr/bin/env bash
#
# install.sh — Download and install pre-built adbclaw binary.
#
# Usage:
#   curl -fsSL https://github.com/llm-net/adbclaw/releases/latest/download/install.sh | bash
#   curl -fsSL .../install.sh | INSTALL_DIR=/custom/path bash
#
# Environment variables:
#   INSTALL_DIR  — where to put the binary (default: ~/.local/bin)
#   VERSION      — specific version to install (default: latest)
#

set -euo pipefail

REPO="llm-net/adbclaw"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BOLD='\033[1m'
NC='\033[0m'

info()  { echo -e "${GREEN}[adbclaw]${NC} $*"; }
warn()  { echo -e "${YELLOW}[adbclaw]${NC} $*"; }
error() { echo -e "${RED}[adbclaw]${NC} $*" >&2; }

# --- Detect OS and Arch ---
detect_platform() {
    local os arch
    os="$(uname -s | tr '[:upper:]' '[:lower:]')"
    arch="$(uname -m)"

    case "$os" in
        darwin) os="darwin" ;;
        linux)  os="linux" ;;
        *)      error "Unsupported OS: $os (adbclaw supports darwin and linux)"; exit 1 ;;
    esac

    case "$arch" in
        x86_64|amd64)  arch="amd64" ;;
        aarch64|arm64) arch="arm64" ;;
        *)             error "Unsupported architecture: $arch (adbclaw supports amd64 and arm64)"; exit 1 ;;
    esac

    echo "${os}-${arch}"
}

# --- Get latest release tag ---
get_latest_version() {
    local url="https://api.github.com/repos/${REPO}/releases/latest"
    local tag

    if command -v curl &>/dev/null; then
        tag="$(curl -fsSL "$url" 2>/dev/null | grep '"tag_name"' | sed 's/.*"tag_name":[[:space:]]*"\([^"]*\)".*/\1/')"
    elif command -v wget &>/dev/null; then
        tag="$(wget -qO- "$url" 2>/dev/null | grep '"tag_name"' | sed 's/.*"tag_name":[[:space:]]*"\([^"]*\)".*/\1/')"
    else
        error "Neither curl nor wget found."
        return 1
    fi

    if [ -z "$tag" ]; then
        error "Failed to get latest release version from GitHub."
        return 1
    fi

    echo "$tag"
}

# --- Verify checksum ---
verify_checksum() {
    local binary="$1"
    local version="$2"
    local platform="$3"
    local checksums_url="https://github.com/${REPO}/releases/download/${version}/checksums.txt"
    local expected actual

    info "Verifying checksum..."

    local checksums
    if command -v curl &>/dev/null; then
        checksums="$(curl -fsSL "$checksums_url" 2>/dev/null)" || { warn "Checksum file not available, skipping verification."; return 0; }
    else
        checksums="$(wget -qO- "$checksums_url" 2>/dev/null)" || { warn "Checksum file not available, skipping verification."; return 0; }
    fi

    expected="$(echo "$checksums" | grep "adbclaw-${platform}" | awk '{print $1}')"
    if [ -z "$expected" ]; then
        warn "No checksum found for adbclaw-${platform}, skipping verification."
        return 0
    fi

    if command -v sha256sum &>/dev/null; then
        actual="$(sha256sum "$binary" | awk '{print $1}')"
    elif command -v shasum &>/dev/null; then
        actual="$(shasum -a 256 "$binary" | awk '{print $1}')"
    else
        warn "No sha256sum or shasum found, skipping verification."
        return 0
    fi

    if [ "$expected" != "$actual" ]; then
        error "Checksum mismatch!"
        error "  expected: $expected"
        error "  actual:   $actual"
        rm -f "$binary"
        return 1
    fi

    info "Checksum verified."
}

# --- Download binary ---
download_binary() {
    local platform="$1"
    local version="$2"
    local dest="$3"
    local asset_name="adbclaw-${platform}"
    local download_url="https://github.com/${REPO}/releases/download/${version}/${asset_name}"

    info "Downloading adbclaw ${version} for ${platform}..."

    if command -v curl &>/dev/null; then
        curl -fSL --progress-bar -o "$dest" "$download_url"
    elif command -v wget &>/dev/null; then
        wget -q --show-progress -O "$dest" "$download_url"
    fi

    chmod +x "$dest"
}

# --- Check ADB ---
check_adb() {
    if command -v adb &>/dev/null; then
        info "adb found: $(command -v adb)"
    else
        echo ""
        warn "adb not found in PATH. adbclaw requires adb to work."
        warn "Install Android SDK Platform-Tools:"
        warn "  macOS:   brew install android-platform-tools"
        warn "  Linux:   sudo apt install android-tools-adb"
        warn "  Manual:  https://developer.android.com/tools/releases/platform-tools"
    fi
}

# --- Main ---
main() {
    local platform version dest

    echo -e "${BOLD}adbclaw installer${NC}"
    echo ""

    platform="$(detect_platform)"

    if [ -n "${VERSION:-}" ]; then
        version="$VERSION"
        info "Installing specified version: ${version}"
    else
        version="$(get_latest_version)"
        info "Latest version: ${version}"
    fi

    mkdir -p "$INSTALL_DIR"
    dest="${INSTALL_DIR}/adbclaw"

    # Download
    download_binary "$platform" "$version" "$dest"

    # Verify checksum
    verify_checksum "$dest" "$version" "$platform"

    # Verify it works
    if "$dest" --version &>/dev/null; then
        info "Installed successfully: $dest"
        info "Version: $("$dest" --version 2>&1 || true)"
    else
        info "Installed to: $dest"
    fi

    # Check PATH
    if ! echo "$PATH" | tr ':' '\n' | grep -q "^${INSTALL_DIR}$"; then
        echo ""
        warn "${INSTALL_DIR} is not in your PATH."
        warn "Add it to your shell profile:"
        warn "  echo 'export PATH=\"${INSTALL_DIR}:\$PATH\"' >> ~/.zshrc"
    fi

    # Check ADB
    check_adb

    echo ""
    info "Done! Run 'adbclaw doctor' to verify your setup."
}

main
