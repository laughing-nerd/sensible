package secret

import (
	"os"
	"sensible/internal/constants"
	"sensible/internal/utils"
	"strings"
)

func Get(key string, all bool, env string) error {
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
	password, err := utils.AskPassword("Enter password to get the secret: ")
	if err != nil {
		return err
	}

	decrypted, err := utils.Decrypt(data, password)
	if err != nil {
		return err
	}

	// if user wants all secrets, print the whole decrypted content
	// no need to iterate through lines
	if all {
		os.Stdout.WriteString(string(decrypted) + "\n")
		return nil
	}

	lines := strings.SplitSeq(string(decrypted), "\n")
	for line := range lines {
		l := strings.TrimSpace(line)
		if l == "" || strings.HasPrefix(l, "#") {
			continue // skip empty lines and comments
		}

		if strings.HasPrefix(l, key+"=") {
			os.Stdout.WriteString(l + "\n")
			return nil
		}

	}

	return nil
}
