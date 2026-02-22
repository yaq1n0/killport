package cmd

import (
	"fmt"
	"strconv"

	"github.com/tarantino19/killport/internal/output"
	"github.com/tarantino19/killport/internal/process"

	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill processes on specified ports",
	Long:  "Kill processes running on one or more specified ports",
	Args:  cobra.MinimumNArgs(1),
	Run:   runKill,
}

func runKill(cmd *cobra.Command, args []string) {
	manager := process.NewManager()
	formatter := output.NewFormatter()

	for _, port := range args {
		if port == "list" || port == "all" {
			continue
		}

		portNum, err := strconv.Atoi(port)
		if err != nil || portNum < 1 || portNum > 65535 {
			formatter.Error(fmt.Sprintf("Invalid port: %s (must be 1-65535)", port))
			continue
		}

		fmt.Printf("Attempting to kill process on port %s...\n", port)

		if killAll {
			runKillAllOnPort(manager, formatter, port)
		} else {
			runKillListener(manager, formatter, port)
		}
	}

	if len(args) > 1 {
		formatter.Info(fmt.Sprintf("Processed %d ports", len(args)))
	}
}

func runKillListener(manager *process.Manager, formatter *output.Formatter, port string) {
	processInfo, err := manager.GetProcessByPort(port)
	if err != nil {
		formatter.Error(fmt.Sprintf("Port %s: %v", port, err))
		return
	}

	err = manager.KillProcess(processInfo.PID)
	if err != nil {
		formatter.Error(fmt.Sprintf("Port %s: Failed to kill process (PID: %s) - %v", port, processInfo.PID, err))
		return
	}

	formatter.Success(fmt.Sprintf("Port %s: Killed %s (PID: %s)", port, processInfo.Name, processInfo.PID))

	// Check for other connected processes
	allProcs, err := manager.GetAllProcessesByPort(port)
	if err != nil {
		return
	}

	// Count non-listener processes (excluding the one we just killed)
	var others int
	for _, p := range allProcs {
		if p.PID != processInfo.PID {
			others++
		}
	}
	if others > 0 {
		formatter.Info(fmt.Sprintf("%d other process(es) connected to port %s (use --all to include)", others, port))
	}
}

func runKillAllOnPort(manager *process.Manager, formatter *output.Formatter, port string) {
	allProcs, err := manager.GetAllProcessesByPort(port)
	if err != nil {
		formatter.Error(fmt.Sprintf("Port %s: %v", port, err))
		return
	}

	killed := make(map[string]bool)
	for _, proc := range allProcs {
		if killed[proc.PID] {
			continue
		}
		err := manager.KillProcess(proc.PID)
		if err != nil {
			formatter.Error(fmt.Sprintf("Port %s: Failed to kill %s (PID: %s) - %v", port, proc.Name, proc.PID, err))
		} else {
			formatter.Success(fmt.Sprintf("Port %s: Killed %s (PID: %s, %s)", port, proc.Name, proc.PID, proc.Status))
		}
		killed[proc.PID] = true
	}
}
