package kong

import (
	"sensible/internal/secret"
)

type SecretCommand struct {
	Set    SecretSetCmd    `cmd:"" help:"Set a secret"`
	Get    SecretGetCmd    `cmd:"" help:"Get a secret"`
	Remove SecretRemoveCmd `cmd:"" help:"Delete a secret"`
}

// set secret
type SecretSetCmd struct {
	CommonFlags
	Key   string `help:"secret key" short:"k"`
	Value string `help:"secret value" short:"v"`
}

func (set *SecretSetCmd) Run() error {
	return secret.Set(set.Key, set.Value, set.Env)
}

// get secret
type SecretGetCmd struct {
	CommonFlags
	Key string `help:"secret key whose value needs to be read" short:"k"`
	All bool   `help:"Returns all secret key=value pairs" short:"a"`
}

func (get *SecretGetCmd) Run() error {
	return secret.Get(get.Key, get.All, get.Env)
}

// remove secret
type SecretRemoveCmd struct {
	CommonFlags
	Key string `help:"secret key whose value needs to be deleted" short:"k"`
}

func (rem *SecretRemoveCmd) Run() error {
	return secret.Remove(rem.Key, rem.Env)
}
