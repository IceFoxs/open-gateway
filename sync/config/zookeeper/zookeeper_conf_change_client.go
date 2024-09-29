package zookeeper

import (
	"errors"
	"fmt"
	"github.com/IceFoxs/open-gateway/common"
	"github.com/IceFoxs/open-gateway/conf"
	sy "github.com/IceFoxs/open-gateway/sync"
	"github.com/IceFoxs/open-gateway/sync/config"
	"github.com/IceFoxs/open-gateway/util/zkutils"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-zookeeper/zk"
	cmap "github.com/orcaman/concurrent-map/v2"
	"strings"
	"sync"
	"time"
)

const (
	Separator = "/"
	PATH      = "/"
	Scheme    = "digest" // For auth
)

var (
	zookeeperConfChangeClient config.ConfChangeClient
	once                      sync.Once
)

type ZookeeperConfChangeClient struct {
	conn           *zk.Conn
	authOpen       bool
	user, password string
	listeners      cmap.ConcurrentMap[string, common.Listener]
}

func newZookeeperRegistry(servers []string, sessionTimeout time.Duration) {
	conn, _, err := zk.Connect(servers, sessionTimeout)
	if err != nil {
		return
	}
	zookeeperConfChangeClient = &ZookeeperConfChangeClient{conn: conn, listeners: cmap.New[common.Listener]()}
}
func newZookeeperRegistryWithAuth(servers []string, sessionTimeout time.Duration, user, password string) {
	if user == "" || password == "" {
		return
	}
	conn, _, err := zk.Connect(servers, sessionTimeout)
	if err != nil {
		return
	}
	auth := []byte(fmt.Sprintf("%s:%s", user, password))
	err = conn.AddAuth(Scheme, auth)
	if err != nil {
		return
	}
	zookeeperConfChangeClient = &ZookeeperConfChangeClient{conn: conn, authOpen: true, user: user, password: password, listeners: cmap.New[common.Listener]()}
}

func (zc *ZookeeperConfChangeClient) Publish(confType string, confGroup string, confContent string) {
	zc.createNode(Separator+confGroup+Separator+confType, []byte(confContent), false)
}

func (zc *ZookeeperConfChangeClient) Subscribe(confType string, confGroup string, listener common.Listener) {
	path := Separator + confGroup + Separator + confType
	_, isExist := zc.listeners.Get(path)
	if !isExist {
		zc.listeners.Set(path, listener)
	}
	keepWatcher := zkutils.NewKeepWatcher(zc.conn)
	go keepWatcher.WatchData(path, func(data []byte, err error) {
		if err != nil {
			hlog.Errorf("%s", err.Error())
			return
		}
		l, ok := zc.listeners.Get(path)
		if ok {
			l(confGroup, confType, string(data))
		}
		hlog.Infof("%s", string(data))
	})
}

func (z *ZookeeperConfChangeClient) createNode(path string, content []byte, ephemeral bool) error {
	i := strings.LastIndex(path, Separator)
	if i > 0 {
		err := z.createNode(path[0:i], nil, false)
		if err != nil && !errors.Is(err, zk.ErrNodeExists) {
			return err
		}
	}
	var flag int32
	if ephemeral {
		flag = zk.FlagEphemeral
	}
	var acl []zk.ACL
	if z.authOpen {
		acl = zk.DigestACL(zk.PermAll, z.user, z.password)
	} else {
		acl = zk.WorldACL(zk.PermAll)
	}
	f, _, err := z.conn.Exists(path)
	if f {
		_, err := z.conn.Set(path, content, -1)
		if err != nil {
			hlog.Infof("Set node [%s] error, cause %s", path, err.Error())
			return err
		}
		return nil
	}
	_, err = z.conn.Create(path, content, flag, acl)
	if err != nil {
		hlog.Infof("create node [%s] error, cause %s", path, err.Error())
		return fmt.Errorf("create node [%s] error, cause %w", path, err)
	}
	return nil
}

func (z *ZookeeperConfChangeClient) deleteNode(path string) error {
	err := z.conn.Delete(path, -1)
	if err != nil && err != zk.ErrNoNode {
		return fmt.Errorf("delete node [%s] error, cause %w", path, err)
	}
	return nil
}

func GetConfChangeClient() config.ConfChangeClient {
	once.Do(initZookeeperConfChangeClient)
	return zookeeperConfChangeClient
}
func initZookeeperConfChangeClient() {
	zkConf := conf.GetConf().Zookeeper
	sessionTimeout := zkConf.SessionTimeout
	if sessionTimeout == 0 {
		sessionTimeout = 40
	}
	if len(zkConf.Username) == 0 && len(zkConf.Password) == 0 {
		hlog.Infof("initZookeeperConfChangeClient NewZookeeperRegistry")
		newZookeeperRegistry(zkConf.Address, time.Duration(sessionTimeout)*time.Second)
	} else {
		newZookeeperRegistryWithAuth(zkConf.Address, time.Duration(sessionTimeout)*time.Second, zkConf.Username, zkConf.Password)
		hlog.Infof("initZookeeperConfChangeClient NewZookeeperRegistryWithAuth")
	}
	if zookeeperConfChangeClient.(*ZookeeperConfChangeClient).conn != nil {
		sy.GetConfChangeClientHelper().BuildConfChangeClient(&zookeeperConfChangeClient)
	}
}
