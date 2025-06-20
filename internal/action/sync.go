package action

import (
	"errors"
	"path/filepath"
	"sensible/internal/constants"
	"sensible/models"
	"sensible/pkg/hclparser"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty"
)

func Sync() (map[string]cty.Value, map[string]map[string]models.Host, error) {
	var (
		variables = make(map[string]cty.Value)
		groups    = make(map[string]map[string]models.Host)
	)

	hostsFilePath := filepath.Join(constants.MainDir, constants.HostFile)
	valuesFilePath := filepath.Join(constants.MainDir, constants.VariablesFile)

	// 0. Get variables map first
	blocks, err := hclparser.GetBlocks(valuesFilePath)
	if err != nil {
		return nil, nil, err
	}

	if len(blocks) > 1 {
		return nil, nil, errors.New("There should be atmost 1 `variables` block")
	}

	if blocks[0].Type != constants.Variables {
		return nil, nil, errors.New("There should be a `variables` block")
	}

	if err := hclparser.GetBlockAttributes(blocks[0], variables); err != nil {
		return nil, nil, err
	}

	// 1. Get hosts map
	evalCtx := &hcl.EvalContext{Variables: variables}
	var hostConfig models.HostConfig
	if err := hclsimple.DecodeFile(hostsFilePath, evalCtx, &hostConfig); err != nil {
		return nil, nil, err
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

	return variables, groups, nil
}
