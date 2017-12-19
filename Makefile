.PHONY: build

# VERSION is currently based on the last commit
VERSION=`git describe --tags`
COMMIT=`git rev-parse HEAD`
BUILD=`date +%FT%T%z`

all: build

build:
	go build -o _output/bin/kube-cab github.com/kube-metrics-test

clean:
	rm -rf _output

