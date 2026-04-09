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
	t := &Terminal{}
	t.clearFuncs = map[string]func(){
		"windows": func() {
			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()
		},
		"linux": func() {
			cmd := exec.Command("clear")
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
