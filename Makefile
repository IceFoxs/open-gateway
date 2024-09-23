cur_mkfile :=$(abspath $(lastword $(MAKEFILE_LIST)))
currentPath :=$(patsubst %/, %, $(dir $(cur_mkfile)))
targetName := open-gateway
run:
	go run   $(currentPath)/cmd/opengateway/
build:
	go build -o $(currentPath)/dist/$(targetName)  $(currentPath)/cmd/opengateway/
run-with-skywalking:
	go run  -toolexec="$(currentPath)/skywalking-go/bin/skywalking-go-agent-0.5.0-darwin-amd64 -config $(currentPath)/skywalking-go/config/config.yaml" -a $(currentPath)/cmd/opengatewaysky/
build-with-skywalking:
	go build -o $(currentPath)/dist/$(targetName) -toolexec="$(currentPath)/skywalking-go/bin/skywalking-go-agent-0.5.0-darwin-amd64 -config $(currentPath)/skywalking-go/config/config.yaml" -a $(currentPath)/cmd/opengatewaysky/
clean:
	@rm -rf $(currentPath)/dist/