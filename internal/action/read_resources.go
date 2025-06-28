package action

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sensible/internal/constants"
	"sensible/models"
	"sensible/pkg/hclparser"
	"sensible/pkg/logger"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty"
)

// TODO: Add env variable support
func GetValues(env string) (map[string]cty.Value, error) {
	var variables = make(map[string]cty.Value)

	valuesFilePath, err := getFile(constants.ResourcesDir, constants.VariablesFile, env)
	if err != nil {
		return nil, err
	}

	baseValuesFilePath, err := getFile(constants.ResourcesDir, constants.VariablesFile, "base")
	if err != nil {
		return nil, err
	}

	if env != "base" {
		if !fileExists(baseValuesFilePath) {
			logger.Warn("Unable to read base values file, ignoring default values")
		} else {
			_ = getVariablesFromFile(baseValuesFilePath, variables)
		}
	}

	// this will read the intended values file for the env
	// and will overwrite variables if base/resources/values.hcl file exists for env != "base"
	if err := getVariablesFromFile(valuesFilePath, variables); err != nil {
		return nil, err
	}

	return variables, nil
}

func GetHosts(variables map[string]cty.Value, env string) (map[string]map[string]models.Host, error) {
	var groups = make(map[string]map[string]models.Host)

	hostsFilePath, err := getFile(constants.ResourcesDir, constants.HostFile, env)
	if err != nil {
		return nil, err
	}

	// 1. Get hosts map
	evalCtx := &hcl.EvalContext{Variables: variables}
	var hostConfig models.HostConfig
	if err := hclsimple.DecodeFile(hostsFilePath, evalCtx, &hostConfig); err != nil {
		return nil, err
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

	return groups, nil
}

// helper func ...
func fileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil || !os.IsNotExist(err)
}

func getVariablesFromFile(file string, variables map[string]cty.Value) error {
	blocks, err := hclparser.GetBlocks(file)
	if err != nil {
		return err
	}

	if len(blocks) > 1 || blocks == nil {
		return errors.New("There should be atmost 1 `variables` block")
	}

	if blocks[0].Type != constants.Variables {
		return errors.New("There should be a `variables` block")
	}

	return hclparser.GetBlockAttributes(blocks[0], variables)
}

func getFile(dir, file, env string) (string, error) {
	dir = fmt.Sprintf(dir, env)
	joinedFile := filepath.Join(dir, file)
	return filepath.Abs(joinedFile)

}
