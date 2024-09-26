package http

import (
	"context"
	"github.com/IceFoxs/open-gateway/conf"
	con "github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/rpc/http/nacos_client"
	"github.com/cloudwego/hertz/pkg/common/hlog"
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
	register := conf.GetConf().Registry.Register
	if register == con.REGISTRY_NACOS {
		discoveryClient, _ := nacos_client.NewDiscoveryClient()
		httpClient = &HttpClient{discoveryClient: discoveryClient}
		hlog.Infof("init nacos  http client success")
	}
}
