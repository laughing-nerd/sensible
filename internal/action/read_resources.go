package action

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sensible/internal/constants"
	"sensible/internal/utils"
	"sensible/models"
	"sensible/pkg/hclparser"
	"sensible/pkg/logger"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty"
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

	if shouldReadValues {
		readValuesFile(constants.ResourcesDir, constants.ValuesFile, env, variables["values"])
	}

	if shouldReadSecrets {
		readSecretsFile(constants.ResourcesDir, constants.SecretsFile, env, variables["secrets"])
	}

	return variables, nil
}

func GetHosts(variables map[string]map[string]cty.Value, env string) (map[string]map[string]models.Host, error) {
	var groups = make(map[string]map[string]models.Host)

	hostsFilePath, err := utils.GetFilePath(constants.ResourcesDir, constants.HostFile, env)
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
func shouldRead(file string) (bool, bool) {
	contents, err := os.ReadFile(file)
	if err != nil {
		return false, false
	}

	return bytes.Contains(contents, []byte("${values.")), bytes.Contains(contents, []byte("${secrets."))
}

// Reads base and env-specific files into variables map for the given key ("values" or "secrets")
func readValuesFile(dir, file, env string, variables map[string]cty.Value) error {
	mainFilePath, err := utils.GetFilePath(dir, file, env)
	if err != nil {
		return err
	}

	baseFilePath, err := utils.GetFilePath(dir, file, "base")
	if err != nil {
		return err
	}

	// Read base file first if env != base
	if env != "base" {
		if !utils.FileExists(baseFilePath) {
			logger.Warn("Unable to read base values file, ignoring defaults")
		} else {
			getVariablesFromFile(baseFilePath, "values", variables)
		}
	}

	// intentional override of base file with env-specific file
	return getVariablesFromFile(mainFilePath, "values", variables)
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

func readSecretsFile(dir, file, env string, variables map[string]cty.Value) error {
	// read the secrets file
	secretFile, err := utils.GetFilePath(constants.SecretsDir, constants.SecretsFile, env)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(secretFile)
	if err != nil {
		return err
	}

	// get user password to decrypt the secret
	password, err := utils.AskPassword("Enter password to read secrets: ")
	if err != nil {
		return err
	}

	decrypted, err := utils.Decrypt(data, password)
	if err != nil {
		return err
	}

	lines := strings.SplitSeq(string(decrypted), "\n")
	for line := range lines {
		l := strings.TrimSpace(line)
		if l == "" || strings.HasPrefix(l, "#") {
			continue // skip empty lines and comments
		}

		kvPair := strings.SplitN(l, "=", 2)
		if len(kvPair) != 2 {
			return fmt.Errorf("Invalid secret format in line: %s", l)
		}
		key := strings.TrimSpace(kvPair[0])
		value := strings.TrimSpace(kvPair[1])

		variables[key] = cty.StringVal(value)
	}

	return nil

}
