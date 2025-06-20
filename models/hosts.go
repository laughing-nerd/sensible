package models

import "golang.org/x/crypto/ssh"

type HostConfig struct {
	Groups []Group `hcl:"group,block"`
	Hosts  []Host  `hcl:"host,block"`
}

type Group struct {
	Name  string `hcl:"group,label"`
	Hosts []Host `hcl:"host,block"`
}

type Host struct {
	Name       string `hcl:"host,label"`
	Address    string `hcl:"address"`
	Username   string `hcl:"username"`
	Password   string `hcl:"password,optional"`
	PrivateKey string `hcl:"private_key,optional"`
	Timeout    int    `hcl:"timeout,optional"`
	SshClient  *ssh.Client
}
