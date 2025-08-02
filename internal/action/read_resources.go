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
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/zclconf/go-cty/cty"
)

// TODO: Add env variable support
func GetVariables(actionFile, env string) (map[string]map[string]cty.Value, error) {
	// first check if the action file has any variables or secrets
	shouldReadValues, shouldReadSecrets, shouldReadFacts := shouldRead(actionFile)
	if !shouldReadValues && !shouldReadSecrets && !shouldReadFacts {
		return nil, nil
	}

	variables := map[string]map[string]cty.Value{
		"values":  {},
		"secrets": {},
		"facts":   {},
	}

	if shouldReadValues {
		if err := readValuesFile(constants.ResourcesDir, constants.ValuesFile, env, variables["values"]); err != nil {
			return nil, err
		}
	}

	if shouldReadSecrets {
		if err := readSecretsFile(env, variables["secrets"]); err != nil {
			return nil, err
		}
	}

	if shouldReadFacts {
		readFacts(variables["facts"])
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
func shouldRead(file string) (bool, bool, bool) {
	contents, err := os.ReadFile(file)
	if err != nil {
		return false, false, false
	}

	return bytes.Contains(contents, []byte("${values.")),
		bytes.Contains(contents, []byte("${secrets.")),
		bytes.Contains(contents, []byte("${facts."))
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
			if err := getVariablesFromFile(baseFilePath, "values", variables); err != nil {
				return err
			}
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
		return errors.New("there should be atmost 1 `variables` block")
	}

	if blocks[0].Type != blockType {
		return fmt.Errorf("expected a `%s` block, but found `%s`", blockType, blocks[0].Type)
	}

	return hclparser.GetBlockAttributes(blocks[0], variables)
}

func readSecretsFile(env string, variables map[string]cty.Value) error {
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
			return fmt.Errorf("invalid secret format in line: %s", l)
		}
		key := strings.TrimSpace(kvPair[0])
		value := strings.TrimSpace(kvPair[1])

		variables[key] = cty.StringVal(value)
	}

	return nil
}

func readFacts(facts map[string]cty.Value) {
	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)

	wg.Add(4)
	go func() {
		defer wg.Done()
		hostInfo, err := host.Info()
		if err == nil {
			mu.Lock()
			facts["os"] = cty.StringVal(hostInfo.OS)
			facts["os_version"] = cty.StringVal(hostInfo.PlatformVersion)
			facts["kernel_version"] = cty.StringVal(hostInfo.KernelVersion)
			facts["architecture"] = cty.StringVal(hostInfo.KernelArch)
			facts["hostname"] = cty.StringVal(hostInfo.Hostname)
			facts["uptime"] = cty.NumberIntVal(int64(hostInfo.Uptime))
			facts["boot_time"] = cty.NumberIntVal(int64(hostInfo.BootTime))
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		cpuInfo, err := cpu.Info()
		if err == nil && len(cpuInfo) > 0 {
			mu.Lock()
			facts["cpu_model"] = cty.StringVal(cpuInfo[0].ModelName)
			facts["cpu_cores"] = cty.NumberIntVal(int64(cpuInfo[0].Cores))
			mu.Unlock()
		}

		cores, err := cpu.Counts(true)
		if err == nil {
			mu.Lock()
			facts["cpu_logical_cores"] = cty.NumberIntVal(int64(cores))
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		memStats, err := mem.VirtualMemory()
		if err == nil {
			mu.Lock()
			facts["memory_total_mb"] = cty.NumberIntVal(int64(memStats.Total / 1024 / 1024))
			facts["memory_available_mb"] = cty.NumberIntVal(int64(memStats.Available / 1024 / 1024))
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		diskStats, err := disk.Usage("/")
		if err == nil {
			mu.Lock()
			facts["disk_total"] = cty.NumberIntVal(int64(diskStats.Total / 1024 / 1024))
			facts["disk_used"] = cty.NumberIntVal(int64(diskStats.Used / 1024 / 1024))
			facts["disk_free"] = cty.NumberIntVal(int64(diskStats.Free / 1024 / 1024))
			mu.Unlock()
		}
	}()
	wg.Wait()
}
