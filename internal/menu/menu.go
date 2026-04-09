package menu

import (
	"fmt"

	"github.com/leandrojesus-enterprise/put-in-list/internal/i18n"
	"github.com/leandrojesus-enterprise/put-in-list/internal/terminal"
)

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

type Menu struct {
	term *terminal.Terminal
	tr   *i18n.Translator
}

func New(term *terminal.Terminal, tr *i18n.Translator) *Menu {
	return &Menu{term: term, tr: tr}
}

func (m *Menu) ShowMain(activeList, currentInstaller string, firstRun bool) {
	if firstRun {
		fmt.Printf("           %v{%v %s %v}%v\n\n", Green, Reset, m.tr.Trans("welcome"), Green, Reset)
		fmt.Println("          Press <Enter> to continue...")
		fmt.Scanln()
		m.term.Clear()
	}

	fmt.Printf("%v{%v Put In List %v}%v\n\n", Green, Reset, Green, Reset)
	fmt.Printf("1 %v»%v %s\n", Green, Reset, m.tr.Trans("create_list"))
	fmt.Printf("2 %v»%v %s %v[%v Current: %s %v]%v\n", Green, Reset, m.tr.Trans("set_list_active"), Green, Reset, activeList, Green, Reset)
	fmt.Printf("3 %v»%v %s\n", Green, Reset, m.tr.Trans("show_lists"))
	fmt.Printf("4 %v»%v %s\n", Green, Reset, m.tr.Trans("uninstall_list"))
	fmt.Printf("%v-------------------------%v\n", Green, Reset)
	fmt.Printf("5 %v»%v %s\n", Green, Reset, m.tr.Trans("install_packages"))
	fmt.Printf("6 %v»%v %s %v[%v Current: %s %v]%v\n", Green, Reset, m.tr.Trans("choose_installer"), Green, Reset, currentInstaller, Green, Reset)
	fmt.Printf("7 %v»%v %s\n", Green, Reset, m.tr.Trans("search_package"))
	fmt.Printf("%v-------------------------%v\n", Green, Reset)
	fmt.Printf("8 %v»%v %s\n", Green, Reset, m.tr.Trans("change_language"))
	fmt.Printf("%v-------------------------%v\n", Green, Reset)
	fmt.Printf("9 %v»%v %s\n\n", Green, Reset, m.tr.Trans("exit"))
}

func (m *Menu) ReadChoice() string {
	fmt.Print("Enter your choice: ")
	var choice string
	fmt.Scanln(&choice)
	return choice
}

func (m *Menu) Pause() {
	fmt.Println(m.tr.Trans("press_enter"))
	fmt.Scanln()
}
