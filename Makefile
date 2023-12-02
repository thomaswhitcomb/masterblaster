all: clean compile test

clean:
	@echo "==> Cleaning up previous builds."
	@rm -rf ./bin/masterblaster

compile:
	@echo "==> Compiling source code."
	@go build -v -o ./bin/masterblaster 

coverage:
	@go test -coverprofile cover.out
	@go tool cover -html=cover.out

deps:
	@echo "==> Downloading dependencies."
	@godep save $(shell go list ./... | grep -v /vendor/)

fmt:
	@echo "==> Formatting source code."
	@gofmt -w ./

race_compile:
	@echo "==> Compiling source code."
	@go build -v -race -o ./bin/masterblaster $(shell go list ./... | grep -v /vendor/)

test: fmt vet
	@echo "==> Running tests."
	@go test -cover $(shell go list ./... | grep -v /vendor/)
	@echo "==> Tests complete."

vet:
	@go vet $(shell go list ./... | grep -v /vendor/)

help:
	@echo "clean\t\tremove previous builds"
	@echo "compile\t\tbuild the code"
	@echo "coverage\tgenerate and view code coverage"
	@echo "deps\t\tdownload dependencies"
	@echo "fmt\t\tformat the code"
	@echo "race_compile\tbuild the code with race detection"
	@echo "test\t\ttest the code"
	@echo "vet\t\tvet the code"
	@echo ""
	@echo "default will test, format, and compile the code"

.PNONY: all clean deps fmt help test
