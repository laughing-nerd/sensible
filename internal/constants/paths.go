package constants

const (
	// files
	HostFile         = "hosts.hcl"
	ValuesFile       = "values.hcl"
	SampleActionFile = "sample-action.hcl"
	SecretsFile      = "secrets.hcl.enc" // secrets is always encrypted

	// dirs
	ResourcesDir = "./.sensible/%s/resources"
	ActionsDir   = "./.sensible/%s/actions"
	SecretsDir   = "./.sensible/%s/secrets"
)
