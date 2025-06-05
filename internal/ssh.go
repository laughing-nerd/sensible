package internal

import "golang.org/x/crypto/ssh"

func SshIntoHost(hostname, port, authType, username, creds string) (*ssh.Client, error) {
	client, err := ssh.Dial("tcp", hostname+":"+port, &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(creds),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})

	return client, err
}
