#!/bin/sh
# VibeDiff Installation Script
# Usage: curl -sSL https://raw.githubusercontent.com/vibediff/vibediff/main/install.sh | sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() { printf "${GREEN}%s${NC}\n" "$1"; }
error() { printf "${RED}%s${NC}\n" "$1"; }
warn() { printf "${YELLOW}%s${NC}\n" "$1"; }

# Detect platform
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux*)     OS=linux ;;
    Darwin*)    OS=darwin ;;
    MINGW*|MSYS*|CYGWIN*) OS=windows ;;
    *)          error "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
    x86_64|amd64) ARCH=amd64 ;;
    aarch64|arm64) ARCH=arm64 ;;
    *) error "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest version from GitHub API
info "Detecting latest version..."
LATEST_URL="https://api.github.com/repos/vibediff/vibediff/releases/latest"
VERSION=$(curl -s "$LATEST_URL" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    warn "Could not detect latest version, using 'latest'"
    VERSION="latest"
else
    info "Latest version: $VERSION"
fi

# Determine install directory
BINDIR="${BINDIR:-$HOME/.local/bin}"
if [ "$OS" = "windows" ]; then
    BINDIR="$USERPROFILE/bin"
fi

# Create bin directory if it doesn't exist
mkdir -p "$BINDIR"

# Download binary
BINARY_NAME="vibediff"
if [ "$OS" = "windows" ]; then
    BINARY_NAME="vibediff.exe"
fi

DOWNLOAD_URL="https://github.com/shubhamgurunglama007-oss/vibediff/releases/download/$VERSION/vibediff-${OS}-${ARCH}${EXT}"

info "Downloading from $DOWNLOAD_URL"
if command -v curl >/dev/null 2>&1; then
    curl -sSL -o "$BINDIR/$BINARY_NAME" "$DOWNLOAD_URL"
elif command -v wget >/dev/null 2>&1; then
    wget -q -O "$BINDIR/$BINARY_NAME" "$DOWNLOAD_URL"
else
    error "Neither curl nor wget found. Please install one and try again."
    exit 1
fi

# Make executable
chmod +x "$BINDIR/$BINARY_NAME"

info ""
info "Successfully installed VibeDiff!"
info ""
info "Add to PATH (if not already):"
warn "  export PATH=\"\$PATH:$BINDIR\""
info ""
info "Run:"
warn "  vibediff version"
info ""
info "Uninstall with:"
warn "  rm $BINDIR/$BINARY_NAME"
info ""
