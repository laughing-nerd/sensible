package initialize

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"sensible/internal/constants"
	"sensible/pkg/logger"
)

var (
	//go:embed templates/hosts.hcl
	hostFileTemplate string

	//go:embed templates/values.hcl
	valuesTemplateFile string

	//go:embed templates/action.hcl
	actionFileTemplate string

	//go:embed templates/secrets.hcl
	secretsFileTemplate string
)

type Starter struct {
	FileName string
	Template string
}

var starter = map[string][]Starter{
	constants.ResourcesDir: {
		{FileName: constants.HostFile, Template: hostFileTemplate},
		{FileName: constants.ValuesFile, Template: valuesTemplateFile},
	},
	constants.ActionsDir: {
		{FileName: constants.SampleActionFile, Template: actionFileTemplate},
	},
	constants.SecretsDir: {
		{FileName: constants.SecretsFile, Template: secretsFileTemplate},
	},
}

func Start(env string) error {
	for dir, files := range starter {

		// create the directory
		dirPath := fmt.Sprintf(dir, env)
		if err := os.MkdirAll(dirPath, 0756); err != nil {
			return fmt.Errorf("error creating `%s` directory: %w", dir, err)
		}

		// create the files in the directory
		for _, file := range files {
			if err := os.WriteFile(filepath.Join(dirPath, file.FileName), []byte(file.Template), 0644); err != nil {
				return fmt.Errorf("error creating `%s` file in `%s` directory: %w", file.FileName, dir, err)
			}
		}

	}

	logger.Success("Sensible initialized successfully!!", "\n")
	return nil
}
