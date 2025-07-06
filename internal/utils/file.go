package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func GetFilePath(dir, filename, env string) (string, error) {
	if env == "" {
		env = "base"
	}

	dir = fmt.Sprintf(dir, env)
	return filepath.Abs(filepath.Join(dir, filename))
}
