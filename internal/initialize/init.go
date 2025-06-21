package initialize

import (
	_ "embed"
	"errors"
	"os"
	"sensible/internal/constants"
	"sensible/pkg/logger"
)

type sensibleFile struct {
	name    string
	content string
}

var (
	//go:embed templates/hosts.hcl
	hostFileTemplate string

	//go:embed templates/variables.hcl
	variablesFileTemplate string

	//go:embed templates/action.hcl
	actionFileTemplate string
)

func Start() error {

	// make required directories
	if err := os.MkdirAll(constants.ResourcesDir, 0756); err != nil {
		return errors.New("Error creating `resources` directory\n" + err.Error())
	}

	if err := os.MkdirAll(constants.ActionsDir, 0756); err != nil {
		return errors.New("Error creating `actions` directory\n" + err.Error())
	}

	// write the required files
	if err := os.WriteFile(constants.HostFile, []byte(hostFileTemplate), 0644); err != nil {
		return errors.New("Error creating `hosts.hcl` file\n" + err.Error())
	}

	if err := os.WriteFile(constants.VariablesFile, []byte(variablesFileTemplate), 0644); err != nil {
		return errors.New("Error creating `values.hcl` file\n" + err.Error())
	}

	if err := os.WriteFile(constants.ActionsDir+"/sample-action.hcl", []byte(actionFileTemplate), 0644); err != nil {
		return errors.New("Error creating `sample-action.hcl` file\n" + err.Error())
	}

	logger.Success("Sensible initialized successfully!!", "\n")
	return nil
}
