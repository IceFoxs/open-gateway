package main

import (
	"github.com/IceFoxs/open-gateway/conf"
	"github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/server/consul"
	na "github.com/IceFoxs/open-gateway/server/nacos"
)

func main() {
	address := conf.GetConf().Registry.RegistryAddress[0]
	username := conf.GetConf().Registry.Username
	password := conf.GetConf().Registry.Password
	register := conf.GetConf().App.Register
	host := conf.GetConf().App.Host
	appName := conf.GetConf().App.Name
	if register == "" {
		panic("app register is empty")
	}
	if register == constant.REGISTRY_NACOS {
		na.CreateNacosServer(host, appName, address, username, password)
	}
	if register == constant.REGISTRY_CONSUL {
		consul.CreateConsulServer(host, appName, address)
	}
}
