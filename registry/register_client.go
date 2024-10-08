package registry

import (
	"github.com/IceFoxs/open-gateway/common"
	"github.com/IceFoxs/open-gateway/conf"
	"github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/registry/nacos"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"sync"
)

var (
	rgs  *Registry
	once sync.Once
)

type Registry struct {
	registerClient RegisterClient
}

type RegisterClient interface {
	PublishConfig(key string, group string, value string) error
	GetConfig(key string, group string) (string, error)
	Subscribe(key string, group string, listener common.Listener)
}

func (r *Registry) PublishConfig(key string, group string, value string) error {
	return r.registerClient.PublishConfig(key, group, value)
}
func (r *Registry) GetConfig(key string, group string) (string, error) {
	return r.registerClient.GetConfig(key, group)
}
func (r *Registry) Subscribe(key string, group string, listener common.Listener) {
	r.registerClient.Subscribe(key, group, listener)
}
func GetRegisterClient() *Registry {
	once.Do(initRegisterClient)
	return rgs
}

func initRegisterClient() {
	register := conf.GetConf().Registry.Register
	if register == constant.REGISTRY_NACOS {
		var rc, err = nacos.NewRegisterClient()
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
