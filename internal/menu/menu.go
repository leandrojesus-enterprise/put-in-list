// Package menu handles rendering the interactive terminal UI and reading user input.
package menu

import (
	"fmt"

	"github.com/leandrojesus-enterprise/put-in-list/internal/i18n"
	"github.com/leandrojesus-enterprise/put-in-list/internal/terminal"
)

// ANSI escape codes for terminal color output.
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

// Menu wraps the terminal and translator to provide UI rendering and input helpers.
type Menu struct {
	term *terminal.Terminal
	tr   *i18n.Translator
}

// New creates a Menu with the given terminal and translator.
func New(term *terminal.Terminal, tr *i18n.Translator) *Menu {
	return &Menu{term: term, tr: tr}
}

// ShowMain renders the main menu, displaying the active list, current installer,
// app version, and a welcome message on the very first run.
func (m *Menu) ShowMain(activeList, currentInstaller, version string, firstRun bool) {
	if firstRun {
		// Display a one-time welcome banner and wait for the user to acknowledge it.
		fmt.Printf("           %v{%v %s %v}%v\n\n", Green, Reset, m.tr.Trans("welcome"), Green, Reset)
		fmt.Println("          Press <Enter> to continue...")
		fmt.Scanln()
		m.term.Clear()
	}

	// Header with app name and version.
	fmt.Printf("%v{%v Put In List %v}%v %v%s%v\n\n", Green, Reset, Green, Reset, Gray, version, Reset)

	// List management options.
	fmt.Printf("1 %v»%v %s\n", Green, Reset, m.tr.Trans("create_list"))
	fmt.Printf("2 %v»%v %s %v[%v %s: %s %v]%v\n", Green, Reset, m.tr.Trans("set_list_active"), Green, Reset, m.tr.Trans("current"), activeList, Green, Reset)
	fmt.Printf("3 %v»%v %s\n", Green, Reset, m.tr.Trans("show_lists"))
	fmt.Printf("4 %v»%v %s\n", Green, Reset, m.tr.Trans("uninstall_list"))
	fmt.Printf("%v-------------------------%v\n", Green, Reset)

	// Package management options.
	fmt.Printf("5 %v»%v %s\n", Green, Reset, m.tr.Trans("install_packages"))
	fmt.Printf("6 %v»%v %s %v[%v %s: %s %v]%v\n", Green, Reset, m.tr.Trans("choose_installer"), Green, Reset, m.tr.Trans("current"), currentInstaller, Green, Reset)
	fmt.Printf("7 %v»%v %s\n", Green, Reset, m.tr.Trans("search_package"))
	fmt.Printf("%v-------------------------%v\n", Green, Reset)

	// Settings and exit options.
	fmt.Printf("8 %v»%v %s\n", Green, Reset, m.tr.Trans("change_language"))
	fmt.Printf("%v-------------------------%v\n", Green, Reset)
	fmt.Printf("9 %v»%v %s\n\n", Green, Reset, m.tr.Trans("exit"))
}

// ReadChoice prints a prompt and returns the single-character choice entered by the user.
func (m *Menu) ReadChoice() string {
	fmt.Print("Enter your choice: ")
	var choice string = ""
	fmt.Scanln(&choice)
	return choice
}

// Pause prints a translated "press enter" message and waits for the user to press Enter.
func (m *Menu) Pause() {
	fmt.Println(m.tr.Trans("press_enter"))
	fmt.Scanln()
}
