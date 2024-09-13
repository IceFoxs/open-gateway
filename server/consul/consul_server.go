package consul

import (
	"fmt"
	con "github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/server/router"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/utils"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hertz-contrib/registry/consul"
	"log"
	"os"
)

func CreateConsulServer() {
	config := consulapi.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	consulClient, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal(err)
		return
	}
	addr := "127.0.0.1:8888"
	check := &consulapi.AgentServiceCheck{
		CheckID:  "foo-ttl",
		Interval: "10s",
		HTTP:     "http://" + addr + "/health",
	}
	r := consul.NewConsulRegister(consulClient, consul.WithCheck(check))
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(dir)
	h := server.Default(server.WithHostPorts(addr), server.WithRegistry(r, &registry.Info{
		ServiceName: con.TEST_SERVEICE,
		Addr:        utils.NewNetAddr("tcp", addr),
		Weight:      10,
		Tags:        nil,
	}))
	router.AddRouter(h, dir)
	h.Spin()
}
