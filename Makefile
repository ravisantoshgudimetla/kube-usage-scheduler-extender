.PHONY: build

# VERSION is currently based on the last commit
VERSION=`git describe --tags`
COMMIT=`git rev-parse HEAD`
BUILD=`date +%FT%T%z`
GO_TESTS=`go list github.com/kube-usage-scheduler-extender/... | grep -v github.com/kube-usage-scheduler-extender/vendor/`

all: build

build:
	go build -o _output/bin/kube-ext github.com/kube-usage-scheduler-extender
test:
	go test ${GO_TESTS}

clean:
	rm -rf _output

