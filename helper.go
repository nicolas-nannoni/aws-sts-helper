package main

import (
	"./config"
	"./sts"

	"math/rand"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

// App build information
var (
	Version string
	Build   string
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
		cli.Author{
			Name:  "Nicolas Nannoni",
			Email: "nannoni@kth.se",
		},
	}
	app.Usage = "Tool to make your life easier when using Amazon STS (Security Token Service)."
	app.Before = initApp
	app.Version = Version + " (build " + Build + ")"

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
