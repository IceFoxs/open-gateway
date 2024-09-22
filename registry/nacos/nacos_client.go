package nacos

import (
	"github.com/IceFoxs/open-gateway/common"
	"github.com/IceFoxs/open-gateway/nacos"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type NacosRegisterClient struct {
	Client config_client.IConfigClient
}

func NewRegisterClient() (*NacosRegisterClient, error) {
	client := nacos.GetConfigClient()

	rc := &NacosRegisterClient{
		Client: client,
	}
	return rc, nil
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

func (rc *NacosRegisterClient) Subscribe(key string, group string, f common.Listener) {
	err := rc.Client.ListenConfig(vo.ConfigParam{
		DataId: key,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			f(group, dataId, data)
		},
	})
	if err != nil {
		hlog.Errorf("nacos regrister Subscribe key[%s] group[%s] failed", key, group,
			err.Error())
		return
	}
}
