package components

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sensible/internal/connectors"
	"sensible/internal/constants"
	"sensible/models"
	"sensible/pkg/logger"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
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
	command, err := c.getPkgInstallCommand(constants.Local, nil)
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
	var wg sync.WaitGroup

	// iterate over all hosts and create a new ssh session for each
	for _, host := range hosts {
		wg.Add(1)
		go func(host models.Host) {
			defer wg.Done()
			session, err := connectors.NewSshSession(host.SshClient)
			if err != nil {
				logger.Error("failed to create SSH session for host ", host.Name, ": ", err.Error())
				return
			}
			defer session.Close()

			command, err := c.getPkgInstallCommand(constants.Remote, session)
			if err != nil {
				logger.Error("failed to get package install command for host ", host.Name, ": ", err.Error())
				return
			}

			pkg := strings.Join(c.Packages, " ")
			command = fmt.Sprintf(command, pkg)

			if err := session.Run(command); err != nil {
				logger.Error("failed to run command on host ", host.Name, ": ", err.Error())
			}
		}(host)
	}
	wg.Wait()
}

// helper func ...
func (c *InstallerComponent) getPkgInstallCommand(mode string, session *ssh.Session) (string, error) {
	var err error

	// If the preferred pkg manager is set and exists then return it
	if c.Preferred != "" {
		err = pkgLookup(mode, c.Preferred, session)
	}
	command, ok := pkgManagers[c.Preferred]
	if err == nil && ok {
		return command, nil
	}

	for name, command := range pkgManagers {
		if pkgLookup(mode, name, session) == nil {
			return command, nil
		}
	}

	return "", errors.New("no package manager found")
}

func pkgLookup(mode, name string, session *ssh.Session) error {
	var err error
	switch mode {
	case constants.Local:
		_, err = exec.LookPath(name)
	default:
		if session == nil {
			return errors.New("SSH session is required for remote execution")
		}
		err = session.Run(fmt.Sprintf("command -v %s", name))
	}
	return err
}
