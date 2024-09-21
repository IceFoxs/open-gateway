package consul

import (
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hertz-contrib/registry/consul"
	"log"
)

func CreateRegistry(serverHost string, appName string, address string) (registry.Registry, error) {
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
	return r, nil
}
