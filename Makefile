.PHONY: build install clean cross-compile test

# Build for current platform
build:
	@echo "Building killport for current platform..."
	@mkdir -p bin
	@go build -o bin/killport main.go
	@echo "Build complete: bin/killport"

# Install to system
install: build
	@echo "Installing killport to /usr/local/bin..."
	@sudo cp bin/killport /usr/local/bin/
	@sudo chmod +x /usr/local/bin/killport
	@echo "killport installed successfully! You can now use 'killport' from anywhere."

# Uninstall from system
uninstall:
	@echo "Removing killport from /usr/local/bin..."
	@sudo rm -f /usr/local/bin/killport
	@echo "killport uninstalled successfully."

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "Clean complete."

# Cross-compile for all platforms
cross-compile: clean
	@echo "Cross-compiling for all platforms..."
	@mkdir -p bin
	@echo "Building for Linux AMD64..."
	@GOOS=linux GOARCH=amd64 go build -o bin/killport-linux-amd64 main.go
	@echo "Building for Linux ARM64..."
	@GOOS=linux GOARCH=arm64 go build -o bin/killport-linux-arm64 main.go
	@echo "Building for Windows AMD64..."
	@GOOS=windows GOARCH=amd64 go build -o bin/killport-windows-amd64.exe main.go
	@echo "Building for macOS AMD64..."
	@GOOS=darwin GOARCH=amd64 go build -o bin/killport-darwin-amd64 main.go
	@echo "Building for macOS ARM64 (Apple Silicon)..."
	@GOOS=darwin GOARCH=arm64 go build -o bin/killport-darwin-arm64 main.go
	@echo "Cross-compilation complete. Binaries available in bin/"

# Test the application
test:
	@echo "Running tests..."
	@go test ./...
	@echo "Tests complete."

# Run the application with list command for testing
demo: build
	@echo "Running killport list to demonstrate functionality:"
	@./bin/killport list

help:
	@echo "Available commands:"
	@echo "  build         - Build killport for current platform"
	@echo "  install       - Install killport to /usr/local/bin (requires sudo)"
	@echo "  uninstall     - Remove killport from /usr/local/bin (requires sudo)"
	@echo "  clean         - Remove build artifacts"
	@echo "  cross-compile - Build for all supported platforms"
	@echo "  test          - Run tests"
	@echo "  demo          - Build and run 'killport list' for testing"
	@echo "  help          - Show this help message"
