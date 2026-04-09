package terminal

import (
	"os"
	"os/exec"
	"runtime"
)

type Terminal struct {
	clearFuncs map[string]func()
}

func New() *Terminal {
	var t *Terminal = &Terminal{}
	t.clearFuncs = map[string]func(){
		"windows": func() {
			var cmd *exec.Cmd = exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()
		},
		"linux": func() {
			var cmd *exec.Cmd = exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
		},
	}
	return t
}

func (t *Terminal) Clear() {
	if fn, ok := t.clearFuncs[runtime.GOOS]; ok {
		fn()
	}
}
