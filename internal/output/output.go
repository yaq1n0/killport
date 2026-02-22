package output

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/tarantino19/killport/internal/process"

	"github.com/fatih/color"
)

type Formatter struct {
	red    *color.Color
	green  *color.Color
	yellow *color.Color
	blue   *color.Color
}

func NewFormatter() *Formatter {
	return &Formatter{
		red:    color.New(color.FgRed),
		green:  color.New(color.FgGreen),
		yellow: color.New(color.FgYellow),
		blue:   color.New(color.FgBlue),
	}
}

func (f *Formatter) Success(message string) {
	f.green.Println("✓ " + message)
}

func (f *Formatter) Error(message string) {
	f.red.Println("✗ " + message)
}

func (f *Formatter) Warning(message string) {
	f.yellow.Println("⚠ " + message)
}

func (f *Formatter) Info(message string) {
	f.blue.Println("ℹ " + message)
}

func (f *Formatter) PrintProcessTable(processes []*process.ProcessInfo) {
	if len(processes) == 0 {
		f.Warning("No processes found")
		return
	}

	f.Info("Active ports:")

	// Simple table formatting without external library
	fmt.Printf("%-8s %-8s %-20s %s\n", "PID", "Port", "Process Name", "Status")
	fmt.Printf("%-8s %-8s %-20s %s\n", "---", "----", "------------", "------")

	for _, proc := range processes {
		name := truncateName(proc.Name, 20)
		fmt.Printf("%-8s %-8s %-20s %s\n", proc.PID, proc.Port, name, proc.Status)
	}
}

// PrintGroupedProcessTable displays processes grouped by port, sorted numerically.
// Within each port group, LISTEN processes are shown first, then others.
func (f *Formatter) PrintGroupedProcessTable(processes []*process.ProcessInfo) {
	if len(processes) == 0 {
		f.Warning("No processes found")
		return
	}

	// Group by port
	groups := make(map[string][]*process.ProcessInfo)
	for _, proc := range processes {
		groups[proc.Port] = append(groups[proc.Port], proc)
	}

	// Sort ports numerically
	ports := make([]string, 0, len(groups))
	for port := range groups {
		ports = append(ports, port)
	}
	sort.Slice(ports, func(i, j int) bool {
		pi, _ := strconv.Atoi(ports[i])
		pj, _ := strconv.Atoi(ports[j])
		return pi < pj
	})

	f.Info("Active ports:")
	fmt.Println()

	for _, port := range ports {
		procs := groups[port]

		// Sort: LISTEN first, then by PID
		sort.Slice(procs, func(i, j int) bool {
			if isListen(procs[i].Status) != isListen(procs[j].Status) {
				return isListen(procs[i].Status)
			}
			return procs[i].PID < procs[j].PID
		})

		fmt.Printf("Port %s:\n", port)
		fmt.Printf("  %-8s %-24s %s\n", "PID", "Process Name", "Status")

		for _, proc := range procs {
			name := truncateName(proc.Name, 24)
			fmt.Printf("  %-8s %-24s %s\n", proc.PID, name, proc.Status)
		}
		fmt.Println()
	}
}

func truncateName(name string, maxLen int) string {
	if len(name) > maxLen {
		return name[:maxLen-3] + "..."
	}
	return name
}

func isListen(status string) bool {
	return status == "LISTEN" || status == "LISTENING"
}
