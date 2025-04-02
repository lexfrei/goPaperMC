.PHONY: build test clean example install completion version

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X github.com/lexfrei/goPaperMC/cmd/papermc/cmd.Version=$(VERSION) -X github.com/lexfrei/goPaperMC/cmd/papermc/cmd.Commit=$(COMMIT) -X github.com/lexfrei/goPaperMC/cmd/papermc/cmd.BuildDate=$(BUILD_DATE)"

build:
	@mkdir -p bin
	@go build $(LDFLAGS) -o bin/papermc ./cmd/papermc
	@echo "Built papermc cli to bin/papermc"

test:
	@go test -v ./...

clean:
	@rm -rf bin
	@go clean
	@echo "Clean completed"

example:
	@go run ./examples/list_projects.go

download:
	@go run ./examples/download_latest.go ./

install: build
	@cp bin/papermc $(GOPATH)/bin/papermc
	@echo "Installed papermc to $(GOPATH)/bin/papermc"

completion:
	@mkdir -p completions
	@bin/papermc completion bash > completions/papermc.bash
	@bin/papermc completion zsh > completions/papermc.zsh
	@bin/papermc completion fish > completions/papermc.fish
	@bin/papermc completion powershell > completions/papermc.ps1
	@echo "Generated shell completions in completions directory"

version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

all: clean build test completion
