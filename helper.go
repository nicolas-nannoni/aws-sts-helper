package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/nicolas-nannoni/aws-sts-helper/config"
	"github.com/nicolas-nannoni/aws-sts-helper/sts"
	"github.com/urfave/cli"
)

// App build information
var (
	AppVersion string
	AppBuild   string
)

func initApp(c *cli.Context) error {

	if config.Config.Debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug mode enabled!")
	} else {
		log.SetLevel(log.InfoLevel)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	return nil
}

func appDefinition() (app *cli.App) {

	app = cli.NewApp()
	app.Name = "aws-sts-helper"
	app.Authors = []cli.Author{
		{
			Name:  "Nicolas Nannoni",
			Email: "nannoni@kth.se",
		},
	}
	app.Usage = "Tool to make your life easier when using Amazon STS (Security Token Service)."
	app.Before = initApp
	app.Version = AppVersion + " (build " + AppBuild + ")"

	app.Commands = []cli.Command{
		{
			Name:  "get-token",
			Usage: "Get new temporary STS credentials.",
			Subcommands: []cli.Command{
				{
					Name:   "in-new-shell",
					Usage:  "Get new temporary STS credentials and update the current environment to use them.",
					Action: sts.GetTokenAndSetEnvironment,
				},
				{
					Name:   "and-show-export",
					Usage:  "Get new temporary STS credentials and print out the environment variable 'export' commands to use them.",
					Action: sts.GetTokenAndReturnExportEnvironment,
				},
				{
					Name:   "and-serve-via-http",
					Usage:  "Get new temporary STS credentials and start an HTTP server serving the retrieved credentials for use by application compatible with EC2 IAM role retrieval (e.g. Cyberduck).",
					Action: sts.GetTokenAndServeOverHttp,
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:        "port",
							Value:       3000,
							Usage:       "The port on which the HTTP server should expose the temporary credentials.",
							Destination: &config.Config.HttpPort,
						},
						cli.StringFlag{
							Name:        "path",
							Value:       "/credentials",
							Usage:       "The URL path at which the HTTP server should expose the temporary credentials.",
							Destination: &config.Config.HttpPath,
						},
					},
				},
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "token-code",
					Value:       "",
					Usage:       "The MFA (Multi-Factor Authentication) code given by your MFA device, if applicable.",
					Destination: &config.Config.MfaTokenCode,
				},
				cli.StringFlag{
					Name:        "mfa-arn",
					Value:       "",
					Usage:       "The MFA (Multi-Factor Authentication) device ARN (Amazon Resource Name) for your account, if applicable.",
					Destination: &config.Config.MfaArn,
				},
				cli.StringFlag{
					Name:        "role-arn",
					Value:       "",
					Usage:       "The role ARN (Amazon Resource Name) that you want to assume.",
					Destination: &config.Config.RoleArn,
				},
			},
		},
		{
			Name:  "clear-environment",
			Usage: "Unset all the credential-related AWS environment variables.",
			Subcommands: []cli.Command{
				{
					Name:   "in-new-shell",
					Usage:  "Get a new shell without any credential-related AWS environment variables.",
					Action: sts.ClearAwsEnvironmentInNewShell,
				},
				{
					Name:   "and-show-unset",
					Usage:  "Print out the environment variable 'unset' commands to issue to remove AWS credential-related environment variables.",
					Action: sts.ClearAwsEnvironmentAndReturnUnsetEnvironment,
				},
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enable debug mode.",
			Destination: &config.Config.Debug,
		},
		cli.BoolFlag{
			Name:        "keep-aws-environment",
			Usage:       "Keep the AWS environment variables set (e.g. AWS_ACCESS_KEY_ID) before creating AWS sessions.",
			Destination: &config.Config.Debug,
		},
	}

	return
}

func main() {
	appDefinition().Run(os.Args)
}
