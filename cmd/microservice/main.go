// +build real

package main

import (
	"os"

	"gitea-devops-shared-threefld-cicd.apps.c4.us-east-1.dev.aws.ocp.14west.io/cicd/trackmate-message-consumer/pkg/connectors"
	"gitea-devops-shared-threefld-cicd.apps.c4.us-east-1.dev.aws.ocp.14west.io/cicd/trackmate-message-consumer/pkg/handlers"
	"gitea-devops-shared-threefld-cicd.apps.c4.us-east-1.dev.aws.ocp.14west.io/cicd/trackmate-message-consumer/pkg/validator"
	"github.com/microlib/simple"
)

var (
	logger *simple.Logger
)

// Main function : keep things clean and simple
// Allows for simple E2E testing and code coverage
func main() {

	if os.Getenv("LOG_LEVEL") == "" {
		logger = &simple.Logger{Level: "info"}
	} else {
		logger = &simple.Logger{Level: os.Getenv("LOG_LEVEL")}
	}
	err := validator.ValidateEnvars(logger)
	if err != nil {
		os.Exit(-1)
	}
	conn := connectors.NewClientConnectors(logger)
	handlers.Init(conn)
	//defer conn.Close()
}
