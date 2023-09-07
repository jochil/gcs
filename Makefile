.PHONY: test
test: 
	go test -v ./pkg/...

.PHONY: fmt
fmt:
	go fmt ./pkg/...
	go fmt main.go 

.PHONY: coverage
coverage: 
	mkdir -p .cov
	-go test ./pkg/... -cover -args -test.gocoverdir="${PWD}/.cov"
	go tool covdata func -i=.cov/
