package main

import (
	"fmt"
	"github.com/IceFoxs/open-gateway/cache/appmetadata"
	"github.com/IceFoxs/open-gateway/cache/gatewayconfig"
	"github.com/IceFoxs/open-gateway/cache/gatewaymethod"
	"github.com/IceFoxs/open-gateway/cache/gatewaysystem"
	"github.com/IceFoxs/open-gateway/client"
	"github.com/IceFoxs/open-gateway/conf"
	"github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/db"
	"github.com/IceFoxs/open-gateway/registry"
	"github.com/IceFoxs/open-gateway/server/consul"
	na "github.com/IceFoxs/open-gateway/server/nacos"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func main() {
	address := conf.GetConf().Registry.RegistryAddress[0]
	username := conf.GetConf().Registry.Username
	password := conf.GetConf().Registry.Password
	register := conf.GetConf().App.Register
	host := conf.GetConf().App.Host
	appName := conf.GetConf().App.Name
	dsn := conf.GetConf().MySQL.DSN
	db.Init(dsn)
	registry.GetRegisterClient()
	c := gatewayconfig.GetGatewayConfigCache()
	c.RefreshCache()
	gsc := gatewaysystem.GetGatewaySystemCache()
	gsc.RefreshCache()
	amc := appmetadata.GetAppMetadataCache()
	amc.RefreshAllCache()
	gmc := gatewaymethod.GetGatewayMethodCache()
	gmc.RefreshAllCache()
	client.GetHttpClient()
	if register == "" {
		panic("app register can not empty, please check your config")
	}

	h, err := CreateServer(register, host, appName, address, username, password)
	if err != nil {
		hlog.SystemLogger().Errorf("create server failed: %s", err.Error())
	}
	h.Spin()

}

func CreateServer(register string, host string, appName string, address string, username string, password string) (*server.Hertz, error) {
	if register == constant.REGISTRY_NACOS {
		return na.CreateNacosServer(host, appName, address, username, password)
	}
	if register == constant.REGISTRY_CONSUL {
		return consul.CreateConsulServer(host, appName, address)
	}
	return nil, fmt.Errorf("register %s is not supported", register)
}
