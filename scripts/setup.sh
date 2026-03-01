#!/usr/bin/env bash
#
# setup.sh — Download pre-built adbclaw binary from GitHub Releases.
# Called by the SessionStart hook to ensure the binary is available.
# Falls back to building from source if Go is installed.
#

set -euo pipefail

REPO="llm-net/adbclaw"

# Resolve plugin root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PLUGIN_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BIN_DIR="$PLUGIN_ROOT/bin"
BINARY="$BIN_DIR/adbclaw"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

info()  { echo -e "${GREEN}[adbclaw]${NC} $*" >&2; }
warn()  { echo -e "${YELLOW}[adbclaw]${NC} $*" >&2; }
error() { echo -e "${RED}[adbclaw]${NC} $*" >&2; }

# --- Detect OS and Arch ---
detect_platform() {
    local os arch
    os="$(uname -s | tr '[:upper:]' '[:lower:]')"
    arch="$(uname -m)"

    case "$os" in
        darwin) os="darwin" ;;
        linux)  os="linux" ;;
        *)      error "Unsupported OS: $os"; exit 1 ;;
    esac

    case "$arch" in
        x86_64|amd64)  arch="amd64" ;;
        aarch64|arm64) arch="arm64" ;;
        *)             error "Unsupported architecture: $arch"; exit 1 ;;
    esac

    echo "${os}-${arch}"
}

# --- Check ADB ---
check_adb() {
    if command -v adb &>/dev/null; then
        info "adb: $(adb version 2>/dev/null | head -1 || echo 'found')"
    else
        warn "adb not found in PATH. Install Android SDK Platform-Tools:"
        warn "  macOS:   brew install android-platform-tools"
        warn "  Linux:   sudo apt install android-tools-adb"
        warn "  Manual:  https://developer.android.com/tools/releases/platform-tools"
    fi
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
        error "Neither curl nor wget found. Cannot download binary."
        return 1
    fi

    if [ -z "$tag" ]; then
        error "Failed to get latest release version."
        return 1
    fi

    echo "$tag"
}

# --- Download binary ---
download_binary() {
    local platform="$1"
    local version="$2"
    local asset_name="adbclaw-${platform}"
    local download_url="https://github.com/${REPO}/releases/download/${version}/${asset_name}"

    info "Downloading adbclaw ${version} for ${platform}..."

    mkdir -p "$BIN_DIR"

    if command -v curl &>/dev/null; then
        curl -fSL --progress-bar -o "$BINARY" "$download_url"
    elif command -v wget &>/dev/null; then
        wget -q --show-progress -O "$BINARY" "$download_url"
    fi

    chmod +x "$BINARY"
    info "Downloaded to $BINARY"
}

# --- Build from source (fallback) ---
build_from_source() {
    local src_dir="$PLUGIN_ROOT/src"

    if ! command -v go &>/dev/null; then
        error "Go is not installed. Cannot build from source."
        error "Install Go 1.24+ from https://go.dev/dl/"
        return 1
    fi

    info "Building from source..."
    cd "$src_dir"

    # Build directly into plugin bin/ directory
    mkdir -p "$BIN_DIR"
    local version
    version="$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')"
    go mod tidy
    go build -ldflags "-X github.com/llm-net/adbclaw/cmd.Version=${version}" -o "$BINARY" .

    info "Built from source: $BINARY"
}

# --- Main ---
main() {
    # Already installed?
    if [ -f "$BINARY" ]; then
        info "adbclaw is ready: $BINARY"
        check_adb
        exit 0
    fi

    # Try downloading pre-built binary first
    local platform
    platform="$(detect_platform)"

    local version
    if version="$(get_latest_version)"; then
        if download_binary "$platform" "$version"; then
            info "adbclaw ${version} is ready."
            check_adb
            exit 0
        else
            warn "Download failed, trying to build from source..."
        fi
    else
        warn "Could not fetch latest release, trying to build from source..."
    fi

    # Fallback: build from source
    if build_from_source; then
        info "adbclaw is ready (built from source)."
        check_adb
        exit 0
    fi

    error "Failed to install adbclaw. Please check your network or install Go to build from source."
    exit 1
}

main
