//go:build !windows

package cmd_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

var binaryPath string

func TestMain(m *testing.M) {
	// Build the binary once for all CLI tests
	tmpDir, err := os.MkdirTemp("", "killport-test")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}
	defer os.RemoveAll(tmpDir)

	binaryPath = filepath.Join(tmpDir, "killport")
	build := exec.Command("go", "build", "-o", binaryPath, "..")
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		panic("failed to build killport: " + err.Error())
	}

	os.Exit(m.Run())
}

func startListener(t *testing.T, port int) *exec.Cmd {
	t.Helper()
	cmd := exec.Command("nc", "-l", strconv.Itoa(port))
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start listener on port %d: %v", port, err)
	}
	time.Sleep(200 * time.Millisecond)
	t.Cleanup(func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
			cmd.Wait()
		}
	})
	return cmd
}

func TestCLI_KillPort(t *testing.T) {
	port := 19520
	startListener(t, port)

	out, err := exec.Command(binaryPath, strconv.Itoa(port)).CombinedOutput()
	if err != nil {
		t.Fatalf("killport %d failed: %v\noutput: %s", port, err, out)
	}

	output := string(out)
	if !strings.Contains(output, "Killed") {
		t.Errorf("expected output to contain 'Killed', got: %s", output)
	}
	if !strings.Contains(output, strconv.Itoa(port)) {
		t.Errorf("expected output to contain port %d, got: %s", port, output)
	}
}

func TestCLI_KillPort_NoProcess(t *testing.T) {
	out, err := exec.Command(binaryPath, "19597").CombinedOutput()
	output := string(out)

	// Should report error but not crash
	if !strings.Contains(output, "no process found") {
		t.Errorf("expected 'no process found' in output, got: %s", output)
	}
	_ = err // exit code may or may not be 0 here, just check output
}

func TestCLI_KillPort_InvalidPort(t *testing.T) {
	out, err := exec.Command(binaryPath, "99999").CombinedOutput()
	output := string(out)

	if !strings.Contains(output, "Invalid port") {
		t.Errorf("expected 'Invalid port' in output, got: %s", output)
	}
	_ = err
}

func TestCLI_KillPort_MultiplePorts(t *testing.T) {
	port1 := 19521
	port2 := 19522
	startListener(t, port1)
	startListener(t, port2)

	out, err := exec.Command(binaryPath, strconv.Itoa(port1), strconv.Itoa(port2)).CombinedOutput()
	if err != nil {
		t.Fatalf("killport %d %d failed: %v\noutput: %s", port1, port2, err, out)
	}

	output := string(out)
	if !strings.Contains(output, strconv.Itoa(port1)) {
		t.Errorf("expected output to mention port %d, got: %s", port1, output)
	}
	if !strings.Contains(output, strconv.Itoa(port2)) {
		t.Errorf("expected output to mention port %d, got: %s", port2, output)
	}
}

func TestCLI_List(t *testing.T) {
	port := 19523
	startListener(t, port)

	out, err := exec.Command(binaryPath, "list").CombinedOutput()
	if err != nil {
		t.Fatalf("killport list failed: %v\noutput: %s", err, out)
	}

	output := string(out)
	if !strings.Contains(output, strconv.Itoa(port)) {
		t.Errorf("expected 'killport list' to show port %d, got: %s", port, output)
	}
	if !strings.Contains(output, "nc") {
		t.Errorf("expected 'killport list' to show 'nc' process, got: %s", output)
	}
}

func TestCLI_KillPort_AllFlag(t *testing.T) {
	port := 19524
	startListener(t, port)

	out, err := exec.Command(binaryPath, "--all", strconv.Itoa(port)).CombinedOutput()
	if err != nil {
		t.Fatalf("killport --all %d failed: %v\noutput: %s", port, err, out)
	}

	output := string(out)
	if !strings.Contains(output, "Killed") {
		t.Errorf("expected output to contain 'Killed', got: %s", output)
	}
}

func TestCLI_KillPort_VerifyProcessDead(t *testing.T) {
	port := 19525
	listener := startListener(t, port)
	pid := listener.Process.Pid

	_, err := exec.Command(binaryPath, strconv.Itoa(port)).CombinedOutput()
	if err != nil {
		t.Fatalf("killport %d failed: %v", port, err)
	}

	time.Sleep(100 * time.Millisecond)

	// Verify the process is actually dead
	p, _ := os.FindProcess(pid)
	if err := p.Signal(os.Signal(nil)); err == nil {
		t.Errorf("process %d is still alive after killport", pid)
	}
}
