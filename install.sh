#!/bin/bash

# KillPort Auto-Installer for macOS/Linux
# This script downloads and installs killport automatically

set -e

echo "🔫 KillPort Auto-Installer"
echo "=========================="

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "❌ Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

case $OS in
    darwin)
        BINARY_NAME="killport-darwin-$ARCH"
        ;;
    linux)
        BINARY_NAME="killport-linux-$ARCH"
        ;;
    *)
        echo "❌ Unsupported OS: $OS"
        exit 1
        ;;
esac

DOWNLOAD_URL="https://github.com/tarantino19/killport/releases/latest/download/$BINARY_NAME"
INSTALL_DIR="/usr/local/bin"
INSTALL_PATH="$INSTALL_DIR/killport"

echo "📋 Detected: $OS $ARCH"
echo "📥 Downloading: $BINARY_NAME"

# Download the binary
if command -v curl >/dev/null 2>&1; then
    curl -L -o /tmp/killport "$DOWNLOAD_URL"
elif command -v wget >/dev/null 2>&1; then
    wget -O /tmp/killport "$DOWNLOAD_URL"
else
    echo "❌ Neither curl nor wget found. Please install one of them."
    exit 1
fi

echo "✅ Downloaded successfully"

# Make it executable
chmod +x /tmp/killport

# Check if we need sudo
if [ -w "$INSTALL_DIR" ]; then
    echo "📁 Installing to $INSTALL_PATH"
    mv /tmp/killport "$INSTALL_PATH"
else
    echo "📁 Installing to $INSTALL_PATH (requires sudo)"
    sudo mv /tmp/killport "$INSTALL_PATH"
fi

echo "✅ killport installed successfully!"
echo ""
echo "🎯 Try it out:"
echo "   killport list"
echo "   killport 3000"
echo ""
echo "📖 Need help? Visit: https://github.com/tarantino19/killport"