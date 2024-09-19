package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"os"
)

func CreateConfigClient(host string, port uint64, username string, password string) (iClient config_client.IConfigClient, err error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(host, port, constant.WithContextPath("/nacos")),
	}
	path, _ := os.Getwd()
	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithUsername(username),
		constant.WithPassword(password),
		constant.WithNamespaceId(""),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir(path+"/tmp/nacos/log"),
		constant.WithCacheDir(path+"/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)

	// create config client
	iClient, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	return
}

func CreateNamingClient(host string, port uint64, username string, password string) (iClient naming_client.INamingClient, err error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(host, port, constant.WithContextPath("/nacos")),
	}
	path, _ := os.Getwd()
	cc := *constant.NewClientConfig(
		constant.WithUsername(username),
		constant.WithPassword(password),
		constant.WithNamespaceId(""),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir(path+"/tmp/nacos/log"),
		constant.WithCacheDir(path+"/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)

	iClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	return
}
