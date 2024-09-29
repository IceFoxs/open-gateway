package inits

import (
	"github.com/IceFoxs/open-gateway/conf"
	"github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/sync/config/nacos"
	"github.com/IceFoxs/open-gateway/sync/config/zookeeper"
)

func Init() {
	configType := conf.GetConf().SyncConfig.ConfigType
	if configType == constant.REGISTRY_NACOS {
		nacos.GetConfChangeClient()
	}
	if configType == constant.REGISTRY_ZOOKEEPER {
		zookeeper.GetConfChangeClient()
	}
}
