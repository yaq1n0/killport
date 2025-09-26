package cmd

import (
	"fmt"

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

		fmt.Printf("Attempting to kill process on port %s...\n", port)

		processInfo, err := manager.GetProcessByPort(port)
		if err != nil {
			formatter.Error(fmt.Sprintf("Port %s: %v", port, err))
			continue
		}

		err = manager.KillProcess(processInfo.PID)
		if err != nil {
			formatter.Error(fmt.Sprintf("Port %s: Failed to kill process (PID: %s) - %v", port, processInfo.PID, err))
			continue
		}

		formatter.Success(fmt.Sprintf("Port %s: Killed %s (PID: %s)", port, processInfo.Name, processInfo.PID))
	}

	if len(args) > 1 {
		formatter.Info(fmt.Sprintf("Processed %d ports", len(args)))
	}
}
