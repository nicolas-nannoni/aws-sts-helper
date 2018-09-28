package cmd

import (
	"github.com/nicolas-nannoni/aws-sts-helper/sts"
	"github.com/spf13/cobra"
)

type ClearEnvironmentOptions struct {
}

func NewClearEnvironmentCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "clear-environment",
		Short: "Unset all the credential-related AWS environment variables.",
	}

	cmd.AddCommand(newClearInNewShell())
	cmd.AddCommand(newClearAndShowUnset())

	return cmd
}

func newClearInNewShell() *cobra.Command {

	o := ClearEnvironmentOptions{}

	cmd := &cobra.Command{
		Use:   "in-new-shell",
		Short: "Get a new shell without any credential-related AWS environment variables.",
		Run:   o.InNewShell,
	}

	return cmd
}

func newClearAndShowUnset() *cobra.Command {

	o := ClearEnvironmentOptions{}

	cmd := &cobra.Command{
		Use:   "and-show-unset",
		Short: "Print out the environment variable 'unset' commands to issue to remove AWS credential-related environment variables.",
		Run:   o.AndShowUnset,
	}

	return cmd
}

func (o *ClearEnvironmentOptions) InNewShell(cmd *cobra.Command, args []string) {

	sts.ClearAwsEnvironmentInNewShell()
}

func (o *ClearEnvironmentOptions) AndShowUnset(cmd *cobra.Command, args []string) {

	sts.ClearAwsEnvironmentAndReturnUnsetEnvironment()
}
