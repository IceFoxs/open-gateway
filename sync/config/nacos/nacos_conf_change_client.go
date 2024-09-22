package nacos

import (
	"github.com/IceFoxs/open-gateway/common"
	na "github.com/IceFoxs/open-gateway/nacos"
	sy "github.com/IceFoxs/open-gateway/sync"
	"github.com/IceFoxs/open-gateway/sync/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"sync"
)

var (
	nacosConfChangeClient config.ConfChangeClient
	once                  sync.Once
)

type NacosConfChangeClient struct {
	ClientConfig *config_client.IConfigClient
}

func (nc *NacosConfChangeClient) Publish(confType string, confGroup string, confContent string) {
	_, err := (*nc.ClientConfig).PublishConfig(vo.ConfigParam{
		DataId:  confType,
		Group:   confGroup,
		Content: confContent,
	})
	if err != nil {
		hlog.Errorf("publish config dataId:[%s] group:[%s],failed :%s", confType, confGroup, err.Error())
	}
}

func (nc *NacosConfChangeClient) Subscribe(confType string, confGroup string, listener common.Listener) {
	err := (*nc.ClientConfig).ListenConfig(vo.ConfigParam{
		DataId: confType,
		Group:  confGroup,
		OnChange: func(namespace, group, dataId, data string) {
			listener(group, dataId, data)
		},
	})
	if err != nil {
		hlog.Errorf("Listen config dataId:[%s] group:[%s],failed : %s", confType, confGroup, err.Error())
	}
}
func GetConfChangeClient() config.ConfChangeClient {
	once.Do(initNacosConfChangeClient)
	return nacosConfChangeClient
}
func initNacosConfChangeClient() {
	nacosConfChangeClient = &NacosConfChangeClient{
		ClientConfig: na.GetConfigClient(),
	}
	if nacosConfChangeClient.(*NacosConfChangeClient).ClientConfig != nil {
		sy.GetConfChangeClientHelper().BuildConfChangeClient(&nacosConfChangeClient)
	}
}
