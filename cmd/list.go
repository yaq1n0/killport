package cmd

import (
	"github.com/tarantino19/killport/internal/output"
	"github.com/tarantino19/killport/internal/process"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all active ports",
	Long:  "Display all processes currently using ports, grouped by port",
	Run:   runList,
}

func runList(cmd *cobra.Command, args []string) {
	manager := process.NewManager()
	formatter := output.NewFormatter()

	processes, err := manager.GetAllProcessesWithConnections()
	if err != nil {
		formatter.Error("Failed to get process list: " + err.Error())
		return
	}

	formatter.PrintGroupedProcessTable(processes)
}
