package components

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sensible/models"
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
	"pacman":  "sudo pacman -S --noconfirm %s",
	"yum":     "sudo yum install -y %s",
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

func (c *InstallerComponent) RunLocal() error {
	command, err := c.getPkgInstallCommand()
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

func (c *InstallerComponent) RunRemote(hosts map[string]models.Host) {
	c.RunSshInstallCommand(hosts)
}

// helper func ...
func (c *InstallerComponent) getPkgInstallCommand() (string, error) {

	// If the preferred pkg manager is set and exists, return it
	if c.Preferred != "" {
		_, err := exec.LookPath(c.Preferred)
		command, ok := pkgManagers[c.Preferred]
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

func (c *Base) GetName() string {
	return "Hello"
}
