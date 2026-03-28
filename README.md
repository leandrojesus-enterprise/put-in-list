# 📦 Put In List (CLI)

> Automate package installation & uninstallation using organized lists.

**Put In List** is a cross-platform CLI tool written in Go that lets you manage software packages through lists — making it easy to install or remove multiple packages in bulk.

---

## ✨ Features

- 📋 Create and manage multiple package lists  
- 🎯 Set an active list for operations  
- 📦 Install packages and automatically save them to a list  
- 🗑️ Uninstall all packages from a list (mass uninstall)  
- 🔎 Search for packages across all lists  
- ⚙️ Choose your preferred package manager  
- 💾 Persistent storage using JSON (`config.json`, `lists.json`)  

---

## 🖥️ Supported Systems

### Windows
- `winget`
- `choco` (soon)

### Linux
- `apt` (soon)
- `snap`

> The app automatically detects available installers in your system `PATH`.

---

## 🚀 Installation

### 1. Clone the repository
```bash
git clone https://github.com/yourusername/put-in-list.git
cd put-in-list
```

### 2. Build the app
```bash
go build -o putinlist
```

### 3. Run
```bash
./putinlist
```

---

## 📂 File Structure

- `config.json` → Stores selected package manager  
- `lists.json` → Stores all lists and packages  

These files are automatically created if they don't exist.

---

## 🧠 How It Works

### Main Menu

```
1 » Create List
2 » Set Active List
3 » Show Lists
4 » Uninstall List
-------------------------
5 » Install Package
6 » Choose Current Package Installer
7 » Search For Package In a List
-------------------------
8 » Exit
```

---

## 🧪 Usage Examples

### ➕ Create a list
```
Enter the desired new list name: dev-tools
```

### 🎯 Set active list
Choose from existing lists by index.

### 📦 Install a package
```
Name of the package that you want to install: vscode
```

For snap classic installs:
```bash 
code --classic
```

### 🗑️ Uninstall all packages from a list
- Removes packages one by one using their respective installer  
- Optionally deletes the list after uninstall  

### 🔎 Search for a package
```
Enter the name of the package: vscode
```

---

## ⚠️ Notes

- You **must select an installer** before installing packages  
- You **must set an active list** before installing/uninstalling  
- Some commands may require **sudo/admin privileges**  
- Only supported installers are allowed per OS  

---

## 💡 Use Cases

 - 🧑‍💻 Quick development environment setup (soon)
 
   - Install all your tools (VSCode, Git, Node, etc.) in one go.

 - 🔄 Restore your setup after OS reinstall (soon)
 
   - No need to remember every app — everything is saved in a list.

- 🧹 Bulk software cleanup

   - Uninstall dozens of applications at once.

- 📦 Reproducible environments (soon)

   - Keep consistent setups across different machines.

- 🧪 Testing different environments (soon)
                           
  - Create separate lists for different stacks (e.g., backend, frontend, DevOps).

- 🧑‍🏫 Teaching / Workshops (soon)

  - Share a list with all required tools for students or participants.

- 🧰 Personal software management

  - Keep an organized record of what you use and with which installer. 

---

## 🔧 Tech Stack

- Go (Golang)  
- Standard Library only (no external dependencies) til today.

---

## 📜 License

MIT License (or your preferred license)

---

## 👨‍💻 Author

**Leandro Jesus**  
🇵🇹 Portugal
