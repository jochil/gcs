.PHONY: test
test: 
	go test -v ./pkg/...

.PHONY: build 
build: 
	go build -v

.PHONY: lint
lint: 
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.2
	golangci-lint run

.PHONY: fmt
fmt:
	go fmt ./pkg/...
	go fmt main.go 

.PHONY: coverage
coverage: 
	mkdir -p .cov
	-go test ./pkg/... -cover -args -test.gocoverdir="${PWD}/.cov"
	go tool covdata func -i=.cov/

FUNC?=CycloA
.PHONY: coverage
graph/test:
	go test ./pkg/parser/graph_test.go -v 
	dot -Tsvg -O .draw/$(FUNC).gv 
	firefox-developer-edition --new-tab .draw/$(FUNC).gv.svg  &
