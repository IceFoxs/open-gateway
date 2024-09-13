package nacos

import (
	"fmt"
	con "github.com/IceFoxs/open-gateway/constant"
	na "github.com/IceFoxs/open-gateway/nacos"
	"github.com/IceFoxs/open-gateway/server/router"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/registry/nacos/v2"
	"os"
)

func CreateNacosServer() {
	namingClient, err := na.CreateNamingClient("127.0.0.1", 8848, "nacos", "nacos")
	r := nacos.NewNacosRegistry(namingClient)
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(dir)
	addr := "127.0.0.1:8889"
	h := server.Default(server.WithHostPorts(addr), server.WithRegistry(r, &registry.Info{
		ServiceName: con.TEST_SERVEICE,
		Addr:        utils.NewNetAddr("tcp", addr),
		Weight:      10,
		Tags:        nil,
	}))
	router.AddRouter(h, dir)
	h.Spin()
}
