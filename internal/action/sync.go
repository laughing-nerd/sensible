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

func Sync(env string) (map[string]cty.Value, map[string]map[string]models.Host, error) {
	var (
		variables        = make(map[string]cty.Value)
		groups           = make(map[string]map[string]models.Host)
		resourcesDir     = fmt.Sprintf(constants.ResourcesDir, env)
		baseResourcesDir = fmt.Sprintf(constants.ResourcesDir, "base")
		hostsFile        = filepath.Join(resourcesDir, constants.HostFile)
		valuesFile       = filepath.Join(resourcesDir, constants.VariablesFile)
		baseValuesFile   = filepath.Join(baseResourcesDir, constants.VariablesFile)
	)

	hostsFilePath, err := filepath.Abs(hostsFile)
	if err != nil {
		return nil, nil, err
	}

	valuesFilePath, err := filepath.Abs(valuesFile)
	if err != nil {
		return nil, nil, err
	}

	baseValuesFilePath, err := filepath.Abs(baseValuesFile)
	if err != nil {
		return nil, nil, err
	}

	// 0. Get variables map first
	// try to get the base values file first if env is not base (fallback mechanism)
	// if this operation fails, it will ignore the base values file and log a warning
	if env != "base" {
		if !fileExists(baseValuesFilePath) {
			logger.Warn("Unable to read base values file, ignoring default values")
		} else {
			getVariablesFromFile(baseValuesFilePath, variables) //nolint:errcheck
		}
	}

	// this will read the intended values file for the env
	// and will overwrite variables if base values file exists for env != "base"
	if err := getVariablesFromFile(valuesFilePath, variables); err != nil {
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
