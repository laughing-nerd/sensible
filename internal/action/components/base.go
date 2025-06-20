package components

import (
	"sensible/internal/connectors"
	"sensible/internal/constants"
	"sensible/models"
	"sensible/pkg/logger"
	"sync"
)

type Component interface {
	ValidateParams() error
	RunLocal() error
	RunRemote(map[string]models.Host)
}

type Base struct {
	Parallel bool `hcl:"parallel"` // will be implemented later
}

var ComponentMap = map[string]Component{
	"shell":     &ShellComponent{},
	"installer": &InstallerComponent{},
	"cron":      &CronComponent{},
}

func Execute(c Component, mode string, hosts map[string]models.Host) error {
	if err := c.ValidateParams(); err != nil {
		return err
	}

	switch mode {
	case constants.Remote:
		c.RunRemote(hosts)
		return nil
	default:
		return c.RunLocal()
	}
}

func (b *Base) RunSshCommand(hosts map[string]models.Host, command string) {
	var wg sync.WaitGroup

	// iterate over all hosts and create a new ssh session for each
	for _, host := range hosts {
		wg.Add(1)
		go func(host models.Host) {
			defer wg.Done()
			session, err := connectors.NewSshSession(host.SshClient)
			if err != nil {
				logger.Error("failed to create SSH session for host " + host.Name + ": " + err.Error())
				return
			}
			defer session.Close()

			if err := session.Run(command); err != nil {
				logger.Error("failed to run command on host " + host.Name + ": " + err.Error())
			}
		}(host)
	}
	wg.Wait()
}

func (b *Base) RunSshInstallCommand(hosts map[string]models.Host) {
	var wg sync.WaitGroup

	// iterate over all hosts and create a new ssh session for each
	for _, host := range hosts {
		wg.Add(1)
		go func(host models.Host) {
			defer wg.Done()
			session, err := connectors.NewSshSession(host.SshClient)
			if err != nil {
				logger.Error("failed to create SSH session for host " + host.Name + ": " + err.Error())
				return
			}
			defer session.Close()

			// TODO: Fix this to run the install command properly
			// if err := session.Run(command); err != nil {
			// 	logger.Error("failed to run command on host " + host.Name + ": " + err.Error())
			// }
		}(host)
	}
	wg.Wait()
}
