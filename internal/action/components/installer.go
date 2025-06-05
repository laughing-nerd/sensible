package components

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type InstallerComponent struct {
	Base
	Name      string   `hcl:"installer,label"`
	Preferred string   `hcl:"preferred,optional"` // preferred package manager
	Packages  []string `hcl:"packages"`           // packages to install
}

var pkgManagers = map[string]string{
	"apt":     "sudo apt install -y %s",
	"apt-get": "sudo apt-get install -y %s",
	"dnf":     "sudo dnf install -y %s",
	"yum":     "sudo yum install -y %s",
	"pacman":  "sudo pacman -S --noconfirm %s",
	"apk":     "sudo apk add %s",
	"brew":    "brew install %s",
	"port":    "sudo port install %s",
	"nix-env": "nix-env -iA nixpkgs.%s",
	"pkg":     "sudo pkg install -y %s",
	"zypper":  "sudo zypper install -y %s",
}

func (c *InstallerComponent) ValidateParams() error {
	return nil
}

func (c *InstallerComponent) Run() error {
	command, err := getPkgInstallCommand(c.Preferred)
	if err != nil {
		return err
	}

	pkg := strings.Join(c.Packages, " ")
	cmd := fmt.Sprintf(command, pkg)

	execCommand := exec.Command("sh", "-c", cmd)

	execCommand.Stdout = os.Stdout
	execCommand.Stderr = os.Stderr

	return execCommand.Run()
}

// helper func ...
func getPkgInstallCommand(preferred string) (string, error) {

	// If the preferred pkg manager is set and exists, return it
	if preferred != "" {
		_, err := exec.LookPath(preferred)
		command, ok := pkgManagers[preferred]
		if err == nil && ok {
			return command, nil
		}
	}

	for name, command := range pkgManagers {
		if _, err := exec.LookPath(name); err == nil {
			return command, nil
		}
	}

	return "", errors.New("no package manager found")
}
