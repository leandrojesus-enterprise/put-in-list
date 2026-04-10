// Package terminal provides a cross-platform helper for clearing the terminal screen.
package terminal

import (
	"os"
	"os/exec"
	"runtime"
)

// Terminal holds OS-specific functions for terminal operations.
type Terminal struct {
	// clearFuncs maps an OS name to its screen-clearing command.
	clearFuncs map[string]func()
}

// New creates a Terminal and registers clear commands for each supported OS.
func New() *Terminal {
	var t *Terminal = &Terminal{}
	t.clearFuncs = map[string]func(){
		"windows": func() {
			// Use cmd.exe's cls command to clear the Windows console.
			var cmd *exec.Cmd = exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()
		},
		"linux": func() {
			// Use the standard clear utility on Linux.
			var cmd *exec.Cmd = exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
		},
	}
	return t
}

// Clear erases the terminal screen using the command appropriate for the current OS.
// It is a no-op on unsupported platforms.
func (t *Terminal) Clear() {
	if fn, ok := t.clearFuncs[runtime.GOOS]; ok {
		fn()
	}
}
