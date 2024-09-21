package main

import (
	"fmt"
	"github.com/IceFoxs/open-gateway/cache/appmetadata"
	"github.com/IceFoxs/open-gateway/cache/gatewayconfig"
	"github.com/IceFoxs/open-gateway/cache/gatewaymethod"
	"github.com/IceFoxs/open-gateway/cache/gatewaysystem"
	"github.com/IceFoxs/open-gateway/conf"
	"github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/db"
	"github.com/IceFoxs/open-gateway/registry"
	"github.com/IceFoxs/open-gateway/rpc/http"
	"github.com/IceFoxs/open-gateway/server/consul"
	na "github.com/IceFoxs/open-gateway/server/nacos"
	"github.com/IceFoxs/open-gateway/server/router"
	"github.com/cloudwego/hertz/pkg/app/server"
	re "github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"os"
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
	http.GetHttpClient()
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
	var r re.Registry
	var err error
	if register == constant.REGISTRY_NACOS {
		r, err = na.CreateRegistry(address, username, password)
	} else if register == constant.REGISTRY_CONSUL {
		r, err = consul.CreateRegistry(host, appName, address)
	} else {
		return nil, fmt.Errorf("register %s is not supported", register)
	}
	if err != nil {
		return nil, fmt.Errorf("register occur error,%s", err.Error())
	}
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(dir)
	h := server.Default(server.WithHostPorts(host), server.WithRegistry(r, &re.Info{
		ServiceName: appName,
		Addr:        utils.NewNetAddr("tcp", host),
		Weight:      1,
		Tags:        nil,
	}))
	router.AddRouter(h, dir)
	return h, nil
}
