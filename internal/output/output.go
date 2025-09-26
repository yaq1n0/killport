package output

import (
	"fmt"

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
		name := proc.Name
		if len(name) > 20 {
			name = name[:17] + "..."
		}
		fmt.Printf("%-8s %-8s %-20s %s\n", proc.PID, proc.Port, name, proc.Status)
	}
}
