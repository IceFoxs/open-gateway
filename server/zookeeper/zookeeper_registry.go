package zookeeper

import (
	"github.com/IceFoxs/open-gateway/conf"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/registry/zookeeper"
	"time"
)

func CreateRegistry() (registry.Registry, error) {
	zk := conf.GetConf().Zookeeper
	sessionTimeout := zk.SessionTimeout
	if sessionTimeout == 0 {
		sessionTimeout = 40
	}
	if len(zk.Username) == 0 && len(zk.Password) == 0 {
		hlog.Infof("NewZookeeperRegistry")
		return zookeeper.NewZookeeperRegistry(zk.Address, time.Duration(sessionTimeout)*time.Second)
	} else {
		r, err := zookeeper.NewZookeeperRegistryWithAuth(zk.Address, time.Duration(sessionTimeout)*time.Second, zk.Username, zk.Password)
		hlog.Infof("NewZookeeperRegistryWithAuth")
		return r, err
	}
}
