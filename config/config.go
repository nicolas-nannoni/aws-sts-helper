package config

var (
	Config = &ConfigEntries{}
)

type ConfigEntries struct {
	RoleArn      string
	MfaArn       string
	MfaTokenCode string

	HttpPort int
	HttpPath string

	Debug              bool
	KeepAwsEnvironment bool
}
