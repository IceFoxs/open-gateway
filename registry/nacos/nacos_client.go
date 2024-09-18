package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type NacosRegisterClient struct {
	Client config_client.IConfigClient
}

func NewRegisterClient(host string, port uint64, username string, password string) (*NacosRegisterClient, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(host, port, constant.WithContextPath("/nacos")),
	}
	cc := *constant.NewClientConfig(
		constant.WithUsername(username),
		constant.WithPassword(password),
		constant.WithNamespaceId(""),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)
	iclient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	rc := &NacosRegisterClient{
		Client: iclient,
	}
	return rc, err
}

func (rc *NacosRegisterClient) PublishConfig(key string, group string, value string) error {
	_, err := rc.Client.PublishConfig(vo.ConfigParam{
		DataId:  key,
		Group:   group,
		Content: value,
	})
	return err
}

func (rc *NacosRegisterClient) GetConfig(key string, group string) (string, error) {
	data, err := rc.Client.GetConfig(vo.ConfigParam{
		DataId: key,
		Group:  group,
	})
	return data, err
}

func (rc *NacosRegisterClient) Subscribe(key string, group string) {
	err := rc.Client.ListenConfig(vo.ConfigParam{
		DataId: key,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {

		},
	})
	if err != nil {
		return
	}
}
