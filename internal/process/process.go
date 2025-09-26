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

func (m *Manager) KillProcess(pid string) error {
	if m.platform.IsWindows() {
		return m.killProcessWindows(pid)
	}
	return m.killProcessUnix(pid)
}

func (m *Manager) getProcessByPortUnix(port string) (*ProcessInfo, error) {
	cmd := exec.Command("lsof", "-ti:"+port)
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

func (m *Manager) getAllProcessesUnix() ([]*ProcessInfo, error) {
	cmd := exec.Command("lsof", "-i", "-P", "-n")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %v", err)
	}

	var processes []*ProcessInfo
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if strings.Contains(line, "LISTEN") {
			fields := strings.Fields(line)
			if len(fields) >= 9 {
				name := fields[0]
				pid := fields[1]
				address := fields[8]

				re := regexp.MustCompile(`:(\d+)$`)
				matches := re.FindStringSubmatch(address)
				if len(matches) > 1 {
					port := matches[1]
					processes = append(processes, &ProcessInfo{
						PID:    pid,
						Port:   port,
						Name:   name,
						Status: "LISTEN",
					})
				}
			}
		}
	}

	return processes, nil
}

func (m *Manager) getAllProcessesWindows() ([]*ProcessInfo, error) {
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %v", err)
	}

	var processes []*ProcessInfo
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
	cmd := exec.Command("lsof", "-ti:"+port)
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
