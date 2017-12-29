.PHONY: build

# VERSION is currently based on the last commit
VERSION=`git describe --tags`
COMMIT=`git rev-parse HEAD`
BUILD=`date +%FT%T%z`
GO_TESTS=`go list github.com/kube-cab/... | grep -v github.com/kube-cab/vendor/`

all: build

build:
	go build -o _output/bin/kube-cab github.com/kube-cab
test:
	go test ${GO_TESTS}

clean:
	rm -rf _output

