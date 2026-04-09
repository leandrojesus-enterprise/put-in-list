package app

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/leandrojesus-enterprise/put-in-list/internal/config"
	"github.com/leandrojesus-enterprise/put-in-list/internal/i18n"
	"github.com/leandrojesus-enterprise/put-in-list/internal/installer"
	"github.com/leandrojesus-enterprise/put-in-list/internal/menu"
	"github.com/leandrojesus-enterprise/put-in-list/internal/storage"
	"github.com/leandrojesus-enterprise/put-in-list/internal/terminal"
)

var supportedOS []string = []string{"windows", "linux"}

type App struct {
	cfgSvc      *config.JSONService
	listSvc     *storage.JSONStore
	instSvc     *installer.Service
	ui          *menu.Menu
	term        *terminal.Terminal
	tr          *i18n.Translator
	cfg         config.Config
	lists       storage.StoredLists
	firstRun    bool
	menuChoices map[string]func()
}

func New(
	cfgSvc *config.JSONService,
	listSvc *storage.JSONStore,
	instSvc *installer.Service,
	ui *menu.Menu,
	term *terminal.Terminal,
	tr *i18n.Translator,
) *App {
	return &App{
		cfgSvc:  cfgSvc,
		listSvc: listSvc,
		instSvc: instSvc,
		ui:      ui,
		term:    term,
		tr:      tr,
	}
}

func (a *App) Run() {
	a.init()
	for {
		a.term.Clear()
		a.ui.ShowMain(a.lists.ActiveList, a.cfg.CurrentInstaller, a.firstRun)
		a.firstRun = false
		var choice string = a.ui.ReadChoice()
		a.runChoice(choice)
		a.ui.Pause()
	}
}

func (a *App) init() {
	fmt.Println(a.tr.Trans("app_initializing"))

	var currentOS string = runtime.GOOS
	if !a.osSupported(currentOS) {
		fmt.Println(a.tr.Trans("os_not_supported"))
		a.ui.Pause()
		os.Exit(1)
	}

	a.cfgSvc.SetTranslator(a.tr)

	var cfg config.Config
	var isNew bool
	var err error
	cfg, isNew, err = a.cfgSvc.Load()
	if err != nil {
		fmt.Println(err)
		a.ui.Pause()
		os.Exit(1)
	}
	a.cfg = cfg
	if isNew {
		fmt.Println(a.tr.Trans("config_not_found"))
		a.ui.Pause()
		if err := a.cfgSvc.Save(a.cfg); err != nil {
			fmt.Println(err)
		}
		a.firstRun = true
	}

	var lists storage.StoredLists
	lists, isNew, err = a.listSvc.Load()
	if err != nil {
		fmt.Println(a.tr.Trans("lists_read_error"), err)
		a.ui.Pause()
		os.Exit(1)
	}
	a.lists = lists
	if isNew {
		fmt.Println(a.tr.Trans("lists_not_found"))
		a.ui.Pause()
		if err := a.listSvc.Save(a.lists); err != nil {
			fmt.Println(err)
		}
	}

	if a.cfg.Language == "" {
		a.cfg.Language = i18n.EN
	}
	a.tr.SetLang(a.cfg.Language)

	a.setupMenuChoices()
}

func (a *App) osSupported(os string) bool {
	for _, s := range supportedOS {
		if os == s {
			return true
		}
	}
	return false
}

func (a *App) runChoice(choice string) {
	if choice >= "1" && choice <= "9" {
		a.menuChoices[choice]()
	} else {
		fmt.Println(a.tr.Trans("invalid_choice"))
	}
}

func (a *App) setupMenuChoices() {
	a.menuChoices = map[string]func(){
		"1": a.createList,
		"2": a.setActiveList,
		"3": a.showLists,
		"4": a.uninstallList,
		"5": a.installPackage,
		"6": a.chooseInstaller,
		"7": a.searchPackage,
		"8": a.changeLanguage,
		"9": a.exitApp,
	}
}

func (a *App) createList() {
	fmt.Print(a.tr.Trans("enter_list_name"))
	var listName string
	fmt.Scanln(&listName)

	if strings.TrimSpace(listName) == "" {
		fmt.Println(a.tr.Trans("list_name_cannot_be_empty"))
		return
	}

	if _, exists := a.lists.Lists[listName]; exists {
		fmt.Printf(a.tr.Trans("list_already_exists"), listName)
		return
	}

	a.lists.Lists[listName] = []storage.InstallEntry{}
	a.listSvc.Save(a.lists)
}

func (a *App) setActiveList() {
	fmt.Println(a.tr.Trans("available_lists"))
	i := 1
	for listName := range a.lists.Lists {
		fmt.Printf(" %d - %s\n", i, listName)
		i++
	}

	fmt.Print(a.tr.Trans("enter_list_index"))
	var listIndex int
	fmt.Scanln(&listIndex)

	if listIndex < 1 || listIndex > len(a.lists.Lists) {
		fmt.Println(a.tr.Trans("invalid_list_index"))
		return
	}

	var listNames []string = make([]string, 0, len(a.lists.Lists))
	for name := range a.lists.Lists {
		listNames = append(listNames, name)
	}

	var listName string = listNames[listIndex-1]
	if _, exists := a.lists.Lists[listName]; exists {
		a.lists.ActiveList = listName
		a.listSvc.Save(a.lists)
	} else {
		fmt.Printf(a.tr.Trans("list_does_not_exist"), listName)
	}
}

func (a *App) showLists() {
	fmt.Println(a.tr.Trans("available_lists"))
	for listName, entries := range a.lists.Lists {
		var activeMarker string = ""
		if listName == a.lists.ActiveList {
			activeMarker = a.tr.Trans("active_marker")
		}
		fmt.Printf(" - %s%s\n", listName, activeMarker)
		for _, entry := range entries {
			fmt.Printf("   - %s (%s)\n", entry.Name, entry.Installer)
		}
	}
}

func (a *App) uninstallList() {
	if a.lists.ActiveList == a.tr.Trans("none") {
		fmt.Println(a.tr.Trans("no_active_list"))
		return
	}

	var entries []storage.InstallEntry
	var ok bool
	entries, ok = a.lists.Lists[a.lists.ActiveList]
	if !ok {
		fmt.Printf(a.tr.Trans("list_does_not_exist"), a.lists.ActiveList)
		return
	}

	if len(entries) == 0 {
		fmt.Printf(a.tr.Trans("list_is_empty"), a.lists.ActiveList)
		fmt.Print(a.tr.Trans("delete_list_prompt"))
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) == "y" {
			delete(a.lists.Lists, a.lists.ActiveList)
			a.lists.ActiveList = a.tr.Trans("none")
			a.listSvc.Save(a.lists)
		}
		return
	}

	fmt.Print(a.tr.Trans("uninstall_list_prompt"))
	var response string
	fmt.Scanln(&response)
	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		return
	}

	fmt.Print(a.tr.Trans("delete_list_prompt"))
	fmt.Scanln(&response)
	if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
		delete(a.lists.Lists, a.lists.ActiveList)
		a.lists.ActiveList = a.tr.Trans("none")
		a.listSvc.Save(a.lists)
	}

	for i := len(entries) - 1; i >= 0; i-- {
		var entry storage.InstallEntry = entries[i]
		fmt.Printf(a.tr.Trans("uninstalling"), entry.Name, entry.Installer)

		if err := a.instSvc.Uninstall(entry, a.tr); err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf(a.tr.Trans("successfully_uninstalled"), entry.Name)
		entries = append(entries[:i], entries[i+1:]...)
	}

	a.lists.Lists[a.lists.ActiveList] = entries
	a.listSvc.Save(a.lists)

	if len(entries) == 0 {
		fmt.Printf(a.tr.Trans("list_is_empty"), a.lists.ActiveList)
	}
}

func (a *App) installPackage() {
	if a.cfg.CurrentInstaller == "None" {
		fmt.Println(a.tr.Trans("no_installer_selected"))
		return
	}

	if a.lists.ActiveList == "None" {
		fmt.Println(a.tr.Trans("no_active_list"))
		return
	}

	fmt.Print(a.tr.Trans("package_name_prompt"))
	var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()
	var pkg string = scanner.Text()

	if strings.TrimSpace(pkg) == "" {
		fmt.Println(a.tr.Trans("package_name_empty"))
		return
	}

	var pkgSplit []string = strings.Split(pkg, " ")
	var pkgName string = pkgSplit[0]

	var pkgArg string = ""
	if len(pkgSplit) >= 2 {
		pkgArg = pkgSplit[1]
	}

	fmt.Printf(a.tr.Trans("installing"), pkgName, a.cfg.CurrentInstaller)

	var installCmd *exec.Cmd = a.buildInstallCmd(pkgName, pkgArg)
	if installCmd == nil {
		fmt.Printf("Installer '%s' not supported.\n", a.cfg.CurrentInstaller)
		return
	}

	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	if err := installCmd.Run(); err != nil {
		fmt.Printf(a.tr.Trans("failed_to_install"), pkgName, err)
		fmt.Println(a.tr.Trans("wont_be_added_to_active_list"))
		return
	}

	fmt.Printf(a.tr.Trans("successfully_installed"), pkgName)

	var activeEntries []storage.InstallEntry = a.lists.Lists[a.lists.ActiveList]
	for _, existing := range activeEntries {
		if existing.Name == pkg && existing.Installer == a.cfg.CurrentInstaller {
			fmt.Printf(a.tr.Trans("package_already_exists"), pkgName, a.lists.ActiveList, a.cfg.CurrentInstaller)
			return
		}
	}

	a.lists.Lists[a.lists.ActiveList] = append(activeEntries, storage.InstallEntry{Name: pkgName, Installer: a.cfg.CurrentInstaller})
	a.listSvc.Save(a.lists)
	fmt.Printf(a.tr.Trans("installed_added_to_list"), pkgName, a.cfg.CurrentInstaller, a.lists.ActiveList)
}

func (a *App) buildInstallCmd(pkgName, pkgArg string) *exec.Cmd {
	switch a.cfg.CurrentInstaller {
	case "winget":
		if pkgArg != "" {
			return exec.Command("winget", "install", pkgName, pkgArg)
		}
		return exec.Command("winget", "install", pkgName)
	case "snap":
		if pkgArg == "--classic" {
			return exec.Command("snap", "install", pkgName, "--classic")
		}
		return exec.Command("snap", "install", pkgName)
	case "apt":
		if pkgArg != "" {
			return exec.Command("sudo", "apt", "install", "-y", pkgArg, pkgName)
		}
		return exec.Command("sudo", "apt", "install", "-y", pkgName)
	case "choco":
		if pkgArg != "" {
			return exec.Command("choco", "install", pkgName, pkgArg, "-y")
		}
		return exec.Command("choco", "install", pkgName, "-y")
	default:
		return nil
	}
}

func (a *App) chooseInstaller() {
	var currentOS string = runtime.GOOS
	var supported []string = installer.SupportedInstallers[currentOS]
	var available []string = a.instSvc.DetectAvailable(currentOS)

	fmt.Printf(a.tr.Trans("supported_installers"), currentOS, strings.Join(supported, ", "))
	if len(available) == 0 {
		fmt.Println(a.tr.Trans("no_available_installers"))
		fmt.Println(a.tr.Trans("install_installer_instructions"))
		return
	}

	fmt.Println(a.tr.Trans("available_installers"))
	for i, tool := range available {
		var currentMark string = ""
		if tool == a.cfg.CurrentInstaller {
			currentMark = " (current)"
		}
		fmt.Printf("%d) %s%s\n", i+1, tool, currentMark)
	}

	fmt.Print(a.tr.Trans("choose_installer_prompt"))
	var opt int
	var err error
	_, err = fmt.Scanln(&opt)
	if err != nil || opt < 1 || opt > len(available) {
		fmt.Println(a.tr.Trans("invalid_index"))
		return
	}

	a.cfg.CurrentInstaller = available[opt-1]
	a.cfgSvc.Save(a.cfg)
	fmt.Printf(a.tr.Trans("installer_set_current"), a.cfg.CurrentInstaller)
}

func (a *App) searchPackage() {
	fmt.Print(a.tr.Trans("package_search_prompt"))
	var packageName string
	fmt.Scanln(&packageName)
	packageName = strings.ToLower(packageName)

	var found bool = false
	for listName, entries := range a.lists.Lists {
		for _, entry := range entries {
			if strings.ToLower(entry.Name) == packageName {
				fmt.Printf(a.tr.Trans("package_found_in_list"), packageName, listName, entry.Installer)
				found = true
				fmt.Print(a.tr.Trans("do_you_want_to_uninstall"))
				var confirm string
				fmt.Scanln(&confirm)
				if strings.ToLower(confirm) == "y" {
					a.instSvc.Uninstall(entry, a.tr)

					var updatedEntries []storage.InstallEntry = []storage.InstallEntry{}
					for _, e := range entries {
						if !(strings.ToLower(e.Name) == packageName && e.Installer == entry.Installer) {
							updatedEntries = append(updatedEntries, e)
						}
					}
					a.lists.Lists[listName] = updatedEntries
					a.listSvc.Save(a.lists)
					fmt.Printf(a.tr.Trans("package_removed_from_list"), packageName, listName)

					if len(updatedEntries) == 0 {
						fmt.Printf(a.tr.Trans("list_empty_prompt"), listName)
						var delConfirm string
						fmt.Scanln(&delConfirm)
						if strings.ToLower(delConfirm) == "y" {
							delete(a.lists.Lists, listName)
							a.listSvc.Save(a.lists)
							fmt.Printf(a.tr.Trans("list_deleted"), listName)
						}
					}
					break
				}
			}
		}
	}

	if !found {
		fmt.Printf(a.tr.Trans("package_not_found_in_any_list"), packageName)
	}
}

func (a *App) changeLanguage() {
	fmt.Println("1 - English")
	fmt.Println("2 - Português")
	fmt.Println("3 - Español")

	var opt int
	fmt.Scanln(&opt)

	switch opt {
	case 1:
		a.tr.SetLang(i18n.EN)
	case 2:
		a.tr.SetLang(i18n.PT)
	case 3:
		a.tr.SetLang(i18n.ES)
	default:
		fmt.Println(a.tr.Trans("invalid_option"))
		return
	}

	a.cfg.Language = a.tr.Lang()
	a.cfgSvc.Save(a.cfg)
	fmt.Printf(a.tr.Trans("language_changed"), a.tr.Lang())
}

func (a *App) exitApp() {
	fmt.Print(a.tr.Trans("exit_app"))
	os.Exit(0)
}
