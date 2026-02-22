//go:build !windows

package process

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

// startListener spawns a nc -l process on the given port.
// Returns the exec.Cmd so we can check if it's still running.
// Registers a cleanup to kill the process if the test fails.
func startListener(t *testing.T, port int) *exec.Cmd {
	t.Helper()
	cmd := exec.Command("nc", "-l", strconv.Itoa(port))
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start listener on port %d: %v", port, err)
	}
	// Give nc a moment to bind
	time.Sleep(200 * time.Millisecond)
	t.Cleanup(func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
			cmd.Wait()
		}
	})
	return cmd
}

// isProcessAlive checks if a process with the given PID is still running.
func isProcessAlive(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// On Unix, FindProcess always succeeds. Signal 0 checks if process exists.
	err = p.Signal(os.Signal(nil))
	// If err is nil, process is alive. If "process already finished", it's dead.
	return err == nil
}

func TestGetProcessByPort(t *testing.T) {
	port := 19501
	listener := startListener(t, port)

	manager := NewManager()
	proc, err := manager.GetProcessByPort(strconv.Itoa(port))
	if err != nil {
		t.Fatalf("GetProcessByPort(%d) error: %v", port, err)
	}

	if proc.Port != strconv.Itoa(port) {
		t.Errorf("Port = %q, want %q", proc.Port, strconv.Itoa(port))
	}
	if proc.PID != strconv.Itoa(listener.Process.Pid) {
		t.Errorf("PID = %q, want %q", proc.PID, strconv.Itoa(listener.Process.Pid))
	}
	if proc.Status != "LISTEN" {
		t.Errorf("Status = %q, want LISTEN", proc.Status)
	}
}

func TestGetProcessByPort_NoProcess(t *testing.T) {
	manager := NewManager()
	_, err := manager.GetProcessByPort("19599")
	if err == nil {
		t.Error("GetProcessByPort on unused port should return error")
	}
	if !strings.Contains(err.Error(), "no process found") {
		t.Errorf("error = %q, want it to contain 'no process found'", err.Error())
	}
}

func TestGetAllProcesses_ContainsListener(t *testing.T) {
	port := 19502
	listener := startListener(t, port)

	manager := NewManager()
	procs, err := manager.GetAllProcesses()
	if err != nil {
		t.Fatalf("GetAllProcesses() error: %v", err)
	}

	found := false
	for _, p := range procs {
		if p.Port == strconv.Itoa(port) && p.PID == strconv.Itoa(listener.Process.Pid) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("GetAllProcesses() did not include listener on port %d (PID %d)", port, listener.Process.Pid)
	}
}

func TestGetAllProcessesByPort(t *testing.T) {
	port := 19503
	listener := startListener(t, port)

	manager := NewManager()
	procs, err := manager.GetAllProcessesByPort(strconv.Itoa(port))
	if err != nil {
		t.Fatalf("GetAllProcessesByPort(%d) error: %v", port, err)
	}

	if len(procs) == 0 {
		t.Fatal("GetAllProcessesByPort returned no processes")
	}

	found := false
	for _, p := range procs {
		if p.PID == strconv.Itoa(listener.Process.Pid) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("GetAllProcessesByPort(%d) did not include our listener PID %d", port, listener.Process.Pid)
	}
}

func TestKillProcess(t *testing.T) {
	port := 19504
	listener := startListener(t, port)
	pid := listener.Process.Pid

	manager := NewManager()
	err := manager.KillProcess(strconv.Itoa(pid))
	if err != nil {
		t.Fatalf("KillProcess(%d) error: %v", pid, err)
	}

	// Wait for process to die
	time.Sleep(100 * time.Millisecond)

	if isProcessAlive(pid) {
		t.Errorf("process %d is still alive after KillProcess", pid)
	}
}

func TestKillProcessByPort(t *testing.T) {
	port := 19505
	listener := startListener(t, port)
	pid := listener.Process.Pid

	manager := NewManager()
	err := manager.KillProcessByPort(strconv.Itoa(port))
	if err != nil {
		t.Fatalf("KillProcessByPort(%d) error: %v", port, err)
	}

	time.Sleep(100 * time.Millisecond)

	if isProcessAlive(pid) {
		t.Errorf("process %d is still alive after KillProcessByPort", pid)
	}
}

func TestKillAllProcessesByPort(t *testing.T) {
	port := 19506
	listener := startListener(t, port)
	pid := listener.Process.Pid

	manager := NewManager()
	err := manager.KillAllProcessesByPort(strconv.Itoa(port))
	if err != nil {
		t.Fatalf("KillAllProcessesByPort(%d) error: %v", port, err)
	}

	time.Sleep(100 * time.Millisecond)

	if isProcessAlive(pid) {
		t.Errorf("process %d is still alive after KillAllProcessesByPort", pid)
	}
}

func TestKillProcessByPort_NoProcess(t *testing.T) {
	manager := NewManager()
	err := manager.KillProcessByPort("19598")
	if err == nil {
		t.Error("KillProcessByPort on unused port should return error")
	}
}

func TestKillProcessByPort_InvalidPort(t *testing.T) {
	manager := NewManager()
	err := manager.KillProcessByPort("99999")
	if err == nil {
		t.Error("KillProcessByPort with invalid port should return error")
	}
	if !strings.Contains(err.Error(), "invalid port") {
		t.Errorf("error = %q, want it to contain 'invalid port'", err.Error())
	}
}

func TestGetAllProcessesWithConnections(t *testing.T) {
	port := 19507
	listener := startListener(t, port)

	manager := NewManager()
	procs, err := manager.GetAllProcessesWithConnections()
	if err != nil {
		t.Fatalf("GetAllProcessesWithConnections() error: %v", err)
	}

	found := false
	for _, p := range procs {
		if p.Port == strconv.Itoa(port) && p.PID == strconv.Itoa(listener.Process.Pid) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("GetAllProcessesWithConnections() did not include listener on port %d", port)
	}
}

func TestMultiplePortListeners(t *testing.T) {
	port1 := 19508
	port2 := 19509
	l1 := startListener(t, port1)
	l2 := startListener(t, port2)

	manager := NewManager()
	procs, err := manager.GetAllProcesses()
	if err != nil {
		t.Fatalf("GetAllProcesses() error: %v", err)
	}

	found1, found2 := false, false
	for _, p := range procs {
		if p.PID == strconv.Itoa(l1.Process.Pid) {
			found1 = true
		}
		if p.PID == strconv.Itoa(l2.Process.Pid) {
			found2 = true
		}
	}

	if !found1 {
		t.Errorf("GetAllProcesses() missing listener on port %d", port1)
	}
	if !found2 {
		t.Errorf("GetAllProcesses() missing listener on port %d", port2)
	}

	// Kill one, verify the other is still found
	err = manager.KillProcessByPort(strconv.Itoa(port1))
	if err != nil {
		t.Fatalf("KillProcessByPort(%d) error: %v", port1, err)
	}
	time.Sleep(100 * time.Millisecond)

	proc2, err := manager.GetProcessByPort(strconv.Itoa(port2))
	if err != nil {
		t.Fatalf("GetProcessByPort(%d) after killing %d: %v", port2, port1, err)
	}
	if proc2.PID != fmt.Sprintf("%d", l2.Process.Pid) {
		t.Errorf("PID = %q, want %d", proc2.PID, l2.Process.Pid)
	}
}
