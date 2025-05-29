.DEFAULT_GOAL:=help

.PHONY: help clean build

help:  ## Display this help and any documented user-facing targets. Other undocumented targets may be present in the Makefile.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

clean:  ## Clean
	rm bootstrap myFUnction.zip

build:  ## Build
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go
	zip myFunction.zip bootstrap
	