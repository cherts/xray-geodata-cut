APPNAME = xray-geodata-cut
APPOS = ${GOOS}
LDFLAGS = -a

MODERNIZE_CMD = go run golang.org/x/tools/go/analysis/passes/modernize/cmd/modernize@latest

.PHONY: help \
		clean lint test race \
		build go-update \
		modernize modernize-fix modernize-check

.DEFAULT_GOAL := help

help: ## Display this help screen
	@echo "Makefile available targets:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  * \033[36m%-25s\033[0m %s\n", $$1, $$2}'

clean: ## Clean
	rm -f ./${APPNAME}

go-update: ## Update go mod
	go mod tidy -compat=1.26
	go get -u ./
	go mod download
	go get -u ./
	go mod download

dep: ## Get the dependencies
	go mod download

lint: ## Lint the source files
	go env -w GOFLAGS="-buildvcs=false"
	go vet ./...
	REVIVE_FORCE_COLOR=1 revive -formatter friendly ./...
	gosec -quiet ./...

test: dep lint ## Run tests
	go test -race -timeout 300s -coverprofile=.test_coverage.txt ./... && \
    	go tool cover -func=.test_coverage.txt | tail -n1 | awk '{print "Total test coverage: " $$3}'
	@rm .test_coverage.txt

race: dep ## Run data race detector
	go test -race -short -timeout 300s -p 1 ./...

build: dep ## Build
	CGO_ENABLED=0 GOOS=${APPOS} GOARCH=${GOARCH} go build ${LDFLAGS} -o ${APPNAME} ./

modernize: modernize-fix ## Run gopls modernize check and fix

modernize-fix: ## Run gopls modernize fix
	@echo "Running gopls modernize with -fix..."
	go env -w GOFLAGS="-buildvcs=false"
	$(MODERNIZE_CMD) -test -fix ./...

modernize-check: ## Run gopls modernize only check
	@echo "Checking if code needs modernization..."
	go env -w GOFLAGS="-buildvcs=false"
	$(MODERNIZE_CMD) -test ./...