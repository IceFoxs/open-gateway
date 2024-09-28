package nacos

import (
	na "github.com/IceFoxs/open-gateway/nacos"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/hertz-contrib/registry/nacos/v2"
)

func CreateRegistry() (registry.Registry, error) {
	namingClient := na.GetNamingClient()
	r := nacos.NewNacosRegistry(namingClient)
	return r, nil
}
