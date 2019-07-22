package cmd

import (
	"github.com/nicolas-nannoni/aws-sts-helper/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func NewRootCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:     filepath.Base(os.Args[0]) + " <command>",
		Short:   "Tool to make your life easier when using Amazon STS (Security Token Service).",
		Version: config.AppVersion + " (build " + config.AppBuild + ")",

		PersistentPreRun: initApp,
	}
	cmd.AddCommand(NewGetTokenCmd())
	cmd.AddCommand(NewClearEnvironmentCmd())

	cmd.PersistentFlags().BoolVar(&config.Config.Debug, "debug", false, "Enable debug mode.")
	cmd.PersistentFlags().BoolVar(&config.Config.KeepAwsEnvironment, "keep-aws-environment", false,
		"Keep the AWS environment variables set (e.g. AWS_ACCESS_KEY_ID) before creating AWS sessions.")

	return cmd
}

func initApp(cmd *cobra.Command, args []string) {

	if config.Config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debug mode enabled!")
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	rand.Seed(time.Now().UTC().UnixNano())

}
