package consul

import (
	"context"
	con "github.com/IceFoxs/open-gateway/constant"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/middlewares/client/sd"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/dubbogo/gost/log/logger"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hertz-contrib/registry/consul"
	"log"
)

func main1() {
	client, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = "127.0.0.1:8500"
	consulClient, err := consulapi.NewClient(consulConfig)
	if err != nil {
		log.Fatal(err)
		return
	}
	r := consul.NewConsulResolver(consulClient)
	client.Use(sd.Discovery(r))
	for i := 0; i < 10; i++ {
		status, body, err := client.Get(context.Background(), nil, "http://"+con.TEST_SERVEICE+"/ping?", config.WithSD(true))
		if err != nil {
			hlog.Fatal(err)
		}
		logger.Infof("code=%d,body=%s\n", status, string(body))
	}

}
