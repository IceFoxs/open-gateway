package nacos

import (
	"github.com/IceFoxs/open-gateway/conf"
	con "github.com/IceFoxs/open-gateway/constant"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	configClient *config_client.IConfigClient
	namingClient naming_client.INamingClient
	once         sync.Once
)

func GetConfigClient() *config_client.IConfigClient {
	once.Do(initConfigClient)
	return configClient
}

func GetNamingClient() naming_client.INamingClient {
	once.Do(initConfigClient)
	return namingClient
}

func initConfigClient() {
	address := conf.GetConf().Registry.RegistryAddress[0]
	username := conf.GetConf().Registry.Username
	password := conf.GetConf().Registry.Password
	register := conf.GetConf().Registry.Register
	if register == con.REGISTRY_NACOS {
		addresses := strings.Split(address, ":")
		host := addresses[0]
		port, _ := strconv.ParseUint(addresses[1], 0, 64)
		client, err := CreateConfigClient(host, port, username, password)
		configClient = client
		if err != nil {
			hlog.Errorf("initNacosConfigClient failed %s", err)
			return
		}
		hlog.Infof("initNacosConfigClient success")
		n, err := CreateNamingClient(host, port, username, password)
		namingClient = n
		if err != nil {
			hlog.Errorf("initNacosConfigClient failed %s", err)
			return
		}
		hlog.Infof("initNacosConfigClient success")
	}
}
func initNamingClient() {
	address := conf.GetConf().Registry.RegistryAddress[0]
	username := conf.GetConf().Registry.Username
	password := conf.GetConf().Registry.Password
	register := conf.GetConf().Registry.Register
	if register == con.REGISTRY_NACOS {
		addresses := strings.Split(address, ":")
		host := addresses[0]
		port, _ := strconv.ParseUint(addresses[1], 0, 64)
		n, err := CreateNamingClient(host, port, username, password)
		namingClient = n
		if err != nil {
			hlog.Errorf("initNacosConfigClient failed %s", err)
			return
		}
		hlog.Infof("initNacosConfigClient success")
	}
}
func CreateConfigClient(host string, port uint64, username string, password string) (c *config_client.IConfigClient, err error) {
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
		constant.WithLogDir(path+"/logs/nacos/log"),
		constant.WithCacheDir(path+"/logs/nacos/cache"),
		constant.WithLogLevel("debug"),
	)
	// create config client
	iClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	c = &iClient
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
		constant.WithLogDir(path+"/logs/nacos/log"),
		constant.WithCacheDir(path+"/logs/nacos/cache"),
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
