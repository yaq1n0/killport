# KillPort 🔫

**Kill any process running on any port with one simple command.**

Works on macOS, Windows, and Linux. No setup required.

## What You Can Do (Copy & Paste These Commands)

```bash
# See what's running on your ports
killport list

# Kill whatever's on port 3000
killport 3000

# Kill multiple ports at once
killport 3000 8080 5000

# Kill everything (careful!)
killport all
```

**That's it.** Simple as that.

---

## How to Install (Pick One Method)

** Easiest Way (If you have Go installed):**
```bash
go install github.com/tarantino19/killport@latest
```
Then just use `killport` anywhere.

**🥈 Automated Install (Recommended):**

**macOS/Linux (one command):**
```bash
curl -sSL https://raw.githubusercontent.com/tarantino19/killport/main/install.sh | bash
```

**Windows (run as Administrator):**
```cmd
curl -o install.bat https://raw.githubusercontent.com/tarantino19/killport/main/install.bat && install.bat
```

**Manual Download:**

Download the pre-built binary for your system:

**macOS:**
- Intel Macs: [Download killport-darwin-amd64](https://github.com/tarantino19/killport/releases/latest/download/killport-darwin-amd64)
- Apple Silicon (M1/M2): [Download killport-darwin-arm64](https://github.com/tarantino19/killport/releases/latest/download/killport-darwin-arm64)

**Linux:**
- 64-bit: [Download killport-linux-amd64](https://github.com/tarantino19/killport/releases/latest/download/killport-linux-amd64)
- ARM64: [Download killport-linux-arm64](https://github.com/tarantino19/killport/releases/latest/download/killport-linux-arm64)

**Windows:**
- 64-bit: [Download killport-windows-amd64.exe](https://github.com/tarantino19/killport/releases/latest/download/killport-windows-amd64.exe)

After manual download:
1. Rename the file to `killport` (or `killport.exe` on Windows)
2. Make it executable: `chmod +x killport` (macOS/Linux)
3. Move it to a directory in your PATH (e.g., `/usr/local/bin/`)

**🔧 Build from Source:**
```bash
git clone https://github.com/tarantino19/killport.git
cd killport
make build
sudo cp bin/killport /usr/local/bin/
```

---

## Examples (What You'll See)

```bash
$ killport list
ℹ Active ports:
PID      Port     Process Name         Status
---      ----     ------------         ------
1234     3000     node                 LISTEN
5678     8080     java                 LISTEN
9012     5432     postgres             LISTEN
```

**When you run `killport 3000`:**
```bash
$ killport 3000
Attempting to kill process on port 3000...
✓ Port 3000: Killed node (PID: 1234)
```

**When you run `killport 3000 8080 5432`:**
```bash
$ killport 3000 8080 5432
Attempting to kill process on port 3000...
✓ Port 3000: Killed node (PID: 1234)
Attempting to kill process on port 8080...
✓ Port 8080: Killed java (PID: 5678)
Attempting to kill process on port 5432...
✓ Port 5432: Killed postgres (PID: 9012)
ℹ Processed 3 ports
```

**When you run `killport all` (BE CAREFUL!):**
```bash
$ killport all
⚠ This will kill 3 processes listening on ports:
PID      Port     Process Name         Status
---      ----     ------------         ------
1234     3000     node                 LISTEN
5678     8080     java                 LISTEN
9012     5432     postgres             LISTEN

Are you sure you want to kill all these processes? (y/N): y
✓ Killed node (PID: 1234, Port: 3000)
✓ Killed java (PID: 5678, Port: 8080)
✓ Killed postgres (PID: 9012, Port: 5432)
ℹ Summary: 3 killed, 0 failed
```

---

## Important Notes

- **`killport all` is dangerous** - it kills ALL processes using ports
- **Save your work first** before killing processes
- **You can't undo** killing a process
- **Some system processes will restart** automatically
- **Development servers will stop** and you'll lose unsaved changes

---

## Problems?

[Open an issue here](https://github.com/tarantino19/killport/issues) and I'll help you out.

---

## For Developers

**Build from source**

```bash
git clone https://github.com/tarantino19/killport.git
cd killport
make build
```

**Available commands:**
- `make build` - Build it
- `make install` - Install it
- `make cross-compile` - Build for all platforms
- `make clean` - Clean up

**Built with:** Go + [Cobra CLI](https://github.com/spf13/cobra) + [Color](https://github.com/fatih/color)

**Works on:** macOS, Linux, Windows (AMD64, ARM64)
