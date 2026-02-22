package process

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/tarantino19/killport/internal/platform"
)

type ProcessInfo struct {
	PID    string
	Port   string
	Name   string
	Status string
}

type Manager struct {
	platform *platform.Platform
}

func NewManager() *Manager {
	return &Manager{
		platform: platform.New(),
	}
}

func (m *Manager) GetProcessByPort(port string) (*ProcessInfo, error) {
	if m.platform.IsWindows() {
		return m.getProcessByPortWindows(port)
	}
	return m.getProcessByPortUnix(port)
}

func (m *Manager) GetAllProcesses() ([]*ProcessInfo, error) {
	if m.platform.IsWindows() {
		return m.getAllProcessesWindows()
	}
	return m.getAllProcessesUnix()
}

// GetAllProcessesWithConnections returns all processes on all ports including
// non-LISTEN states (ESTABLISHED, CLOSE_WAIT, etc.) for the grouped list view.
func (m *Manager) GetAllProcessesWithConnections() ([]*ProcessInfo, error) {
	if m.platform.IsWindows() {
		return m.getAllProcessesWindows()
	}
	return m.getAllProcessesWithConnectionsUnix()
}

// GetAllProcessesByPort returns all processes associated with a specific port,
// including connected clients (ESTABLISHED, etc.), not just listeners.
func (m *Manager) GetAllProcessesByPort(port string) ([]*ProcessInfo, error) {
	if m.platform.IsWindows() {
		return m.getAllProcessesByPortWindows(port)
	}
	return m.getAllProcessesByPortUnix(port)
}

func (m *Manager) KillProcess(pid string) error {
	if m.platform.IsWindows() {
		return m.killProcessWindows(pid)
	}
	return m.killProcessUnix(pid)
}

func (m *Manager) getProcessByPortUnix(port string) (*ProcessInfo, error) {
	cmd := exec.Command("lsof", "-ti:"+port, "-sTCP:LISTEN")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("no process found on port %s", port)
	}

	pidOutput := strings.TrimSpace(string(output))
	if pidOutput == "" {
		return nil, fmt.Errorf("no process found on port %s", port)
	}

	// Handle multiple PIDs - take the first one
	pids := strings.Split(pidOutput, "\n")
	pid := strings.TrimSpace(pids[0])
	if pid == "" {
		return nil, fmt.Errorf("no process found on port %s", port)
	}

	nameCmd := exec.Command("ps", "-p", pid, "-o", "comm=")
	nameOutput, err := nameCmd.Output()
	if err != nil {
		return &ProcessInfo{PID: pid, Port: port, Name: "unknown", Status: "LISTEN"}, nil
	}

	name := strings.TrimSpace(string(nameOutput))
	return &ProcessInfo{
		PID:    pid,
		Port:   port,
		Name:   name,
		Status: "LISTEN",
	}, nil
}

func (m *Manager) getProcessByPortWindows(port string) (*ProcessInfo, error) {
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get process info: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, ":"+port) && strings.Contains(line, "LISTENING") {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				pid := fields[len(fields)-1]

				nameCmd := exec.Command("tasklist", "/fi", "PID eq "+pid, "/fo", "csv", "/nh")
				nameOutput, err := nameCmd.Output()
				name := "unknown"
				if err == nil {
					csvFields := strings.Split(strings.TrimSpace(string(nameOutput)), ",")
					if len(csvFields) > 0 {
						name = strings.Trim(csvFields[0], "\"")
					}
				}

				return &ProcessInfo{
					PID:    pid,
					Port:   port,
					Name:   name,
					Status: "LISTENING",
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("no process found on port %s", port)
}

// parseLsofLine parses a single lsof output line into a ProcessInfo.
// Returns nil if the line cannot be parsed.
func parseLsofLine(line string) *ProcessInfo {
	fields := strings.Fields(line)
	if len(fields) < 10 {
		return nil
	}

	name := strings.ReplaceAll(fields[0], "\\x20", " ")
	pid := fields[1]
	// The address is the second-to-last field, status is the last
	address := fields[len(fields)-2]
	status := strings.Trim(fields[len(fields)-1], "()")

	re := regexp.MustCompile(`:(\d+)$`)
	matches := re.FindStringSubmatch(address)
	if len(matches) < 2 {
		return nil
	}

	return &ProcessInfo{
		PID:    pid,
		Port:   matches[1],
		Name:   name,
		Status: status,
	}
}

func (m *Manager) getAllProcessesUnix() ([]*ProcessInfo, error) {
	cmd := exec.Command("lsof", "-i", "-P", "-n", "+c", "0")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %v", err)
	}

	var processes []*ProcessInfo
	seen := make(map[string]bool)
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if !strings.Contains(line, "LISTEN") {
			continue
		}
		proc := parseLsofLine(line)
		if proc == nil {
			continue
		}
		key := proc.PID + ":" + proc.Port
		if seen[key] {
			continue
		}
		seen[key] = true
		processes = append(processes, proc)
	}

	return processes, nil
}

// getAllProcessesWithConnectionsUnix returns all processes across all ports,
// including ESTABLISHED, CLOSE_WAIT, etc. — not just LISTEN.
func (m *Manager) getAllProcessesWithConnectionsUnix() ([]*ProcessInfo, error) {
	cmd := exec.Command("lsof", "-i", "-P", "-n", "+c", "0")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %v", err)
	}

	var processes []*ProcessInfo
	seen := make(map[string]bool)
	lines := strings.Split(string(output), "\n")

	// Skip header line
	for _, line := range lines {
		if strings.HasPrefix(line, "COMMAND") || line == "" {
			continue
		}
		proc := parseLsofLine(line)
		if proc == nil {
			continue
		}
		key := proc.PID + ":" + proc.Port + ":" + proc.Status
		if seen[key] {
			continue
		}
		seen[key] = true
		processes = append(processes, proc)
	}

	return processes, nil
}

// getAllProcessesByPortUnix returns all processes on a specific port (all states).
func (m *Manager) getAllProcessesByPortUnix(port string) ([]*ProcessInfo, error) {
	cmd := exec.Command("lsof", "-i:"+port, "-P", "-n", "+c", "0")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("no processes found on port %s", port)
	}

	var processes []*ProcessInfo
	seen := make(map[string]bool)
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "COMMAND") || line == "" {
			continue
		}
		proc := parseLsofLine(line)
		if proc == nil {
			continue
		}
		key := proc.PID + ":" + proc.Port + ":" + proc.Status
		if seen[key] {
			continue
		}
		seen[key] = true
		processes = append(processes, proc)
	}

	return processes, nil
}

// getAllProcessesByPortWindows returns all processes on a specific port (Windows).
// Currently only returns LISTENING processes due to netstat parsing limitations.
func (m *Manager) getAllProcessesByPortWindows(port string) ([]*ProcessInfo, error) {
	proc, err := m.getProcessByPortWindows(port)
	if err != nil {
		return nil, err
	}
	return []*ProcessInfo{proc}, nil
}

func (m *Manager) getAllProcessesWindows() ([]*ProcessInfo, error) {
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %v", err)
	}

	var processes []*ProcessInfo
	seen := make(map[string]bool)
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if strings.Contains(line, "LISTENING") {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				address := fields[1]
				pid := fields[len(fields)-1]

				parts := strings.Split(address, ":")
				if len(parts) >= 2 {
					port := parts[len(parts)-1]

					key := pid + ":" + port
					if seen[key] {
						continue
					}
					seen[key] = true

					nameCmd := exec.Command("tasklist", "/fi", "PID eq "+pid, "/fo", "csv", "/nh")
					nameOutput, _ := nameCmd.Output()
					name := "unknown"
					if len(nameOutput) > 0 {
						csvFields := strings.Split(strings.TrimSpace(string(nameOutput)), ",")
						if len(csvFields) > 0 {
							name = strings.Trim(csvFields[0], "\"")
						}
					}

					processes = append(processes, &ProcessInfo{
						PID:    pid,
						Port:   port,
						Name:   name,
						Status: "LISTENING",
					})
				}
			}
		}
	}

	return processes, nil
}

func (m *Manager) killProcessUnix(pid string) error {
	cmd := exec.Command("kill", "-9", pid)
	return cmd.Run()
}

func (m *Manager) killProcessWindows(pid string) error {
	cmd := exec.Command("taskkill", "/F", "/PID", pid)
	return cmd.Run()
}

func (m *Manager) KillProcessByPort(port string) error {
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		return fmt.Errorf("invalid port number: %s", port)
	}

	return m.KillAllProcessesByPort(port)
}

func (m *Manager) KillAllProcessesByPort(port string) error {
	if m.platform.IsWindows() {
		return m.killAllProcessesByPortWindows(port)
	}
	return m.killAllProcessesByPortUnix(port)
}

func (m *Manager) killAllProcessesByPortUnix(port string) error {
	cmd := exec.Command("lsof", "-ti:"+port, "-sTCP:LISTEN")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("no process found on port %s", port)
	}

	pidOutput := strings.TrimSpace(string(output))
	if pidOutput == "" {
		return fmt.Errorf("no process found on port %s", port)
	}

	pids := strings.Split(pidOutput, "\n")
	var errors []string

	for _, pid := range pids {
		pid = strings.TrimSpace(pid)
		if pid != "" {
			err := m.killProcessUnix(pid)
			if err != nil {
				errors = append(errors, fmt.Sprintf("PID %s: %v", pid, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to kill some processes: %s", strings.Join(errors, "; "))
	}

	return nil
}

func (m *Manager) killAllProcessesByPortWindows(port string) error {
	// For Windows, we'll use the existing single process approach for now
	process, err := m.GetProcessByPort(port)
	if err != nil {
		return err
	}
	return m.KillProcess(process.PID)
}
