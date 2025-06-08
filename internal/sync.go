package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"sensible/internal/constants"
	"sensible/models"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty"
)

var (
	variables = make(map[string]cty.Value)
	groups    = make(map[string]map[string]models.Host)
)

func Sync() {
	hostsFilePath := filepath.Join(constants.MainDir, constants.HostFile)
	valuesFilePath := filepath.Join(constants.MainDir, constants.VariablesFile)

	// 0. Parse variables.hcl file
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(valuesFilePath)
	if diags.HasErrors() {
		os.Stdout.WriteString("Error syncing values: " + diags.Error())
		return
	}
	content, _, diags := file.Body.PartialContent(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "variables"},
		},
	})
	if diags.HasErrors() {
		os.Stdout.WriteString("Error syncing values: " + diags.Error())
		return
	}

	n := len(content.Blocks)
	if n < 1 || n > 1 {
		os.Stdout.WriteString("Error syncing values: There should be exactly 1 'variables' block in the file.")
		return
	}

	attrs, diags := content.Blocks[0].Body.JustAttributes()
	if diags.HasErrors() {
		os.Stdout.WriteString("Error syncing values: " + diags.Error())
		return
	}

	for name, attr := range attrs {
		val, diag := attr.Expr.Value(nil) // nil = no context needed here
		if diag.HasErrors() {
			os.Stdout.WriteString("Error evaluating variable " + name + ": " + diag.Error())
			return
		}
		variables[name] = val
	}

	evalCtx := &hcl.EvalContext{
		Variables: variables,
	}

	// 1. Decode the hosts file
	var hostConfig models.HostConfig
	if err := hclsimple.DecodeFile(hostsFilePath, evalCtx, &hostConfig); err != nil {
		os.Stdout.WriteString("Error syncing hosts: " + err.Error())
		return
	}

	// 2. Post process
	for _, group := range hostConfig.Groups {
		groupName := group.Name
		if _, ok := groups[groupName]; !ok {
			groups[groupName] = make(map[string]models.Host)
		}

		for _, host := range group.Hosts {
			hostName := host.Name
			if _, ok := groups[groupName][hostName]; !ok {
				groups[groupName][hostName] = host
			}
		}
	}

	fmt.Println(hostConfig.Groups)
	fmt.Println(hostConfig.Hosts)
}
