package initialize

import (
	_ "embed"
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

func Start() {
	if err := os.MkdirAll("./.sensible/resources", 0756); err != nil {
		os.Stdout.WriteString("Error creating directories: " + err.Error())
		return
	}

	if err := os.MkdirAll("./.sensible/actions", 0756); err != nil {
		os.Stdout.WriteString("Error creating directories: " + err.Error())
		return
	}

	if err := os.WriteFile("./.sensible/resources/hosts.hcl", []byte(hostFileTemplate), 0644); err != nil {
		os.Stdout.WriteString("Error creating host file: " + err.Error())
		return
	}

	if err := os.WriteFile("./.sensible/resources/values.hcl", []byte(variablesFileTemplate), 0644); err != nil {
		os.Stdout.WriteString("Error creating values file: " + err.Error())
		return
	}

	os.Stdout.WriteString("Sensible initialized successfully!!\n")
}
