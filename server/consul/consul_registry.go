package consul

import (
	"github.com/IceFoxs/open-gateway/conf"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hertz-contrib/registry/consul"
)

func CreateRegistry() (registry.Registry, error) {
	config := consulapi.DefaultConfig()
	config.Address = conf.GetConf().Consul.Address[0]
	appName := conf.GetConf().App.Name
	host := conf.GetConf().App.Host
	consulClient, err := consulapi.NewClient(config)
	if err != nil {
		hlog.Fatal(err)
		return nil, err
	}
	check := &consulapi.AgentServiceCheck{
		CheckID:  appName + "-check",
		Interval: "10s",
		HTTP:     "http://" + host + "/health",
	}
	r := consul.NewConsulRegister(consulClient, consul.WithCheck(check))
	return r, nil
}
