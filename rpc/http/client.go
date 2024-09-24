package http

import (
	"context"
	"github.com/IceFoxs/open-gateway/conf"
	con "github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/rpc/http/nacos_client"
	"github.com/dubbogo/gost/log/logger"
	"strconv"
	"strings"
	"sync"
)

var (
	httpClient *HttpClient
	once       sync.Once
)

type HttpClient struct {
	discoveryClient DiscoveryClient
}

type DiscoveryClient interface {
	Post(ctx context.Context, appname string, path string, param interface{}) (interface{}, error)
}

func GetHttpClient() *HttpClient {
	once.Do(initClient)
	return httpClient
}

func (hc *HttpClient) Post(ctx context.Context, appname string, path string, param interface{}) (interface{}, error) {
	return hc.discoveryClient.Post(ctx, appname, path, param)
}

func initClient() {
	address := conf.GetConf().Registry.RegistryAddress[0]
	username := conf.GetConf().Registry.Username
	password := conf.GetConf().Registry.Password
	register := conf.GetConf().Registry.Register
	if register == con.REGISTRY_NACOS {
		addresses := strings.Split(address, ":")
		host := addresses[0]
		port, _ := strconv.ParseUint(addresses[1], 0, 64)
		discoveryClient, _ := nacos_client.NewDiscoveryClient(host, port, username, password)
		httpClient = &HttpClient{discoveryClient: discoveryClient}
		logger.Infof("init nacos  http client success")
	}
}
