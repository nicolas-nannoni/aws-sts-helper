package sts

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	conf "github.com/nicolas-nannoni/aws-sts-helper/config"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/manifoldco/promptui"
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

func GetTokenAndSetEnvironment(ctx context.Context, roleArn, mfaArn, tokenCode string) {

	var credentials *types.Credentials
	if roleArn == "" {
		resp := getSessionToken(ctx, mfaArn, tokenCode)
		credentials = resp.Credentials
	} else {
		resp := getToken(ctx, roleArn, mfaArn, tokenCode)
		credentials = resp.Credentials
	}

	setEnvironmentFromStsCredentials(credentials)
	openNewShell()
}

func GetTokenAndReturnExportEnvironment(ctx context.Context, roleArn, mfaArn, tokenCode string) {

	var credentials *types.Credentials
	if roleArn == "" {
		resp := getSessionToken(ctx, mfaArn, tokenCode)
		credentials = resp.Credentials
	} else {
		resp := getToken(ctx, roleArn, mfaArn, tokenCode)
		credentials = resp.Credentials
	}

	logrus.Info("Run this command wrapped in 'eval $(aws-sts-helper get-token ...)' to automatically set your AWS environment variables.")
	fmt.Println(getSetEnvironmentStringFromCredentials(credentials))
}

func GetTokenAndReturnEnvironment(ctx context.Context, roleArn, mfaArn, tokenCode string) {

	var credentials *types.Credentials
	if roleArn == "" {
		resp := getSessionToken(ctx, mfaArn, tokenCode)
		credentials = resp.Credentials
	} else {
		resp := getToken(ctx, roleArn, mfaArn, tokenCode)
		credentials = resp.Credentials
	}

	fmt.Println(getEnvironmentStringFromCredentials(credentials))
}

func GetTokenAndServeOverHttp(ctx context.Context, roleArn, mfaArn, tokenCode string, httpPort int, httpPath string) {

	var credentials *types.Credentials
	if roleArn == "" {
		resp := getSessionToken(ctx, mfaArn, tokenCode)
		credentials = resp.Credentials
	} else {
		resp := getToken(ctx, roleArn, mfaArn, tokenCode)
		credentials = resp.Credentials
	}

	startHttpServerWithCredentials(httpPort, httpPath, credentials)
}

func ClearAwsEnvironmentInNewShell() {

	clearAwsEnvironment()
	logrus.Info("AWS environment variables unset")
	openNewShell()

}

func ClearAwsEnvironmentAndReturnUnsetEnvironment() {

	logrus.Info("Run this command wrapped in 'eval $(aws-sts-helper clear-environment ...)' to automatically set your AWS environment variables.")
	fmt.Println(getUnsetEnvironmentString())

}

func getSessionToken(ctx context.Context, mfaArn, tokenCode string) *sts.GetSessionTokenOutput {

	if !conf.Config.KeepAwsEnvironment {
		clearAwsEnvironment()
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		logrus.Fatalf("unable to get the AWS configuration: %s", err)
	}

	svc := sts.NewFromConfig(cfg)

	req := &sts.GetSessionTokenInput{}

	if mfaArn != "" {
		logrus.Debugf("using MFA device with serial %s and token code %s", mfaArn, tokenCode)
		req.SerialNumber = aws.String(mfaArn)
		req.TokenCode = aws.String(getMfaToken(tokenCode))
	}
	resp, err := svc.GetSessionToken(ctx, req)

	if err != nil {
		logrus.Fatalf("error while trying to GetSessionToken: %s", err.Error())
	}

	b, err := json.MarshalIndent(resp, "", " ")
	if err != nil {
		logrus.Fatalf("error while trying to marshall the response from the GetSessionToken call: %s", err.Error())
	}

	logrus.Infof("session token successfully received!\n%s", string(b))

	return resp
}

func getToken(ctx context.Context, roleArn, mfaArn, tokenCode string) *sts.AssumeRoleOutput {

	if !conf.Config.KeepAwsEnvironment {
		clearAwsEnvironment()
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		logrus.Fatalf("unable to get the AWS configuration: %s", err)
	}

	svc := sts.NewFromConfig(cfg)

	req := &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(getRandomSessionName()),
	}

	if mfaArn != "" {
		logrus.Debugf("using MFA device with serial %s and token code %s", mfaArn, tokenCode)
		req.SerialNumber = aws.String(mfaArn)
		req.TokenCode = aws.String(getMfaToken(tokenCode))
	}
	resp, err := svc.AssumeRole(ctx, req)

	if err != nil {
		logrus.Fatalf("error while trying to AssumeRole: %s", err.Error())
	}

	b, err := json.MarshalIndent(resp, "", " ")
	if err != nil {
		logrus.Fatalf("error while trying to marshall the response from the AssumeRole call: %s", err.Error())
	}

	logrus.Infof("token successfully received!\n%s", string(b))

	return resp
}

func getMfaToken(givenTokenCode string) string {

	if givenTokenCode != "" {
		return givenTokenCode
	}

	logrus.Debugf("no token code passed in command invocation while --mfa-arn is provided: requesting MFA token interactively")

	prompt := promptui.Prompt{
		Label: "Please type in your MFA code",
	}
	token, err := prompt.Run()
	if err != nil {
		logrus.Fatalf("error while reading the MFA code: %s", err)
	}

	return strings.TrimSpace(token)
}

func clearAwsEnvironment() {

	logrus.Debugf("clearing AWS environment variables (%s) from the current environment: %s", envAwsVariables, os.Environ())

	for _, name := range envAwsVariables {
		logrus.Debugf("clearing %s...", name)
		os.Unsetenv(name)
	}
}

func setEnvironmentFromStsCredentials(credentials *types.Credentials) {

	logrus.Debugf("setting environment variables (%s, %s, %s) based on STS output: %s", envAwsAccessKeyId, envAwsSecretAccessKey, envAwsSessionToken, credentials)

	setEnvironmentVariable(envAwsAccessKeyId, *credentials.AccessKeyId)
	setEnvironmentVariable(envAwsSecretAccessKey, *credentials.SecretAccessKey)
	setEnvironmentVariable(envAwsSessionToken, *credentials.SessionToken)

	logrus.Debugf("Environment: %s", os.Environ())
}

func setEnvironmentVariable(key string, value string) {

	logrus.Debugf("setting environment variable %s to %s", key, value)
	err := os.Setenv(key, value)
	if err != nil {
		logrus.Fatalf("error while setting environment variable %s: %s", key, err)
	}
}

func getSetEnvironmentStringFromCredentials(credentials *types.Credentials) string {

	setCmds := make([]string, 3)
	setCmds = append(setCmds, getExportEnvironmentString(envAwsAccessKeyId, *credentials.AccessKeyId))
	setCmds = append(setCmds, getExportEnvironmentString(envAwsSecretAccessKey, *credentials.SecretAccessKey))
	setCmds = append(setCmds, getExportEnvironmentString(envAwsSessionToken, *credentials.SessionToken))
	return strings.Join(setCmds, "\n")
}

func getEnvironmentStringFromCredentials(credentials *types.Credentials) string {

	setCmds := make([]string, 3)
	setCmds = append(setCmds, getEnvironmentString(envAwsAccessKeyId, *credentials.AccessKeyId))
	setCmds = append(setCmds, getEnvironmentString(envAwsSecretAccessKey, *credentials.SecretAccessKey))
	setCmds = append(setCmds, getEnvironmentString(envAwsSessionToken, *credentials.SessionToken))
	return strings.Join(setCmds, "\n")
}

func getExportEnvironmentString(key string, value string) string {
	return fmt.Sprintf("export %s=%s", key, value)
}

func getEnvironmentString(key string, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
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

func startHttpServerWithCredentials(httpPort int, httpPath string, credentials *types.Credentials) error {

	iamBody := IamRoleResponse{
		Code:            "Success",
		LastUpdate:      time.Now().Format(time.RFC3339),
		Type:            "AWS-HMAC",
		AccessKeyId:     *credentials.AccessKeyId,
		SecretAccessKey: *credentials.SecretAccessKey,
		Token:           *credentials.SessionToken,
		Expiration:      (*credentials.Expiration).Format(time.RFC3339),
	}

	json, err := json.Marshal(iamBody)
	if err != nil {
		return fmt.Errorf("unable to convert the STS credentials into a JSON EC2 IAM Role output: %s", err)
	}

	logrus.Infof("Starting server on port %d, exposing credentials at %s", httpPort, httpPath)

	http.HandleFunc(httpPath, serveCredentials(json))
	http.ListenAndServe(fmt.Sprintf("localhost:%d", httpPort), nil)

	return nil
}

func serveCredentials(jsonIamBody []byte) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		logrus.Infof("Received credentials request: %s", r.UserAgent())
		logrus.Debugf("received request: %s", *r)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonIamBody)
	}
}
