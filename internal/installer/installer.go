package installer

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/leandrojesus-enterprise/put-in-list/internal/i18n"
	"github.com/leandrojesus-enterprise/put-in-list/internal/storage"
)

var SupportedInstallers = map[string][]string{
	"windows": {"winget", "choco"},
	"linux":   {"apt", "snap"},
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) DetectAvailable(osName string) []string {
	available := []string{}
	supported, ok := SupportedInstallers[osName]
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

func (s *Service) Uninstall(entry storage.InstallEntry, tr *i18n.Translator) error {
	fmt.Printf(tr.Trans("uninstalling_package"), entry.Name, entry.Installer)

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
