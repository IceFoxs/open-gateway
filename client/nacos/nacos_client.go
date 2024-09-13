package main

import (
	"context"
	con "github.com/IceFoxs/open-gateway/constant"
	na "github.com/IceFoxs/open-gateway/nacos"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/client/loadbalance"
	"github.com/cloudwego/hertz/pkg/app/middlewares/client/sd"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	nacos "github.com/hertz-contrib/registry/nacos/v2"
	"time"
)

func main() {
	client, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	namingClient, err := na.CreateNamingClient("127.0.0.1", 8848, "nacos", "nacos")
	r := nacos.NewNacosResolver(namingClient)
	client.Use(sd.Discovery(r, sd.WithLoadBalanceOptions(loadbalance.NewWeightedBalancer(), loadbalance.Options{
		RefreshInterval: 5 * time.Second,
		ExpireInterval:  15 * time.Second,
	})))
	for i := 0; i < 10; i++ {
		status, body, err := client.Get(context.Background(), nil, "http://"+con.TEST_SERVEICE+"/ping?", config.WithSD(true))
		if err != nil {
			hlog.Fatal(err)
		}
		hlog.Infof("code=%d,body=%s\n", status, string(body))
	}

}
