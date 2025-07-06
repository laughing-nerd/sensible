package secret

import (
	"fmt"
	"os"
	"path/filepath"
	"sensible/internal/constants"
	"sensible/internal/utils"
	"sensible/pkg/logger"
)

func Set(key, value, env string) error {
	var (
		file      *os.File
		err       error
		dataToEnc string
	)

	// get the secrets file path
	path := fmt.Sprintf(constants.SecretsDir, env)
	secretFile, err := filepath.Abs(filepath.Join(path, constants.SecretsFile))
	if err != nil {
		return err
	}

	// create the file if it does not exist
	if !utils.FileExists(secretFile) {
		err = os.MkdirAll(filepath.Dir(secretFile), 0700) // create the directory if it does not exist
		if err != nil {
			return err
		}
		file, err = os.Create(secretFile)
		if err != nil {
			return err
		}
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close the secrets file: %w", err)
	}

	// get user password to encrypt the secret
	password, err := utils.AskPassword("Enter password to set secret: ")
	if err != nil {
		return err
	}

	// read the existing secrets file and decrypt it if not empty
	encrypted, err := os.ReadFile(secretFile)
	if err != nil {
		return err
	}

	dataToEnc = fmt.Sprintf("%s=%s", key, value)
	if len(encrypted) > 0 {
		decrypted, err := utils.Decrypt(encrypted, password)
		if err != nil {
			return err
		}
		dataToEnc = fmt.Sprintf("%s\n%s", string(decrypted), dataToEnc)
	}

	// encrypt the data and append it to the file
	encrypted, err = utils.Encrypt([]byte(dataToEnc), password)
	if err != nil {
		return err
	}

	err = os.WriteFile(secretFile, encrypted, 0600)
	if err != nil {
		return err
	}

	logger.Success("Secret set successfully!")
	return nil
}
