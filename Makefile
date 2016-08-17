.PHONY: build fmt lint run test vendor_clean vendor_get vendor_update vet list

# Prepend our _vendor directory to the system GOPATH
# so that import path resolution will prioritize
# our third party snapshots.
GOPATH := ${PWD}/_vendor:${GOPATH}
export GOPATH

BINARY=./bin/aws-sts-helper

VERSION=0.0.1
BUILD=`git rev-parse HEAD`

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

.DEFAULT_GOAL: build

build:
	go build ${LDFLAGS} -o ${BINARY} ./*.go

run: build
	${BINARY}

osx:
	GOOS="darwin" GOARCH="amd64" go build ${LDFLAGS} -o ${BINARY}-osx ./*.go

linux:
	GOOS="linux" GOARCH="amd64" go build ${LDFLAGS} -o ${BINARY}-linux ./*.go

windows:
	GOOS="windows" GOARCH="amd64" go build ${LDFLAGS} -o ${BINARY}.exe ./*.go

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	go fmt ./...

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	golint ./src

test:
	go test ./...

install:
	go install ${LDFLAGS} -o ${BINARY} ./*.go

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	if [ -f ${BINARY}-osx ] ; then rm ${BINARY}-osx ; fi
	if [ -f ${BINARY}-linux ] ; then rm ${BINARY}-linux ; fi
	if [ -f ${BINARY}-windows ] ; then rm ${BINARY}-windows ; fi

list:
		@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs

# http://godoc.org/code.google.com/p/go.tools/cmd/vet
# go get code.google.com/p/go.tools/cmd/vet
vet:
	go vet ./*
