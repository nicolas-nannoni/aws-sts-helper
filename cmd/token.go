package cmd

import (
	"github.com/nicolas-nannoni/aws-sts-helper/sts"
	"github.com/spf13/cobra"
)

type GetTokenOptions struct {
	TokenCode string
	MfaArn    string
	RoleArn   string

	HttpPort int
	HttpPath string
}

func NewGetTokenCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "get-token",
		Short: "Get new temporary STS credentials.",
	}

	cmd.AddCommand(newInNewShellCmd())
	cmd.AddCommand(newAndShowExportCmd())
	cmd.AddCommand(newAndServeViaHttpCmd())
	cmd.AddCommand(newAndShowEnvCmd())

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

func newAndShowEnvCmd() *cobra.Command {

	o := GetTokenOptions{}

	cmd := &cobra.Command{
		Use:   "and-show-environment",
		Short: "Get new temporary STS credentials and print out the environment variables to use them.",
		Run:   o.AndShowEnv,
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

	cmd.PersistentFlags().StringVar(&o.RoleArn, "role-arn", "", "The role ARN (Amazon Resource Name) that you want to assume (leave empty to request a session token).")
	cmd.PersistentFlags().StringVar(&o.MfaArn, "mfa-arn", "", "The MFA (Multi-Factor Authentication) device ARN (Amazon Resource Name) for your account, if applicable.")
	cmd.PersistentFlags().StringVar(&o.TokenCode, "token-code", "", "The MFA (Multi-Factor Authentication) code given by your MFA device, if applicable.")
}

func (o *GetTokenOptions) InNewShell(cmd *cobra.Command, args []string) {

	sts.GetTokenAndSetEnvironment(cmd.Context(), o.RoleArn, o.MfaArn, o.TokenCode)
}

func (o *GetTokenOptions) AndShowExport(cmd *cobra.Command, args []string) {

	sts.GetTokenAndReturnExportEnvironment(cmd.Context(), o.RoleArn, o.MfaArn, o.TokenCode)
}

func (o *GetTokenOptions) AndShowEnv(cmd *cobra.Command, args []string) {

	sts.GetTokenAndReturnEnvironment(cmd.Context(), o.RoleArn, o.MfaArn, o.TokenCode)
}

func (o *GetTokenOptions) AndServeViaHttp(cmd *cobra.Command, args []string) {

	sts.GetTokenAndServeOverHttp(cmd.Context(), o.RoleArn, o.MfaArn, o.TokenCode, o.HttpPort, o.HttpPath)
}
