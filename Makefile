include *.mf

# Default value for params. You need define GO_BIN in *.mf file.
ifndef GO_BIN
override GO_BIN = "test_http"
endif

VERSION := $(shell git describe --tags 2> /dev/null || echo no-tag)
BRANCH := $(shell git symbolic-ref -q --short HEAD)
COMMIT := $(shell git rev-parse HEAD)

# Default value for params. You need define PROJECT_NAME(like GOMODULE name) in *.mf file.
ifndef PROJECT_NAME
override PROJECT_NAME = github.com/AlexandrGurkin/http_template
endif

# Default value for params. You need define PROJECT_NAME(like GOMODULE name) in *.mf file.
ifndef MAIN_PATH
override MAIN_PATH = ./test/main.go
endif

# Use linker flags to provide version/build settings
# https://stackoverflow.com/questions/47509272/how-to-set-package-variable-using-ldflags-x-in-golang-build
LDFLAGS := -X $(PROJECT_NAME)/internal/ver.version=$(VERSION) -X $(PROJECT_NAME)/internal/ver.commit=$(COMMIT) -X $(PROJECT_NAME)/internal/ver.branch=$(BRANCH) -X $(PROJECT_NAME)/internal/ver.buildTime=`date '+%Y-%m-%d_%H:%M:%S_%Z'`

BUILD_ARG = -ldflags "$(LDFLAGS)" $(MAIN_PATH)

trash:
	@echo $(GO_BIN) $(BUILD_ARG)

download: ##Download go.mod dependencies
	@echo Download go.mod dependencies
	@go mod download
	@echo Download completed

swagger: ##Run swagger generation
	@echo Delete generated files
	@rm -rf restapi/operations restapi/doc.go restapi/embedded_spec.go restapi/server.go models
	@echo Delete completed
	@echo Code generation
	@docker run --rm -it -e GOPATH=/go -v $$(pwd):/work -w /work quay.io/goswagger/swagger:v0.25.0 generate server --exclude-main -f "./api/swagger.yaml"
	@echo Generation completed

build: ## Build app
	@echo Build project
	@go build -o $(GO_BIN) $(BUILD_ARG)
	@echo Build completed

build-full: swagger build