package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ---------- I18N ----------
type Lang string

const (
	EN Lang = "en"
	PT Lang = "pt"
	ES Lang = "es"
)

var currentLang = EN

var translations = map[Lang]map[string]string{
	EN: {
		"app_initializing":   "** App initializing **",
		"os_not_supported": "Your operating system is not supported. Supported OS: windows, linux.",
		"config_not_found":  "config.json not found! Will create one with default values.",
		"config_read_error": "Error reading config.json:",
		"lists_not_found":   "lists.json not found! Will create one with default values.",
		"lists_read_error":  "Error reading lists.json:",
		"save_config_error": "Error saving config.json:",
		"save_config_error_marshal": "Error marshaling config.json:",
		"save_config_error_write": "Error writing config.json:",
		"create_list":       "Create List",
		"set_list_active":   "Set Active List",
		"show_lists":        "Show Lists",
		"uninstall_list":    "Uninstall List",
		"install_packages":  "Install Packages",
		"choose_installer":  "Choose Installer",
		"search_package":    "Search Package in Lists",
		"change_language":   "Change Language",
		"exit":              "Exit",
		"list_already_exists": "List '%s' already exists!\n",
		"enter_list_name": "Enter the name of the new list: ",
		"list_name_cannot_be_empty": "List name cannot be empty!",
		"available_lists": "Available lists:",
		"invalid_list_index": "Invalid list index!",
		"enter_list_index": "Enter the index of the list to set as active: ",
		"list_does_not_exist": "List '%s' does not exist!\n",
		"active_marker": " (active)",
		"none":              "None",
		"list_is_empty": "List '%s' is empty.",
		"delete_list_prompt": "Do you want to delete the list? (y/n)",
		"uninstall_list_prompt": "Do you really want to uninstall this list? this will uninstall all packages in the list. (y/n)",
		"uninstalling": "Uninstalling '%s' with '%s'...",
		"error_uninstalling": "Failed to uninstall '%s': %v\n",
		"successfully_uninstalled": "'%s' uninstalled successfully.\n",
		"no_installer_selected": "Error: no installer selected. Use option 6.",
		"no_active_list": "Error: no active list. Use option 2.",
		"package_name_prompt": "Enter the name of the package to install: ",
		"package_name_empty": "Package name cannot be empty.",
		"installing": "Installing '%s' with %s...\n",
		"failed_to_install": "Failed to install '%s': %v\n",
		"wont_be_added_to_active_list": "Will not be added to active list.",
		"successfully_installed": "'%s' installed successfully.\n",
		"installed_added_to_list": "'%s' [%s] added to active list '%s'.\n",
		"package_already_exists": "'%s' already exists in active list '%s' with installer '%s'.\n",
		"supported_installers": "Supported installers for %s: %s\n",
		"no_available_installers": "No available installers found in PATH.",
		"install_installer_instructions": "To use this option, please install one of them or add it to your PATH.",
		"available_installers": "Available installers detected:",
		"choose_installer_prompt": "Choose the number of the installer to set as current: ",
		"invalid_index": "Invalid index! Choose a valid installer number.",
		"installer_set_current": "Current installer set to: %s\n",
		"package_search_prompt": "Enter the name of the package you want to search for in all lists: ",
		"package_found_in_list": "Package '%s' found in list '%s' (installer=%s)\n",
		"do_you_want_to_uninstall": "Do you want to uninstall it? (y/n): ",
		"package_not_found_in_any_list": "Package '%s' not found in any list.\n",
		"list_empty_prompt": "List '%s' is now empty. Do you want to delete it? (y/n): ",
		"package_removed_from_list": "Package '%s' removed from list '%s'.\n",
		"list_deleted": "List '%s' deleted.\n",
		"language_changed": "Language changed to %s.\n",
		"uninstalling_package": "Uninstalling '%s' with '%s'...\n",
		"instlaller_not_supported": "Installer '%s' not supported for package '%s'.\n",
		"uninstall_failed": "Failed to uninstall '%s': %v\n",
		"uninstalled_successfully": "'%s' uninstalled successfully.\n",
		"exit_app": "Exiting the app!\nSee you soon.\n",
		"welcome":           "Welcome to Put In List",
		"enter_choice":      "Enter your choice:",
		"invalid_choice":    "Invalid choice! Enter a number between 1-9.",
		"press_enter":       "Press <Enter> to continue.",
		"installer_not_set": "No installer selected. Use option 6.",

	},
	PT: {
		"app_initializing":   "** Aplicação inicializando **",
		"os_not_supported": "O seu sistema operacional não é suportado. Sistemas operacionais suportados: windows, linux.", 
		"config_not_found":  "config.json não foi encontrado! Será criado um com valores padrão.",
		"config_read_error": "Erro ao ler config.json:",
		"lists_not_found":   "lists.json não foi encontrado! Será criado um com valores padrão.",
		"lists_read_error":  "Erro ao ler lists.json:",
		"save_config_error": "Erro ao salvar config.json:",
		"save_config_error_marshal": "Erro ao serializar config.json:",
		"save_config_error_write": "Erro ao escrever config.json:",
		"create_list":       "Criar Lista",
		"set_list_active":   "Definir Lista Ativa",
		"show_lists":        "Mostrar Listas",
		"uninstall_list":    "Desinstalar Lista",
		"install_packages":  "Instalar Pacotes",
		"choose_installer":  "Escolher Instalador",
		"search_package":    "Pesquisar Pacote em Listas",
		"change_language":   "Alterar Idioma",
		"exit":              "Sair",
		"list_already_exists": "A lista '%s' já existe!\n",
		"enter_list_name": "Insira o nome da nova lista: ",
		"list_name_cannot_be_empty": "O nome da lista não pode ser vazio!",
		"available_lists": "Listas disponíveis:",
		"invalid_list_index": "Índice de lista inválido!",
		"enter_list_index": "Digite o índice da lista para definir como ativa: ",
		"list_does_not_exist": "A lista '%s' não existe!\n",
		"active_marker": " (ativa)",
		"none":              "None",
		"list_is_empty": "A lista '%s' está vazia.",
		"delete_list_prompt": "Você deseja excluir a lista? (y/n)",
		"uninstall_list_prompt": "Você realmente deseja desinstalar esta lista? Isso irá desinstalar todos os pacotes na lista. (y/n)",
		"uninstalling": "Desinstalando '%s' com '%s'...",
		"error_uninstalling": "Falha ao desinstalar '%s': %v\n",
		"successfully_uninstalled": "'%s' desinstalado com sucesso.\n",
		"no_installer_selected": "Erro: Nenhum instalador selecionado. Usar opção 6.",
		"no_active_list": "Erro: Nenhuma lista ativa. Usar opção 2.",
		"package_name_prompt": "Insira o nome do pacote para instalar: ",
		"package_name_empty": "O nome do pacote não pode ser vazio.",
		"installing": "Instalando '%s' com %s...\n",
		"failed_to_install": "Falha ao instalar '%s': %v\n",
		"wont_be_added_to_active_list": "Não será adicionado à lista ativa.",
		"successfully_installed": "'%s' instalado com sucesso.\n",
		"installed_added_to_list": "'%s' [%s] adicionado à lista ativa '%s'.\n",
		"package_already_exists": "'%s' já existe na lista ativa '%s' com instalador '%s'.\n",
		"supported_installers": "Installers suportados para %s: %s\n",
		"no_available_installers": "Nenhum instalador disponível encontrado no PATH.",
		"install_installer_instructions": "Para usar esta opção, por favor instale um deles ou adicione ao seu PATH.",
		"available_installers": "Installers disponíveis detectados:",
		"choose_installer_prompt": "Escolha o número do instalador para definir como atual: ",
		"invalid_index": "Índice inválido! Escolha um número de instalador válido.",
		"installer_set_current": "Instalador atual definido para: %s\n",
		"package_search_prompt": "Insira o nome do pacote que você deseja pesquisar em todas as listas: ",
		"package_found_in_list": "Pacote '%s' encontrado na lista '%s' (installer=%s)\n",
		"do_you_want_to_uninstall": "Você deseja desinstalá-lo? (y/n): ",
		"package_not_found_in_any_list": "Pacote '%s' não encontrado em nenhuma lista.\n",
		"list_empty_prompt": "Lista '%s' está agora vazia. Você deseja excluí-la? (y/n): ",
		"package_removed_from_list": "Pacote '%s' removido da lista '%s'.\n",
		"list_deleted": "Lista '%s' excluída.\n",
		"language_changed": "Idioma alterado para %s.\n",
		"uninstalling_package": "Desinstalando '%s' com '%s'...\n",
		"instlaller_not_supported": "Installer '%s' not supported for package '%s'.\n",
		"uninstall_failed": "Falha ao desinstalar '%s': %v\n",
		"uninstalled_successfully": "'%s' desinstalado com sucesso.\n",
		"exit_app": "Saindo do aplicativo!\nAté mais.\n",
		"welcome":           "Bem-vindo ao Put In List",
		"enter_choice":      "Digite sua escolha:",
		"invalid_choice":    "Escolha inválida! Digite um número entre 1-9.",
		"press_enter":       "Pressione <Enter> para continuar.",
		"installer_not_set": "Nenhum instalador selecionado. Use a opção 6.",
	},
	ES: {
		"welcome":           "Bienvenido a Put In List",
		"enter_choice":      "Elige una opción:",
		"invalid_choice":    "Opción inválida. Elige entre 1-8.",
		"press_enter":       "Presiona <Enter> para continuar.",
		"no_active_list":    "Ninguna lista activa. Usa la opción 2.",
		"installer_not_set": "Ningún instalador seleccionado. Usa la opción 6.",
	},
}

func trans(key string) string {
	return translations[currentLang][key]
}

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

type Config struct {
	CurrentInstaller string `json:"currentInstaller"`
	Language         Lang   `json:"language"`
}

var config Config

var supportedInstallers = map[string][]string{
	"windows": {"winget", "choco"},
	"linux":   {"apt", "snap"},
}

type InstallEntry struct {
	Name      string `json:"name"`
	Installer string `json:"installer"`
}

type StoredLists struct {
	ActiveList string                    `json:"activeList"`
	Lists      map[string][]InstallEntry `json:"lists"`
}

var loadedLists StoredLists

var supportedOS = []string{"windows", "linux"}
var currentOS string

var clearTerminalFunc map[string]func()

var choice string

var storedLists StoredLists

var mainMenuChoices map[string]func()

func main() {
	_init()
	for {
		clearTerminal()

		showMainMenu()
		readMainMenuChoice()

		runMainMenuChoice()
		pause()
	}
}

func _init() {
	fmt.Println(trans("app_initializing"))

	currentOS = runtime.GOOS
	if !currentOSSupported() {
		fmt.Println(trans("os_not_supported"))
		pause()
		exit(1)
	}

	config = Config{CurrentInstaller: "None"}

	cfgFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(trans("config_not_found"))
		pause()
		saveConfig(config)
	} else {
		byteValue, err := io.ReadAll(cfgFile)
		if err != nil {
			fmt.Println(trans("config_read_error"), err)
			pause()
			exit(1)
		}
		cfgFile.Close()
		json.Unmarshal(byteValue, &config)
	}

	storedLists = StoredLists{}
	storedLists.Lists = make(map[string][]InstallEntry)

	//Check if lists.json is created
	jsonFile, err := os.Open("lists.json")
	if err != nil {
		fmt.Println(trans("lists_not_found"))
		pause()

		storedLists = StoredLists{
			ActiveList: "None",
			Lists:      map[string][]InstallEntry{},
		}

		saveListsJson(storedLists)
	} else {
		// Read opened file as a byte array
		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			fmt.Println(trans("lists_read_error"), err)
			pause()
			exit(1)
		}

		// Unmarshal JSON data into our StoredLists type
		json.Unmarshal(byteValue, &loadedLists)
		jsonFile.Close()

		storedLists.ActiveList = loadedLists.ActiveList
		storedLists.Lists = loadedLists.Lists
	}

	setupClearTerminal()
	setupMainMenuChoices()

	if config.Language == "" {
		config.Language = EN
	}
	
	currentLang = config.Language
}

func saveListsJson(data StoredLists) {
	// Convert data structure to JSON
	byteValue, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		pause()
		exit(1)
		return
	}

	// Write to JSON file
	err = os.WriteFile("lists.json", byteValue, 0644)
	if err != nil {
		fmt.Println(err)
		pause()
		exit(1)
	}
}

func saveConfig(data Config) {
	byteValue, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(trans("save_config_error_marshal"), err)
		pause()
		exit(1)
	}
	err = os.WriteFile("config.json", byteValue, 0644)
	if err != nil {
		fmt.Println(trans("save_config_error_write"), err)
		pause()
		exit(1)
	}
}

func clearTerminal() {
	currentOS := runtime.GOOS
	clearTerminalFunc[currentOS]()
}

func showMainMenu() {
	fmt.Printf("           %v{%v %s %v}%v\n\n", Green, Reset, trans("welcome"), Green, Reset)
	fmt.Printf("1 %v»%v %s\n", Green, Reset, trans("create_list"))
	fmt.Printf("2 %v»%v %s %v[%v Current: %s %v]%v\n", Green, Reset, trans("set_list_active"), Green, Reset, storedLists.ActiveList, Green, Reset)
	fmt.Printf("3 %v»%v %s\n", Green, Reset, trans("show_lists"))
	fmt.Printf("4 %v»%v %s\n", Green, Reset, trans("uninstall_list"))
	fmt.Printf("%v-------------------------%v\n", Green, Reset)
	fmt.Printf("5 %v»%v %s\n", Green, Reset, trans("install_packages"))
	fmt.Printf("6 %v»%v %s %v[%v Current: %s %v]%v\n", Green, Reset, trans("choose_installer"), Green, Reset, config.CurrentInstaller, Green, Reset)
	fmt.Printf("7 %v»%v %s\n", Green, Reset, trans("search_package"))
	fmt.Printf("%v-------------------------%v\n", Green, Reset)
	fmt.Printf("8 %v»%v %s\n", Green, Reset, trans("change_language"))
	fmt.Printf("%v-------------------------%v\n", Green, Reset)
	fmt.Printf("9 %v»%v %s\n\n", Green, Reset, trans("exit"))
}

func readMainMenuChoice() {
	fmt.Print("Enter your choice: ")
	fmt.Scanln(&choice)
}

func runMainMenuChoice() {
	if choice >= "1" && choice <= "9" {
		mainMenuChoices[choice]()
	} else {
		fmt.Println(trans("invalid_choice"))
	}
}

func pause() {
	fmt.Println(trans("press_enter"))
	fmt.Scanln()
}

func currentOSSupported() bool {
	var isSupported bool = false
	for _, os := range supportedOS {
		if currentOS == os {
			isSupported = true
			break
		}
	}

	return isSupported
}

func detectAvailableInstallers(osName string) []string {
	available := []string{}
	supported, ok := supportedInstallers[osName]
	if !ok {
		return available
	}

	for _, tool := range supported {
		if _, err := exec.LookPath(tool); err == nil {
			available = append(available, tool)
		}
	}
	return available
}

func setupClearTerminal() {
	clearTerminalFunc = make(map[string]func())

	clearTerminalFunc["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	clearTerminalFunc["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func setupMainMenuChoices() {
	mainMenuChoices = make(map[string]func())
	mainMenuChoices["1"] = func() {
		fmt.Print(trans("enter_list_name"))
		var listName string
		fmt.Scanln(&listName)

		if strings.TrimSpace(listName) == "" {
			fmt.Println(trans("list_name_cannot_be_empty"))
			return
		}
		
		if _, exists := storedLists.Lists[listName]; exists {
			fmt.Printf(trans("list_already_exists"), listName)
			return
		}

		storedLists.Lists[listName] = []InstallEntry{}

		saveListsJson(storedLists)
	}

	mainMenuChoices["2"] = func() {
		fmt.Println(trans("available_lists"))
		i := 1
		for listName := range storedLists.Lists {
			fmt.Printf(" %d - %s\n", i, listName)
			i++
		}

		fmt.Print(trans("enter_list_index"))
		var listIndex int
		fmt.Scanln(&listIndex)

		if listIndex < 1 || listIndex > len(storedLists.Lists) {
			fmt.Println(trans("invalid_list_index"))
			return
		}

		listNames := make([]string, 0, len(storedLists.Lists))
		for name := range storedLists.Lists {
			listNames = append(listNames, name)
		}

		listName := listNames[listIndex-1]
		if _, exists := storedLists.Lists[listName]; exists {
			storedLists.ActiveList = listName
			saveListsJson(storedLists)
		} else {
			fmt.Printf(trans("list_does_not_exist"), listName)
		}
	}

	mainMenuChoices["3"] = func() {
		fmt.Println(trans("available_lists"))
		for listName, entries := range storedLists.Lists {
			activeMarker := ""
			if listName == storedLists.ActiveList {
				activeMarker = trans("active_marker")
			}
			fmt.Printf(" - %s%s\n", listName, activeMarker)
			for _, entry := range entries {
				fmt.Printf("   - %s (%s)\n", entry.Name, entry.Installer)
			}
		}
	}

	mainMenuChoices["4"] = func() {
		if storedLists.ActiveList == trans("none") {
			fmt.Println(trans("no_active_list"))
			return
		}

		entries, ok := storedLists.Lists[storedLists.ActiveList]
		if !ok {
			fmt.Printf(trans("list_does_not_exist"), storedLists.ActiveList)
			return
		}

				if len(entries) == 0 {
			fmt.Printf(trans("list_is_empty"), storedLists.ActiveList)
			fmt.Print(trans("delete_list_prompt"))
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) == "y" {
				delete(storedLists.Lists, storedLists.ActiveList)
				storedLists.ActiveList = trans("none")
				saveListsJson(storedLists)
			}
			return
		} else {
			fmt.Print(trans("uninstall_list_prompt"))
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
				return
			}

			fmt.Print(trans("delete_list_prompt"))
			fmt.Scanln(&response)
			if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
				delete(storedLists.Lists, storedLists.ActiveList)
				storedLists.ActiveList = trans("none")
				saveListsJson(storedLists)
			}
			
		}

		for i := len(entries) - 1; i >= 0; i-- {
			entry := entries[i]
			fmt.Printf(trans("uninstalling"), entry.Name, entry.Installer)

			var cmd *exec.Cmd
			switch entry.Installer {
			case "winget":
				cmd = exec.Command("winget", "uninstall", entry.Name)
			case "choco":
				cmd = exec.Command("choco", "uninstall", entry.Name, "-y")
			case "apt":
				cmd = exec.Command("sudo", "apt", "remove", "-y", entry.Name)
			case "snap":
				cmd = exec.Command("snap", "remove", entry.Name)
			default:
				fmt.Printf("Instalador '%s' não suportado para o pacote '%s'.\n", entry.Installer, entry.Name)
				continue
			}

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				fmt.Printf(trans("error_uninstalling"), entry.Name, err)
				continue
			}

			fmt.Printf(trans("successfully_uninstalled"), entry.Name)

			entries = append(entries[:i], entries[i+1:]...)
		}

		storedLists.Lists[storedLists.ActiveList] = entries
		saveListsJson(storedLists)

		if len(entries) == 0 {
			fmt.Printf(trans("list_is_empty"), storedLists.ActiveList)
		}
	}

	mainMenuChoices["5"] = func() {
		if config.CurrentInstaller == trans("none") {
			fmt.Println(trans("no_installer_selected"))
			return
		}

		if storedLists.ActiveList == trans("none") {
			fmt.Println(trans("no_active_list"))
			return
		}

		fmt.Print(trans("package_name_prompt"))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan() // use `for scanner.Scan()` to keep reading
		pkg := scanner.Text()

		if strings.TrimSpace(pkg) == "" {
			fmt.Println(trans("package_name_empty"))
			return
		}
		
		var classicInstallFlag = false
		pkgSplit := strings.Split(pkg, " ")
		pkgName := pkgSplit[0]
		
		var pkgArg string = ""
		if len(pkgSplit) >= 2 {
			pkgArg = pkgSplit[1]
		}

		if len(pkgSplit) == 2 {
			if pkgArg == "--classic"{
				classicInstallFlag = true
			}
		}

		fmt.Printf(trans("installing"), pkgName, config.CurrentInstaller)

		var cmd *exec.Cmd
		switch config.CurrentInstaller {
		case "winget":
			cmd = exec.Command("winget", "install", pkg)
		case "snap":
			if classicInstallFlag {
				cmd = exec.Command("snap", "install", pkgName, "--classic")
			} else {
				cmd = exec.Command("snap", "install", pkgName)
			}
		default:
			fmt.Printf("Instalador '%s' não suportado.\n", config.CurrentInstaller)
			return
		}

		// output em tempo real
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf(trans("failed_to_install"), pkgName, err)
			fmt.Println(trans("wont_be_added_to_active_list"))
			return
		}

		fmt.Printf(trans("successfully_installed"), pkgName)

		activeEntries := storedLists.Lists[storedLists.ActiveList]
		for _, existing := range activeEntries {
			if existing.Name == pkg && existing.Installer == config.CurrentInstaller {
				fmt.Printf(trans("package_already_exists"), pkgName, storedLists.ActiveList, config.CurrentInstaller)
				return
			}
		}

		storedLists.Lists[storedLists.ActiveList] = append(activeEntries, InstallEntry{Name: pkgName, Installer: config.CurrentInstaller})
		saveListsJson(storedLists)
		fmt.Printf(trans("installed_added_to_list"), pkgName, config.CurrentInstaller, storedLists.ActiveList)
	}

	mainMenuChoices["6"] = func() {
		supported := supportedInstallers[currentOS]
		available := detectAvailableInstallers(currentOS)

		fmt.Printf(trans("supported_installers"), currentOS, strings.Join(supported, ", "))
		if len(available) == 0 {
			fmt.Println(trans("no_available_installers"))
			fmt.Println(trans("install_installer_instructions"))
			return
		}

		fmt.Println(trans("available_installers"))
		for i, tool := range available {
			currentMark := ""
			if tool == config.CurrentInstaller {
				currentMark = " (current)"
			}
			fmt.Printf("%d) %s%s\n", i+1, tool, currentMark)
		}

		fmt.Print(trans("choose_installer_prompt"))
		var opt int
		_, err := fmt.Scanln(&opt)
		if err != nil || opt < 1 || opt > len(available) {
			fmt.Println(trans("invalid_index"))
			return
		}

		selected := available[opt-1]
		config.CurrentInstaller = selected
		saveConfig(config)

		fmt.Printf(trans("installer_set_current"), selected)
	}

	mainMenuChoices["7"] = func() {
		fmt.Print(trans("package_search_prompt"))
		var packageName string
		fmt.Scanln(&packageName)
		packageName = strings.ToLower(packageName)

		found := false
		for listName, entries := range storedLists.Lists {
			for _, entry := range entries {
				if strings.ToLower(entry.Name) == packageName {
					fmt.Printf(trans("package_found_in_list"), packageName, listName, entry.Installer)
					found = true
					fmt.Print(trans("do_you_want_to_uninstall"))
					var confirm string
					fmt.Scanln(&confirm)
					if strings.ToLower(confirm) == "y" {
						uninstallPackageIndividually(entry)
						
						// Remove the package from the list
						updatedEntries := []InstallEntry{}
						for _, e := range entries {
							if !(strings.ToLower(e.Name) == packageName && e.Installer == entry.Installer) {
								updatedEntries = append(updatedEntries, e)
							}
						}
						storedLists.Lists[listName] = updatedEntries
						saveListsJson(storedLists)
						fmt.Printf(trans("package_removed_from_list"), packageName, listName)
						
						//If list is empty after removal, ask if user wants to delete the list
						if len(updatedEntries) == 0 {
							fmt.Printf(trans("list_empty_prompt"), listName)
							var delConfirm string
							fmt.Scanln(&delConfirm)
							if strings.ToLower(delConfirm) == "y" {
								delete(storedLists.Lists, listName)
								saveListsJson(storedLists)
								fmt.Printf(trans("list_deleted"), listName)
							}
						}
						break
					}
				}
			}
		}

		if !found {
			fmt.Printf(trans("package_not_found_in_any_list"), packageName)
		}
	}

	mainMenuChoices["8"] = func() {
		fmt.Println("1 - English")
		fmt.Println("2 - Português")
		fmt.Println("3 - Español")

		var opt int
		fmt.Scanln(&opt)

		switch opt {
			case 1:
				currentLang = EN
			case 2:
				currentLang = PT
			case 3:
				currentLang = ES
			default:
				fmt.Println(trans("invalid_option"))
				return
		}

		config.Language = currentLang
		saveConfig(config)

		fmt.Printf(trans("language_changed"), currentLang)
	}
}

func uninstallPackageIndividually(entry InstallEntry) {
	fmt.Printf(trans("uninstalling_package"), entry.Name, entry.Installer)
	
	var cmd *exec.Cmd
	switch entry.Installer {
	case "winget":
		cmd = exec.Command("winget", "uninstall", entry.Name)
	case "choco":
		cmd = exec.Command("choco", "uninstall", entry.Name, "-y")
	case "apt":
		cmd = exec.Command("sudo", "apt", "remove", "-y", entry.Name)
	case "snap":
		cmd = exec.Command("snap", "remove", entry.Name)
	default:
		fmt.Printf(trans("installer_not_supported"), entry.Installer, entry.Name)
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf(trans("uninstall_failed"), entry.Name, err)
		return
	}


	fmt.Printf(trans("uninstalled_successfully"), entry.Name)
}		

func exit(status int) {
	fmt.Println(trans("exit_app"))

	os.Exit(status)
}
