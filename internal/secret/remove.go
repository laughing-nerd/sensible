package secret

import (
	"os"
	"sensible/internal/constants"
	"sensible/internal/utils"
	"sensible/pkg/logger"
	"strings"
)

func Remove(key, env string) error {
	// read the secrets file
	secretFile, err := utils.GetFilePath(constants.SecretsDir, constants.SecretsFile, env)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(secretFile)
	if err != nil {
		return err
	}

	// get user password to decrypt the secret file
	password, err := utils.AskPassword("Enter password to delete the secret: ")
	if err != nil {
		return err
	}

	decrypted, err := utils.Decrypt(data, password)
	if err != nil {
		return err
	}

	var newSecrets []string
	lines := strings.SplitSeq(string(decrypted), "\n")
	for line := range lines {
		l := strings.TrimSpace(line)
		if l == "" || strings.HasPrefix(l, "#") {
			continue // skip empty lines or comments
		}
		if !strings.HasPrefix(l, key+"=") {
			newSecrets = append(newSecrets, l)
		}
	}

	secretToEnc := strings.Join(newSecrets, "\n")

	// encrypt the data and append it to the file
	encrypted, err := utils.Encrypt([]byte(secretToEnc), password)
	if err != nil {
		return err
	}

	err = os.WriteFile(secretFile, encrypted, 0600)
	if err != nil {
		return err
	}

	logger.Success("Secret deleted successfully!")
	return nil
}
