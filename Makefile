cur_mkfile := $(abspath $(lastword $(MAKEFILE_LIST)))
currentPath := $(patsubst %/, %, $(dir $(cur_mkfile)))
targetName := open-gateway
run:
	go run $(currentPath)/cmd/

build:
	go build -o $(currentPath)/dist/$(targetName) $(currentPath)/cmd/