#!/bin/bash

set -e

echo "🔨 Building killport..."
go build -o bin/killport main.go

echo "📦 Installing killport to /usr/local/bin..."
sudo cp bin/killport /usr/local/bin/
sudo chmod +x /usr/local/bin/killport

echo "✅ killport installed successfully!"
echo "ℹ️  You can now use 'killport' from anywhere in your terminal."
echo ""
echo "Usage examples:"
echo "  killport list          - List all active ports"
echo "  killport 3000          - Kill process on port 3000"
echo "  killport 3000 4000     - Kill processes on multiple ports"
echo "  killport all           - Kill all port processes (with confirmation)"
echo ""
echo "To uninstall, run: sudo rm /usr/local/bin/killport"
