package components

import (
	"os"
	"os/exec"
	"sensible/models"
)

type ShellComponent struct {
	Base
	Name    string `hcl:"shell,label"`
	Command string `hcl:"command"`
}

func (c *ShellComponent) ValidateParams() error {
	return nil
}

func (c *ShellComponent) RunLocal() error {
	cmd := exec.Command("sh", "-c", c.Command)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (c *ShellComponent) RunRemote(hosts map[string]models.Host) {
	c.RunSshCommand(hosts, c.Command)
}
