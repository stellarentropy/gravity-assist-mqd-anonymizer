.PHONY: agent

all: anonymizer

GOOS ?= "linux"
GOARCH ?= "amd64"
TRIMPATH = "-trimpath"
BUILDVCS = "false"
LDFLAGS = "-s -w -buildid="

generate:
	go generate ./...

anonymizer: generate
	GOOS=${GOOS} GOARCH=${GOARCH} go build ${TRIMPATH} -buildvcs=${BUILDVCS} -ldflags=${LDFLAGS} -o bin/anonymizer_${GOOS}_${GOARCH} cmd/anonymizer/main.go
