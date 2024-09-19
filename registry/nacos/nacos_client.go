package nacos

import (
	"github.com/IceFoxs/open-gateway/nacos"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type NacosRegisterClient struct {
	Client config_client.IConfigClient
}

func NewRegisterClient(host string, port uint64, username string, password string) (*NacosRegisterClient, error) {
	client, err := nacos.CreateConfigClient(host, port, username, password)
	if err != nil {
		return nil, err
	}
	rc := &NacosRegisterClient{
		Client: client,
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
