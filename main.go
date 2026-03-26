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
	fmt.Println("** App initializing **")

	currentOS = runtime.GOOS
	if !currentOSSupported() {
		fmt.Println("Current OS not supported.")
		pause()
		exit(1)
	}

	config = Config{CurrentInstaller: "None"}

	cfgFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("config.json não existe, criando com None.")
		pause()
		saveConfig(config)
	} else {
		byteValue, err := io.ReadAll(cfgFile)
		if err != nil {
			fmt.Println("Erro ler config.json:", err)
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
		fmt.Println("lists.json doesn't exist! Will create one.")
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
			fmt.Println(err)
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
		fmt.Println("Erro marshal config:", err)
		pause()
		exit(1)
	}
	err = os.WriteFile("config.json", byteValue, 0644)
	if err != nil {
		fmt.Println("Erro escrever config.json:", err)
		pause()
		exit(1)
	}
}

func clearTerminal() {
	currentOS := runtime.GOOS
	clearTerminalFunc[currentOS]()
}

func showMainMenu() {
	fmt.Printf("           %v{%v Welcome to Put In List %v}%v\n\n", Green, Reset, Green, Reset)
	fmt.Printf("1 %v»%v Create List\n", Green, Reset)
	fmt.Printf("2 %v»%v Set List Active %v[%v Current: %s %v]%v\n", Green, Reset, Green, Reset, storedLists.ActiveList, Green, Reset)
	fmt.Printf("3 %v»%v Show lists\n", Green, Reset)
	fmt.Printf("4 %v»%v Uninstall List\n", Green, Reset)
	fmt.Printf("%v-------------------------%v\n", Green, Reset)
	fmt.Printf("5 %v»%v Install Package/s\n", Green, Reset)
	fmt.Printf("6 %v»%v Choose Current Package Installer %v[%v Current: %s %v]%v\n", Green, Reset, Green, Reset, config.CurrentInstaller, Green, Reset)
	fmt.Printf("7 %v»%v Search For Package In a List\n", Green, Reset)
	fmt.Printf("%v-------------------------%v\n", Green, Reset)
	fmt.Printf("8 %v»%v Exit\n\n", Green, Reset)
}

func readMainMenuChoice() {
	fmt.Print("Enter your choice: ")
	fmt.Scanln(&choice)
}

func runMainMenuChoice() {
	if choice >= "1" && choice <= "8" {
		mainMenuChoices[choice]()
	} else {
		fmt.Println("Invalid choice! Enter a choice between <1-8>.")
	}
}

func pause() {
	fmt.Println("Press <Enter> to continue.")
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
		fmt.Print("Enter the desired new list name: ")
		var listName string
		fmt.Scanln(&listName)

		if strings.TrimSpace(listName) == "" {
			fmt.Println("List name cannot be empty.")
			return
		}
		
		if _, exists := storedLists.Lists[listName]; exists {
			fmt.Printf("List '%s' already exists!\n", listName)
			return
		}

		storedLists.Lists[listName] = []InstallEntry{}

		saveListsJson(storedLists)
	}

	mainMenuChoices["2"] = func() {
		fmt.Println("Available lists ordered by index:")
		i := 1
		for listName := range storedLists.Lists {
			fmt.Printf(" %d - %s\n", i, listName)
			i++
		}

		fmt.Print("Enter the index of the list you want to set as active: ")
		var listIndex int
		fmt.Scanln(&listIndex)

		if listIndex < 1 || listIndex > len(storedLists.Lists) {
			fmt.Println("Invalid list index!")
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
			fmt.Printf("List '%s' does not exist!\n", listName)
		}
	}

	mainMenuChoices["3"] = func() {
		fmt.Println("Lists:")
		for listName, entries := range storedLists.Lists {
			activeMarker := ""
			if listName == storedLists.ActiveList {
				activeMarker = " (active)"
			}
			fmt.Printf(" - %s%s\n", listName, activeMarker)
			for _, entry := range entries {
				fmt.Printf("   - %s (%s)\n", entry.Name, entry.Installer)
			}
		}
	}

	mainMenuChoices["4"] = func() {
		if storedLists.ActiveList == "None" {
			fmt.Println("Erro: nenhuma lista ativa. Use opção 2 para definir uma lista.")
			return
		}

		entries, ok := storedLists.Lists[storedLists.ActiveList]
		if !ok {
			fmt.Printf("Erro: lista ativa '%s' não existe.\n", storedLists.ActiveList)
			return
		}

		if len(entries) == 0 {
			fmt.Printf("Lista '%s' está vazia.\n", storedLists.ActiveList)
			fmt.Println("Do you want to delete the list? (y/n)")
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) == "y" {
				delete(storedLists.Lists, storedLists.ActiveList)
				storedLists.ActiveList = "None"
				saveListsJson(storedLists)
			}
			return
		} else {
			fmt.Print("Do you really want to uninstall this list? this will unistall all packages in the list.")
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
				return
			}

			fmt.Print("Do you want to remove the list? (y/n)")
			fmt.Scanln(&response)
			if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
				delete(storedLists.Lists, storedLists.ActiveList)
				storedLists.ActiveList = "None"
				saveListsJson(storedLists)
			}
			
		}

		for i := len(entries) - 1; i >= 0; i-- {
			entry := entries[i]
			fmt.Printf("Desinstalando '%s' com '%s'...\n", entry.Name, entry.Installer)

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
				fmt.Printf("Falha ao desinstalar '%s': %v\n", entry.Name, err)
				continue
			}

			fmt.Printf("'%s' desinstalado com sucesso.\n", entry.Name)

			entries = append(entries[:i], entries[i+1:]...)
		}

		storedLists.Lists[storedLists.ActiveList] = entries
		saveListsJson(storedLists)

		if len(entries) == 0 {
			fmt.Printf("Todos os pacotes removidos da lista '%s'.\n", storedLists.ActiveList)
		}
	}

	mainMenuChoices["5"] = func() {
		if config.CurrentInstaller == "None" {
			fmt.Println("Erro: nenhum instalador selecionado. Use opção 6.")
			return
		}

		if storedLists.ActiveList == "None" {
			fmt.Println("Erro: nenhuma lista ativa. Use opção 2.")
			return
		}

		fmt.Print("Nome do pacote para instalar: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan() // use `for scanner.Scan()` to keep reading
		pkg := scanner.Text()

		if strings.TrimSpace(pkg) == "" {
			fmt.Println("Nome do pacote não pode ser vazio.")
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

		fmt.Printf("Instalando '%s' com %s...\n", pkgName, config.CurrentInstaller)

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
			fmt.Printf("Falha ao instalar '%s': %v\n", pkgName, err)
			fmt.Println("Não será adicionado à lista ativa.")
			return
		}

		fmt.Printf("'%s' instalado com sucesso.\n", pkgName)

		activeEntries := storedLists.Lists[storedLists.ActiveList]
		for _, existing := range activeEntries {
			if existing.Name == pkg && existing.Installer == config.CurrentInstaller {
				fmt.Printf("'%s' já existe na lista ativa '%s' com instalador '%s'.\n", pkgName, storedLists.ActiveList, config.CurrentInstaller)
				return
			}
		}

		storedLists.Lists[storedLists.ActiveList] = append(activeEntries, InstallEntry{Name: pkgName, Installer: config.CurrentInstaller})
		saveListsJson(storedLists)
		fmt.Printf("'%s' [%s] adicionado à lista ativa '%s'.\n", pkgName, config.CurrentInstaller, storedLists.ActiveList)
	}

	mainMenuChoices["6"] = func() {
		supported := supportedInstallers[currentOS]
		available := detectAvailableInstallers(currentOS)

		fmt.Printf("Instaladores suportados para %s: %s\n", currentOS, strings.Join(supported, ", "))
		if len(available) == 0 {
			fmt.Println("Nenhum instalador detectado disponível no PATH.")
			fmt.Println("Para usar a opção, instale um deles ou adicione ao PATH.")
			return
		}

		fmt.Println("Instaladores detectados disponíveis:")
		for i, tool := range available {
			currentMark := ""
			if tool == config.CurrentInstaller {
				currentMark = " (atual)"
			}
			fmt.Printf("%d) %s%s\n", i+1, tool, currentMark)
		}

		fmt.Print("Escolha o número do instalador para definir como atual: ")
		var opt int
		_, err := fmt.Scanln(&opt)
		if err != nil || opt < 1 || opt > len(available) {
			fmt.Println("Escolha inválida.")
			return
		}

		selected := available[opt-1]
		config.CurrentInstaller = selected
		saveConfig(config)

		fmt.Printf("Instalador atual definido: %s\n", selected)
	}

	mainMenuChoices["7"] = func() {
		fmt.Print("Enter the name of the package you want to search for in all lists: ")
		var packageName string
		fmt.Scanln(&packageName)
		packageName = strings.ToLower(packageName)

		found := false
		for listName, entries := range storedLists.Lists {
			for _, entry := range entries {
				if strings.ToLower(entry.Name) == packageName {
					fmt.Printf("Package '%s' found in list '%s' (installer=%s)\n", packageName, listName, entry.Installer)
					found = true
				}
			}
		}

		if !found {
			fmt.Printf("Package '%s' not found in any list.\n", packageName)
		}
	}

	mainMenuChoices["8"] = func() {
		exit(0)
	}
}

func exit(status int) {
	fmt.Println("Exiting the app!")
	fmt.Println("See you soon.")

	os.Exit(status)
}
