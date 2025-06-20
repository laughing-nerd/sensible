package connectors

import (
	"errors"
	"os"
	"path/filepath"
	"sensible/internal/constants"
	"time"

	"golang.org/x/crypto/ssh"
)

func NewSshConnection(address string, authType, username, creds string, timeout int) (*ssh.Client, error) {
	var sshTimeout = 10 * time.Second
	if timeout != 0 {
		sshTimeout = time.Duration(timeout) * time.Second
	}

	// select the auth method
	var authMethod ssh.AuthMethod
	switch authType {
	case constants.Password:
		authMethod = ssh.Password(creds)
	case constants.PrivateKey:
		signer, err := parsePrivateKey(creds)
		if err != nil {
			return nil, err
		}
		authMethod = ssh.PublicKeys(signer)
	}

	return ssh.Dial("tcp", address, &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         sshTimeout,
	})
}

func NewSshSession(client *ssh.Client) (*ssh.Session, error) {
	if client == nil {
		return nil, errors.New("SSH client is nil")
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	return session, nil
}

// helper func ...
// parsePrivateKey reads the private key from the given path and returns an ssh.Signer
func parsePrivateKey(keyPath string) (ssh.Signer, error) {
	file, err := filepath.Abs(keyPath)
	if err != nil {
		return nil, err
	}
	keyBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// try unencrypted key
	// if the key is encrypted then it won't work since encrypted keys are not supported yet
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}

	return signer, nil
}
