package nacos_test

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"testing"
	"time"
)

func TestListen(t *testing.T) {
	scs := []constant.ServerConfig{*constant.NewServerConfig("127.0.0.1", 8848)}
	cc := constant.ClientConfig{
		NamespaceId:         "", //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: false,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
		Username:            "nacos",
		Password:            "nacos",
	}
	client, _ := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: scs,
		},
	)
	//client, _ := nacos.NewNacosConfigClient("nacos", true, scs, cc)
	client.ListenConfig(vo.ConfigParam{
		DataId: "nacos-config" + "11",
		Group:  "dubbo",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
		},
	})

	time.Sleep(100000000 * time.Millisecond)
}
