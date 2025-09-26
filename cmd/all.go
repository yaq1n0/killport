package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/tarantino19/killport/internal/output"
	"github.com/tarantino19/killport/internal/process"

	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Kill all processes listening on ports",
	Long:  "Kill all processes currently listening on any port (use with caution)",
	Run:   runAll,
}

func runAll(cmd *cobra.Command, args []string) {
	manager := process.NewManager()
	formatter := output.NewFormatter()

	processes, err := manager.GetAllProcesses()
	if err != nil {
		formatter.Error("Failed to get process list: " + err.Error())
		return
	}

	if len(processes) == 0 {
		formatter.Warning("No processes found listening on ports")
		return
	}

	formatter.Warning(fmt.Sprintf("This will kill %d processes listening on ports:", len(processes)))
	formatter.PrintProcessTable(processes)

	fmt.Print("Are you sure you want to kill all these processes? (y/N): ")
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(strings.TrimSpace(response)) != "y" {
		formatter.Info("Operation cancelled")
		return
	}

	successCount := 0
	failureCount := 0

	for _, processInfo := range processes {
		err := manager.KillProcess(processInfo.PID)
		if err != nil {
			formatter.Error(fmt.Sprintf("Failed to kill %s (PID: %s, Port: %s) - %v",
				processInfo.Name, processInfo.PID, processInfo.Port, err))
			failureCount++
		} else {
			formatter.Success(fmt.Sprintf("Killed %s (PID: %s, Port: %s)",
				processInfo.Name, processInfo.PID, processInfo.Port))
			successCount++
		}
	}

	formatter.Info(fmt.Sprintf("Summary: %d killed, %d failed", successCount, failureCount))

	if failureCount > 0 {
		os.Exit(1)
	}
}
