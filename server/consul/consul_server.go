package consul

import (
	"fmt"
	"github.com/IceFoxs/open-gateway/server/router"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/utils"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hertz-contrib/registry/consul"
	"log"
	"os"
)

func CreateConsulServer(serverHost string, appName string, address string) (*server.Hertz, error) {
	config := consulapi.DefaultConfig()
	config.Address = address
	consulClient, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	check := &consulapi.AgentServiceCheck{
		CheckID:  appName + "-check",
		Interval: "10s",
		HTTP:     "http://" + serverHost + "/health",
	}
	r := consul.NewConsulRegister(consulClient, consul.WithCheck(check))
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(dir)
	h := server.Default(server.WithHostPorts(serverHost), server.WithRegistry(r, &registry.Info{
		ServiceName: appName,
		Addr:        utils.NewNetAddr("tcp", serverHost),
		Weight:      10,
		Tags:        nil,
	}))
	router.AddRouter(h, dir)
	return h, nil
}
