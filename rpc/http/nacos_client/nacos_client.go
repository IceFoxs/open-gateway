package nacos_client

import (
	"context"
	na "github.com/IceFoxs/open-gateway/nacos"
	c "github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/client/loadbalance"
	"github.com/cloudwego/hertz/pkg/app/middlewares/client/sd"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/registry/nacos/v2"
	"time"
)

type NacosDiscoveryClient struct {
	Client *c.Client
}

func NewDiscoveryClient() (*NacosDiscoveryClient, error) {
	namingClient := na.GetNamingClient()
	client, err := c.NewClient()
	if err != nil {
		panic(err)
	}
	r := nacos.NewNacosResolver(namingClient)
	client.Use(sd.Discovery(r, sd.WithLoadBalanceOptions(loadbalance.NewWeightedBalancer(), loadbalance.Options{
		RefreshInterval: 5 * time.Second,
		ExpireInterval:  15 * time.Second,
	})))
	rc := &NacosDiscoveryClient{
		Client: client,
	}
	return rc, err
}

func (nc *NacosDiscoveryClient) Post(ctx context.Context, appname string, path string, param interface{}) (interface{}, error) {
	req := &protocol.Request{}
	res := &protocol.Response{}
	req.Header.SetMethod(consts.MethodPost)
	req.Header.SetContentTypeBytes([]byte("application/json"))
	req.SetRequestURI("http://" + appname + path)
	req.SetOptions(config.WithSD(true))
	req.SetBody([]byte(param.(string)))
	err := nc.Client.Do(context.Background(), req, res)
	if err != nil {
		return nil, err
	}
	return res.Body(), nil
}
