package config

var (
	Config = &ConfigEntries{}

	// App build information
	AppVersion string
	AppBuild   string
)

type ConfigEntries struct {
	Debug              bool
	KeepAwsEnvironment bool
}
