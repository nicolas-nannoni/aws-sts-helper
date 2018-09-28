package main

import (
	"github.com/nicolas-nannoni/aws-sts-helper/cmd"
	"github.com/sirupsen/logrus"
)

func main() {

	if err := cmd.NewRootCommand().Execute(); err != nil {
		logrus.Fatalf("error while executing aws-sts-helper: %s", err)
	}
}
