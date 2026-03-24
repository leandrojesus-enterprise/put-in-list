package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
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

type StoredLists struct {
	ActiveList string              `json:"activeList"`
	Lists      map[string][]string `json:"lists"`
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

	storedLists = StoredLists{}
	storedLists.Lists = make(map[string][]string)

	//Check if lists.json is created
	jsonFile, err := os.Open("lists.json")
	if err != nil {
		fmt.Println("lists.json doesn't exist! Will create one.")
		pause()

		storedLists = StoredLists{
			ActiveList: "None",
			Lists:      map[string][]string{},
		}

		saveJson(storedLists)
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

func saveJson(data StoredLists) {
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
	fmt.Printf("6 %v»%v Choose Current Package Installer\n", Green, Reset)
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

		storedLists.Lists[listName] = []string{}

		saveJson(storedLists)
	}

	mainMenuChoices["2"] = func() {
		fmt.Println("Available lists:")
		for listName := range storedLists.Lists {
			fmt.Printf(" - %s\n", listName)
		}
		
		fmt.Print("Enter the name of the list you want to set as active: ")
		var listName string
		fmt.Scanln(&listName)

		if _, exists := storedLists.Lists[listName]; exists {
			storedLists.ActiveList = listName
			saveJson(storedLists)
		} else {
			fmt.Printf("List '%s' does not exist!\n", listName)
		}	
	}

	mainMenuChoices["3"] = func() {
		fmt.Println("option 3")
	}

	mainMenuChoices["4"] = func() {
		fmt.Println("option 4")
	}

	mainMenuChoices["5"] = func() {
		fmt.Println("option 5")
	}

	mainMenuChoices["6"] = func() {
		fmt.Println("option 6")
	}

	mainMenuChoices["7"] = func() {
		fmt.Println("option 7")
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
