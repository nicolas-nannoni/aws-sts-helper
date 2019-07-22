package cmd

import (
	"github.com/nicolas-nannoni/aws-sts-helper/sts"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
	"os"
	"os/user"
	"path/filepath"
)

type GetTokenOptions struct {
	TokenCode string
	MfaArn	  string
	RoleArn	  string
	Profile	  string

	HttpPort int
	HttpPath string
}

const (
	envAwsSharedCredentialsFile	= "AWS_SHARED_CREDENTIALS_FILE"
	envAwsConfigFile			= "AWS_CONFIG_FILE"
)

func NewGetTokenCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "get-token",
		Short: "Get new temporary STS credentials.",
	}

	cmd.AddCommand(newInNewShellCmd())
	cmd.AddCommand(newAndShowExportCmd())
	cmd.AddCommand(newAndServeViaHttpCmd())

	return cmd
}

func newInNewShellCmd() *cobra.Command {

	o := GetTokenOptions{}

	cmd := &cobra.Command{
		Use:   "in-new-shell",
		Short: "Get new temporary STS credentials and update the current environment to use them.",
		Run:   o.InNewShell,
	}
	o.applyStandardFlags(cmd)

	return cmd
}

func newAndShowExportCmd() *cobra.Command {

	o := GetTokenOptions{}

	cmd := &cobra.Command{
		Use:   "and-show-export",
		Short: "Get new temporary STS credentials and print out the environment variable 'export' commands to use them.",
		Run:   o.AndShowExport,
	}
	o.applyStandardFlags(cmd)

	return cmd
}

func newAndServeViaHttpCmd() *cobra.Command {

	o := GetTokenOptions{}

	cmd := &cobra.Command{
		Use:   "and-serve-via-http",
		Short: "Get new temporary STS credentials and start an HTTP server serving the retrieved credentials for use by application compatible with EC2 IAM role retrieval (e.g. Cyberduck).",
		Run:   o.AndServeViaHttp,
	}
	o.applyStandardFlags(cmd)

	cmd.PersistentFlags().IntVar(&o.HttpPort, "port", 3000, "The URL path at which the HTTP server should expose the temporary credentials.")
	cmd.PersistentFlags().StringVar(&o.HttpPath, "path", "/credentials", "The URL path at which the HTTP server should expose the temporary credentials.")

	return cmd

}

func (o *GetTokenOptions) applyStandardFlags(cmd *cobra.Command) {

	cmd.PersistentFlags().StringVar(&o.Profile, "profile", "", "The profile name from AWS CLI config you want to use to set role-arn and/or mfa-arn")
	cmd.PersistentFlags().StringVar(&o.RoleArn, "role-arn", "", "The role ARN (Amazon Resource Name) that you want to assume (leave empty to request a session token).")
	cmd.PersistentFlags().StringVar(&o.MfaArn, "mfa-arn", "", "The MFA (Multi-Factor Authentication) device ARN (Amazon Resource Name) for your account, if applicable.")
	cmd.PersistentFlags().StringVar(&o.TokenCode, "token-code", "", "The MFA (Multi-Factor Authentication) code given by your MFA device, if applicable.")
}

func (o *GetTokenOptions) InNewShell(cmd *cobra.Command, args []string) {

	if o.Profile != "" {
		o.getRoleAndMfaFromAwsConfigFile()
	}

	sts.GetTokenAndSetEnvironment(o.RoleArn, o.MfaArn, o.TokenCode)
}

func (o *GetTokenOptions) AndShowExport(cmd *cobra.Command, args []string) {

	if o.Profile != "" {
		o.getRoleAndMfaFromAwsConfigFile()
	}

	sts.GetTokenAndReturnExportEnvironment(o.RoleArn, o.MfaArn, o.TokenCode)
}

func (o *GetTokenOptions) AndServeViaHttp(cmd *cobra.Command, args []string) {

	if o.Profile != "" {
		o.getRoleAndMfaFromAwsConfigFile()
	}

	sts.GetTokenAndServeOverHttp(o.RoleArn, o.MfaArn, o.TokenCode, o.HttpPort, o.HttpPath)
}

func (o *GetTokenOptions) getRoleAndMfaFromAwsConfigFile() {
	logrus.Debugf("Getting role_arn/mfa_serial from AWS CLI profile: %s", o.Profile)

	config_file := os.Getenv(envAwsConfigFile)
	credentials_file := os.Getenv(envAwsSharedCredentialsFile)

	user_infos, err := user.Current()
	if err != nil {
		logrus.Fatalf("error while getting user informations: %s", err)
	}

	if config_file == "" {
		config_file = filepath.FromSlash(user_infos.HomeDir + "/.aws/config")
	}
	logrus.Debugf("Path found for AWS CLI config file: %s", config_file)

	if credentials_file == "" {
		credentials_file = filepath.FromSlash(user_infos.HomeDir + "/.aws/credentials")
	}
	logrus.Debugf("Path found for AWS CLI credentials file: %s", credentials_file)

	cfg, err := ini.LooseLoad(config_file, credentials_file)
	if err != nil {
		logrus.Fatalf("error while loading AWS CLI configuration files: %s", err)
	}

	section, err := cfg.GetSection(o.Profile)
	if err != nil {
		logrus.Fatalf("can't find profile %s in AWS CLI configuration files", o.Profile)
	}

	if section.HasKey("role_arn") {
		o.RoleArn = section.Key("role_arn").Value()
		logrus.Debugf("Found role_arn key with value %s in profile", o.RoleArn)
	}

	if section.HasKey("mfa_serial") {
		o.MfaArn = section.Key("mfa_serial").Value()
		logrus.Debugf("Found mfa_serial key with value %s in profile", o.MfaArn)
	}
}
