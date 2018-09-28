.PHONY: all run build clean vet test package

APP_NAME=aws-sts-helper
APP_BUILD=`git log --pretty=format:'%h' -n 1`

APP_VERSION=1.3.0

GO_FLAGS= CGO_ENABLED=0
GO_LDFLAGS= -ldflags="-X 'github.com/nicolas-nannoni/aws-sts-helper/config.AppVersion=$(APP_VERSION)' -X 'github.com/nicolas-nannoni/aws-sts-helper/config.AppBuild=$(APP_BUILD)'"
GO_BUILD_CMD=$(GO_FLAGS) go build $(GO_LDFLAGS)
BUILD_DIR=bin
BINARY_NAME=aws-sts-helper

all: clean build package

vet:
	@go vet

test:
	@go test

build: vet test
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64
	GOOS=darwin GOARCH=amd64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_NAME)-osx-amd64
	GOOS=windows GOARCH=amd64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe

package:
	tar -C $(BUILD_DIR) -zcf $(BUILD_DIR)/$(BINARY_NAME)-$(APP_VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	tar -C $(BUILD_DIR) -zcf $(BUILD_DIR)/$(BINARY_NAME)-$(APP_VERSION)-osx-amd64.tar.gz $(BINARY_NAME)-osx-amd64
	zip -q -j  $(BUILD_DIR)/$(BINARY_NAME)-$(APP_VERSION)-windows-amd64.zip $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe

clean:
	rm -Rf $(BUILD_DIR)

