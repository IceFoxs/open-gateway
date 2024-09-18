package nacos

import (
	"fmt"
	na "github.com/IceFoxs/open-gateway/nacos"
	"github.com/IceFoxs/open-gateway/server/router"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/registry/nacos/v2"
	"os"
	"strconv"
	"strings"
)

func CreateNacosServer(serverHost string, appName string, address string, username string, password string) (*server.Hertz, error) {
	addresses := strings.Split(address, ":")
	host := addresses[0]
	port, _ := strconv.ParseUint(addresses[1], 0, 64)
	namingClient, err := na.CreateNamingClient(host, port, username, password)
	r := nacos.NewNacosRegistry(namingClient)
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(dir)
	h := server.Default(server.WithHostPorts(serverHost), server.WithRegistry(r, &registry.Info{
		ServiceName: appName,
		Addr:        utils.NewNetAddr("tcp", serverHost),
		Weight:      1,
		Tags:        nil,
	}))
	router.AddRouter(h, dir)
	return h, nil
}
