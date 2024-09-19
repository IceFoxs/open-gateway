package registry

import (
	"github.com/IceFoxs/open-gateway/conf"
	"github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/registry/nacos"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"strconv"
	"strings"
	"sync"
)

var (
	rgs  *Registry
	once sync.Once
)

type Registry struct {
	registerClient RegisterClient
}

type Listener func(group, dataId, data string)

type RegisterClient interface {
	PublishConfig(key string, group string, value string) error
	GetConfig(key string, group string) (string, error)
	Subscribe(key string, group string, listener Listener)
}

func (r *Registry) PublishConfig(key string, group string, value string) error {
	return r.registerClient.PublishConfig(key, group, value)
}
func (r *Registry) GetConfig(key string, group string) (string, error) {
	return r.registerClient.GetConfig(key, group)
}
func (r *Registry) Subscribe(key string, group string, listener Listener) {
	r.registerClient.Subscribe(key, group, listener)
}
func GetRegisterClient() *Registry {
	once.Do(initRegisterClient)
	return rgs
}

func initRegisterClient() {
	address := conf.GetConf().Registry.RegistryAddress[0]
	username := conf.GetConf().Registry.Username
	password := conf.GetConf().Registry.Password
	registerType := conf.GetConf().Registry.RegisterType
	if registerType == constant.REGISTRY_NACOS {
		addresses := strings.Split(address, ":")
		host := addresses[0]
		port, _ := strconv.ParseUint(addresses[1], 0, 64)
		var rc, err = nacos.NewRegisterClient(host, port, username, password)
		if err != nil {
			hlog.Errorf("InitRegisterClient failed %s", err)
			return
		}
		rgs = &Registry{
			registerClient: rc,
		}
		hlog.Infof("InitRegisterClient success")
	}
}
