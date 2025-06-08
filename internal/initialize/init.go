package initialize

import (
	_ "embed"
	"errors"
	"os"
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
)

func Start(file string) error {
	if err := os.MkdirAll(file+"/resources", 0756); err != nil {
		return errors.New("Error creating `resources` directory\n" + err.Error())
	}

	if err := os.MkdirAll(file+"/actions", 0756); err != nil {
		return errors.New("Error creating `actions` directory\n" + err.Error())
	}

	if err := os.WriteFile(file+"/resources/hosts.hcl", []byte(hostFileTemplate), 0644); err != nil {
		return errors.New("Error creating `hosts.hcl` file\n" + err.Error())
	}

	if err := os.WriteFile(file+"/resources/values.hcl", []byte(variablesFileTemplate), 0644); err != nil {
		return errors.New("Error creating `values.hcl` file\n" + err.Error())
	}

	os.Stdout.WriteString("Sensible initialized successfully!!\n")
	return nil
}
