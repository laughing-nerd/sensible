package action

import (
	"fmt"
	"path/filepath"
	"sensible/constants"
	"sensible/internal/action/components"
	"sensible/pkg/logger"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func Parse(path string) {
	parser := hclparse.NewParser()
	f, _ := filepath.Abs(path)
	file, diags := parser.ParseHCLFile(f)
	if diags.HasErrors() {
		logger.Error(fmt.Sprintf("Error parsing file %s: %v", f, diags))
		return
	}

	rootBody := file.Body.(*hclsyntax.Body)

	// iterate over all actions in the rootBody
	for _, block := range rootBody.Blocks {
		if block.Type != constants.Action || len(block.Labels) < 1 {
			continue
		}

		actionName := block.Labels[0]
		actionBody := block.Body

		// Extract hosts or groups from the action body
		// hostsAttr, hostsOk := actionBody.Attributes["hosts"]
		// groupsAttr, groupsOk := actionBody.Attributes["groups"]
		//
		// if (!hostsOk && !groupsOk) || (hostsOk && groupsOk) {
		// 	logger.Error(fmt.Sprintf("Action %s must have either 'hosts' or 'groups' attribute", actionName))
		// 	return
		// }
		//
		// if ok {
		// 	value, diags := hostsAttr.Expr.Value(nil)
		// 	if diags.HasErrors() {
		// 		logger.Error(fmt.Sprintf("Error evaluating 'hosts' attribute in action %s: %v", actionName, diags))
		// 		return
		// 	}
		//
		// 	if value.Type().IsTupleType() || value.Type().IsListType() {
		// 		hosts := []string{}
		// 		for _, v := range value.AsValueSlice() {
		// 			hosts = append(hosts, v.AsString())
		// 		}
		// 		logger.Info(fmt.Sprintf("Hosts for action %s: %v", actionName, hosts))
		// 	} else {
		// 		logger.Error(fmt.Sprintf("'hosts' in action %s is not a list", actionName))
		// 	}
		// }
		//
		// return

		logger.Info("Running action: " + actionName)

		for _, componentBlock := range actionBody.Blocks {

			var component = components.ComponentMap[componentBlock.Type]
			if diags := gohcl.DecodeBody(componentBlock.Body, nil, component); diags.HasErrors() {
				logger.Error(fmt.Sprintf("Error decoding component %s: %v", componentBlock.Type, diags))
				return
			}

			logger.Custom("EXECUTING... "+componentBlock.Labels[0], constants.ColorYellow, "🚀")
			if err := components.Execute(component); err != nil {
				logger.Error(fmt.Sprintf("Error executing component %s: %v", componentBlock.Type, err))
				continue
			}
			logger.Plain("\n")

		}

	}

}
