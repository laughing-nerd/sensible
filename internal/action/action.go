package action

import (
	"errors"
	"fmt"
	"maps"
	"sensible/internal/action/components"
	"sensible/internal/connectors"
	"sensible/internal/constants"
	"sensible/models"
	"sensible/pkg/hclparser"
	"sensible/pkg/logger"
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
)

func Do(filePath string, variables map[string]cty.Value, groups map[string]map[string]models.Host) error {
	var (
		wg    = &sync.WaitGroup{}
		mode  string
		hosts = make(map[string]models.Host)
	)

	actionBlocks, err := hclparser.GetBlocks(filePath)
	if err != nil {
		return err
	}

	// iterate over all actions
	for _, actionBlock := range actionBlocks {
		if actionBlock.Type != constants.Action || len(actionBlock.Labels) < 1 {
			continue
		}

		actionName := actionBlock.Labels[0]
		actionBody := actionBlock.Body

		// 0. fetch the host groups
		hgroups := make(map[string]cty.Value)
		if err := hclparser.GetBlockAttributes(actionBlock, hgroups); err != nil {
			return err
		}

		v, ok := hgroups["groups"]
		mode = constants.Remote

		// if groups is not present then the commands will be run locally
		if !ok {
			mode = constants.Local
		} else {

			// if groups is present then it should be an array of strings
			if !v.Type().IsTupleType() {
				return errors.New("`groups` must be an array of host groups")
			}

			// if everything is okay, then extract the hosts map
			for _, group := range v.AsValueSlice() {
				if h, ok := groups[group.AsString()]; ok {
					maps.Copy(hosts, h)
				}
			}

			logger.Info("Connecting to hosts")
			// Connect to the hosts here
			for hostname, hostVal := range hosts {
				var (
					authMethod string
					creds      string
				)

				if hostVal.Password != "" {
					authMethod = constants.Password
					creds = hostVal.Password
				} else if hostVal.PrivateKey != "" {
					authMethod = constants.PrivateKey
					creds = hostVal.PrivateKey
				} else {
					return errors.New("`password` or `private_key` is required for the host:" + hostname)
				}

				wg.Add(1)
				go func(hostname string, hostVal models.Host, hosts map[string]models.Host) {
					defer wg.Done()
					client, err := connectors.NewSshConnection(hostVal.Address, authMethod, hostVal.Username, creds, hostVal.Timeout)
					if err != nil {
						logger.Error("Unabled connect to host:", hostname)
					}
					host := hosts[hostname]
					host.SshClient = client
					hosts[hostname] = host
				}(hostname, hostVal, hosts)
			}
		}
		wg.Wait()

		// 1. Execute the components
		logger.Info("Running:", actionName)

		for _, componentBlock := range actionBody.Blocks {
			var component = components.ComponentMap[componentBlock.Type]

			evalCtx := &hcl.EvalContext{Variables: variables}
			if diags := gohcl.DecodeBody(componentBlock.Body, evalCtx, component); diags.HasErrors() {
				return fmt.Errorf("Error decoding component %s: %v", componentBlock.Type, diags)
			}

			logger.Custom("EXECUTING", constants.ColorYellow, componentBlock.Labels[0], "🚀")
			if err := components.Execute(component, mode, hosts); err != nil {
				logger.Error("Error executing component", componentBlock.Type, ":", err.Error())
				continue
			}
			logger.Plain("\n")
		}

	}

	return nil
}
