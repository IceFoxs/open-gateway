cur_mkfile :=$(abspath $(lastword $(MAKEFILE_LIST)))
currentPath :=$(patsubst %/, %, $(dir $(cur_mkfile)))
targetName := opengateway
run:
	go run   $(currentPath)/cmd/opengateway/
build:
	go build  -trimpath -o $(currentPath)/dist/$(targetName)  $(currentPath)/cmd/opengateway/
run-with-skywalking-macos:
	go run  -toolexec="$(currentPath)/skywalking-go/bin/skywalking-go-agent-0.5.0-darwin-amd64 -config $(currentPath)/skywalking-go/config/config.yaml" -a $(currentPath)/cmd/opengatewaysky/
build-with-skywalking-macos:
	go build -trimpath -o $(currentPath)/dist/$(targetName) -toolexec="$(currentPath)/skywalking-go/bin/skywalking-go-agent-0.5.0-darwin-amd64 -config $(currentPath)/skywalking-go/config/config.yaml" -a $(currentPath)/cmd/opengatewaysky/
run-with-skywalking-linux:
	go run -toolexec="$(currentPath)/skywalking-go/bin/skywalking-go-agent-0.5.0-linux-amd64 -config $(currentPath)/skywalking-go/config/config.yaml" -a $(currentPath)/cmd/opengatewaysky/
build-with-skywalking-linux:
	go build -trimpath -o $(currentPath)/dist/$(targetName) -toolexec="$(currentPath)/skywalking-go/bin/skywalking-go-agent-0.5.0-linux-amd64 -config $(currentPath)/skywalking-go/config/config.yaml" -a $(currentPath)/cmd/opengatewaysky/
clean:
	@rm -rf $(currentPath)/dist/