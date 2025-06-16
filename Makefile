current_os := $(shell uname -s | tr '[:upper:]' '[:lower:]')
bin_ext :=
ifeq ($(OS),Windows_NT)
	current_os = windows
	bin_ext = .exe
endif

.PHONY: deps 
deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.2

.PHONY: setup
setup: deps
	mkdir -p .draw

.PHONY: build 
build: 
	go build -v -o build/bin/dlth_$(current_os)$(bin_ext)

.PHONY: install
install: 
	go install

.PHONY: test
test: setup
	go test -v ./pkg/...

.PHONY: lint
lint: deps 
	golangci-lint run

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: coverage
coverage: 
	mkdir -p .cov
	-go test ./pkg/... -cover -args -test.gocoverdir="${PWD}/.cov"
	go tool covdata func -i=.cov/

FUNC?=CycloA
.PHONY: graph/test
graph/test: setup
	-go test ./pkg/parser/graph_test.go -v -failfast
	dot -Tsvg -O .draw/$(FUNC).gv 
	firefox-developer-edition --new-tab .draw/$(FUNC).gv.svg  &

