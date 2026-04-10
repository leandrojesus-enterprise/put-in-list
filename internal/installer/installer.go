// Package installer provides detection and execution of supported package managers.
package installer

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/leandrojesus-enterprise/put-in-list/internal/i18n"
	"github.com/leandrojesus-enterprise/put-in-list/internal/storage"
)

// SupportedInstallers maps each supported OS to the package managers it can use.
var SupportedInstallers map[string][]string = map[string][]string{
	"windows": {"winget", "choco"},
	"linux":   {"apt", "snap"},
}

// Service performs package install/uninstall operations via OS-level commands.
type Service struct{}

// NewService creates a new installer Service.
func NewService() *Service {
	return &Service{}
}

// DetectAvailable returns the subset of supported installers that are present in PATH
// for the given OS. Returns an empty slice if the OS is not recognised.
func (s *Service) DetectAvailable(osName string) []string {
	var available []string = []string{}
	var supported []string
	var ok bool
	supported, ok = SupportedInstallers[osName]
	if !ok {
		return available
	}
	// Check each supported tool via exec.LookPath to confirm it is in PATH.
	for _, tool := range supported {
		if _, err := exec.LookPath(tool); err == nil {
			available = append(available, tool)
		}
	}
	return available
}

// Uninstall removes the package described by entry using its recorded installer.
// Returns an error if the installer is unsupported or the uninstall command fails.
func (s *Service) Uninstall(entry storage.InstallEntry, tr *i18n.Translator) error {
	fmt.Printf(tr.Trans("uninstalling_package"), entry.Name, entry.Installer)

	// Build the uninstall command based on which installer originally installed the package.
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
		return fmt.Errorf(tr.Trans("instlaller_not_supported"), entry.Installer, entry.Name)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(tr.Trans("uninstall_failed"), entry.Name, err)
	}

	fmt.Printf(tr.Trans("uninstalled_successfully"), entry.Name)
	return nil
}
