package zookeeper

import (
	"github.com/IceFoxs/open-gateway/conf"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/hertz-contrib/registry/zookeeper"
	"time"
)

func CreateRegistry() (registry.Registry, error) {
	zk := conf.GetConf().Zookeeper
	if len(zk.Username) == 0 || len(zk.Password) == 0 {
		return zookeeper.NewZookeeperRegistry(zk.Address, time.Duration(zk.SessionTimeout)*time.Second)
	}
	r, err := zookeeper.NewZookeeperRegistryWithAuth(zk.Address, time.Duration(zk.SessionTimeout)*time.Second, zk.Username, zk.Password)
	return r, err
}
