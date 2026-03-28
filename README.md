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
- `choco`

### Linux
- `apt`
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
Nome do pacote para instalar: vscode
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

- 🧑‍💻 Quickly set up a new development environment  
- 🔄 Reinstall your favorite tools after OS reinstall  
- 🧹 Clean up systems with bulk uninstall  
- 📦 Maintain reproducible software setups  

---

## 🔧 Tech Stack

- Go (Golang)  
- Standard Library only (no external dependencies)

---

## 📜 License

MIT License (or your preferred license)

---

## 👨‍💻 Author

**Leandro Jesus**  
🇵🇹 Portugal
