package nacos

import (
	na "github.com/IceFoxs/open-gateway/nacos"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/hertz-contrib/registry/nacos/v2"
	"strconv"
	"strings"
)

func CreateRegistry(address string, username string, password string) (registry.Registry, error) {
	addresses := strings.Split(address, ":")
	host := addresses[0]
	port, _ := strconv.ParseUint(addresses[1], 0, 64)
	namingClient, err := na.CreateNamingClient(host, port, username, password)
	if err != nil {
		return nil, err
	}
	r := nacos.NewNacosRegistry(namingClient)
	return r, nil
}
