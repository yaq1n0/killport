package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "killport",
	Short: "Kill processes running on specified ports",
	Long: `A cross-platform CLI tool to kill processes running on specified ports.
Supports macOS, Windows, and Linux.

Examples:
  killport 3000          Kill process on port 3000
  killport 3000 4000     Kill processes on ports 3000 and 4000
  killport list          List all active ports
  killport all           Kill all port processes`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 && args[0] == "list" {
			listCmd.Run(cmd, []string{})
			return
		}
		if len(args) == 1 && args[0] == "all" {
			allCmd.Run(cmd, []string{})
			return
		}
		killCmd.Run(cmd, args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
