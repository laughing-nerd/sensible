package action

import (
	"bytes"
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
	"golang.org/x/sync/errgroup"
)

// TODO: Add env variable support
func GetVariables(actionFile, env string) (map[string]map[string]cty.Value, error) {
	// first check if the action file has any variables or secrets
	shouldReadValues, shouldReadSecrets := shouldRead(actionFile)
	if !shouldReadValues && !shouldReadSecrets {
		return nil, nil
	}

	variables := map[string]map[string]cty.Value{
		"values":  {},
		"secrets": {},
	}

	// errCh := make(chan error, 2)
	var g errgroup.Group

	if shouldReadValues {
		g.Go(func() error {
			return readVariableSet(constants.ResourcesDir, constants.ValuesFile, constants.Values, env, variables)
		})
	}

	if shouldReadSecrets {
		g.Go(func() error {
			return readVariableSet(constants.SecretsDir, constants.SecretsFile, constants.Secrets, env, variables)
		})
	}
	// Wait for err results
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return variables, nil
}

func GetHosts(variables map[string]map[string]cty.Value, env string) (map[string]map[string]models.Host, error) {
	var groups = make(map[string]map[string]models.Host)

	hostsFilePath, err := getFile(constants.ResourcesDir, constants.HostFile, env)
	if err != nil {
		return nil, err
	}

	// 1. Get hosts map
	evalCtx := &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"value":  cty.ObjectVal(variables["values"]),
			"secret": cty.ObjectVal(variables["secrets"]),
		},
	}
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

func getVariablesFromFile(file, blockType string, variables map[string]cty.Value) error {
	blocks, err := hclparser.GetBlocks(file)
	if err != nil {
		return err
	}

	if len(blocks) > 1 || blocks == nil {
		return errors.New("There should be atmost 1 `variables` block")
	}

	if blocks[0].Type != blockType {
		return fmt.Errorf("Expected a `%s` block, but found `%s`", blockType, blocks[0].Type)
	}

	return hclparser.GetBlockAttributes(blocks[0], variables)
}

func getFile(dir, file, env string) (string, error) {
	dir = fmt.Sprintf(dir, env)
	joinedFile := filepath.Join(dir, file)
	return filepath.Abs(joinedFile)
}

func shouldRead(file string) (bool, bool) {
	contents, err := os.ReadFile(file)
	if err != nil {
		return false, false
	}

	return bytes.Contains(contents, []byte("${value.")), bytes.Contains(contents, []byte("${secret."))
}

// Reads base and env-specific files into variables map for the given key ("values" or "secrets")
func readVariableSet(dir, file, key, env string, variables map[string]map[string]cty.Value) error {
	mainFilePath, err := getFile(dir, file, env)
	if err != nil {
		return err
	}

	baseFilePath, err := getFile(dir, file, "base")
	if err != nil {
		return err
	}

	// Read base file first if env != base
	if env != "base" {
		if !fileExists(baseFilePath) {
			logger.Warn(fmt.Sprintf("Unable to read base %s file, ignoring defaults", key))
		} else {
			_ = getVariablesFromFile(baseFilePath, key, variables[key])
		}
	}

	// intentional override of base file with env-specific file
	return getVariablesFromFile(mainFilePath, key, variables[key])
}
