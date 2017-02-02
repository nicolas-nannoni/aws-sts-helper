package sts

import (
	"../config"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/urfave/cli"
	"github.com/segmentio/go-prompt"
)

const (
	envAwsAccessKeyId     = "AWS_ACCESS_KEY_ID"
	envAwsAccessKey       = "AWS_ACCESS_KEY"
	envAwsSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	envAwsSecretKey       = "AWS_SECRET_KEY"
	envAwsSessionToken    = "AWS_SESSION_TOKEN"

	randomSessionNameLength = 8
)

var (
	envAwsVariables = []string{envAwsAccessKeyId, envAwsAccessKey, envAwsSecretAccessKey, envAwsSecretKey, envAwsSessionToken}
)

func GetTokenAndSetEnvironment(c *cli.Context) error {

	resp := getToken(c)

	setEnvironmentFromStsReponse(resp)
	openNewShell()

	return nil
}

func GetTokenAndReturnExportEnvironment(c *cli.Context) error {

	resp := getToken(c)

	log.Info("Run this command wrapped in 'eval $(aws-sts-helper get-token ...)' to automatically set your AWS environment variables.")
	fmt.Println(getSetEnvironmentString(resp))

	return nil
}

func ClearAwsEnvironmentInNewShell(c *cli.Context) error {

	clearAwsEnvironment()
	log.Info("AWS environment variables unset")
	openNewShell()

	return nil
}

func ClearAwsEnvironmentAndReturnUnsetEnvironment(c *cli.Context) error {

	log.Info("Run this command wrapped in 'eval $(aws-sts-helper clear-environment ...)' to automatically set your AWS environment variables.")
	fmt.Println(getUnsetEnvironmentString())

	return nil
}

func getToken(c *cli.Context) *sts.AssumeRoleOutput {

	if !config.Config.KeepAwsEnvironment {
		clearAwsEnvironment()
	}

	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("Unable to open an AWS session: %s", err)
	}

	svc := sts.New(sess)

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(config.Config.RoleArn),
		RoleSessionName: aws.String(getRandomSessionName()),
	}

	if config.Config.MfaArn != "" {
		log.Debugf("Using MFA device with serial %s and token code %s", config.Config.MfaArn, config.Config.MfaTokenCode)
		params.SerialNumber = aws.String(config.Config.MfaArn)
		params.TokenCode = aws.String(getMfaToken())
	}
	resp, err := svc.AssumeRole(params)

	if err != nil {
		log.Fatalf("Error while trying to AssumeRole: %s", err.Error())
	}

	log.Infof("Token successfully received!\n%s", resp.GoString())

	return resp
}

func getMfaToken() string {

	if config.Config.MfaTokenCode != "" {
		return config.Config.MfaTokenCode
	}

	log.Debugf("No token code passed in command invocation while --mfa-arn is provided: requesting MFA token interactively")
	return prompt.StringRequired("Please type in your MFA code")
}

func clearAwsEnvironment() {

	log.Debugf("Clearing AWS environement variables (%s) from the current environment: %s", envAwsVariables, os.Environ())

	for _, name := range envAwsVariables {
		log.Debugf("Clearing %s...", name)
		os.Unsetenv(name)
	}
}

func setEnvironmentFromStsReponse(resp *sts.AssumeRoleOutput) {

	log.Debugf("Setting environment variables (%s, %s, %s) based on STS output: %s", envAwsAccessKey, envAwsSecretKey, envAwsSessionToken, resp)

	setEnvironmentVariable(envAwsAccessKey, *resp.Credentials.AccessKeyId)
	setEnvironmentVariable(envAwsSecretKey, *resp.Credentials.SecretAccessKey)
	setEnvironmentVariable(envAwsSessionToken, *resp.Credentials.SessionToken)

	log.Debugf("Environment: %s", os.Environ())
}

func setEnvironmentVariable(key string, value string) {

	log.Debugf("Setting environment variable %s to %s", key, value)
	err := os.Setenv(key, value)
	if err != nil {
		log.Fatalf("Error while setting environment variable %s: %s", key, err)
	}
}

func getSetEnvironmentString(resp *sts.AssumeRoleOutput) string {

	setCmds := make([]string, 3)
	setCmds = append(setCmds, getExportEnvironmentString(envAwsAccessKeyId, *resp.Credentials.AccessKeyId))
	setCmds = append(setCmds, getExportEnvironmentString(envAwsSecretAccessKey, *resp.Credentials.SecretAccessKey))
	setCmds = append(setCmds, getExportEnvironmentString(envAwsSessionToken, *resp.Credentials.SessionToken))
	return strings.Join(setCmds, "\n")
}

func getExportEnvironmentString(key string, value string) string {
	return fmt.Sprintf("export %s=%s", key, value)
}

func getUnsetEnvironmentString() string {

	unsetCmds := make([]string, len(envAwsVariables))
	for _, variable := range envAwsVariables {
		unsetCmds = append(unsetCmds, fmt.Sprintf("unset %s", variable))
	}
	return strings.Join(unsetCmds, "\n")
}

func getRandomSessionName() string {

	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, randomSessionNameLength)
	for i := 0; i < randomSessionNameLength; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return fmt.Sprintf("aws-sts-helper-%s", string(result))
}

func openNewShell() {
	syscall.Exec(os.Getenv("SHELL"), []string{os.Getenv("SHELL")}, syscall.Environ())
}
