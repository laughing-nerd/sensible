package models

type HostConfig struct {
	Groups []Group `hcl:"group,block"`
	Hosts  []Host  `hcl:"host,block"`
}

type Group struct {
	Name  string `hcl:"group,label"`
	Hosts []Host `hcl:"host,block"`
}

type Host struct {
	Name     string `hcl:"host,label"`
	Address  string `hcl:"address,optional"`
	Username string `hcl:"username"`
	Password string `hcl:"password"`
}
