export REPOSITORY=flytestdlib
include boilerplate/lyft/golang_test_targets/Makefile

# Generate golden files. Add test packages that generate golden files here.
golden:
	go test ./cli/pflags/api -update
	go test ./config -update
	go test ./storage -update
	go test ./tests -update


generate:
	@echo "************************ go generate **********************************"
	go generate ./...

# This is the only target that should be overriden by the project. Get your binary into ${GOREPO}/bin
.PHONY: compile
compile:
	mkdir -p ./bin
	go build -o pflags ./cli/pflags/main.go && mv ./pflags ./bin

generate:
	which pflags || (go get github.com/lyft/flytestdlib/cli/pflags)
	which mockery || (go get github.com/enghabu/mockery/cmd/mockery)
	@go generate ./...

.PHONY: test_unit_codecov
test_unit_codecov:
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic; curl -s https://codecov.io/bash > codecov_bash.sh; bash codecov_bash.sh

