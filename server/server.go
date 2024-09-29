package server

import (
	"errors"
	"github.com/IceFoxs/open-gateway/conf"
	"github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/server/consul"
	na "github.com/IceFoxs/open-gateway/server/nacos"
	zk "github.com/IceFoxs/open-gateway/server/zookeeper"
	"github.com/cloudwego/hertz/pkg/app/server"
	re "github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func CreateServer() (*server.Hertz, error) {
	register := conf.GetConf().App.Register
	if register == "" {
		panic("app register can not empty, please check your config")
	}
	appName := conf.GetConf().App.Name
	host := conf.GetConf().App.Host
	var r re.Registry
	var err error
	if register == constant.REGISTRY_NACOS {
		r, err = na.CreateRegistry()
	} else if register == constant.REGISTRY_CONSUL {
		r, err = consul.CreateRegistry()
	} else if register == constant.REGISTRY_ZOOKEEPER {
		r, err = zk.CreateRegistry()
	} else {
		return nil, errors.New("register[" + register + "]is not supported")
	}
	if err != nil {
		return nil, err
	}
	weight := conf.GetConf().App.Weight
	if weight <= 0 {
		weight = 1
	}
	h := server.Default(server.WithHostPorts(host), server.WithRegistry(r, &re.Info{
		ServiceName: appName,
		Addr:        utils.NewNetAddr("tcp", host),
		Weight:      weight,
		Tags:        nil,
	}))
	return h, nil
}
