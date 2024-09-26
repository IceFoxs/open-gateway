package nacos

import (
	"github.com/IceFoxs/open-gateway/conf"
	con "github.com/IceFoxs/open-gateway/constant"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/common/logger"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	configClient *config_client.IConfigClient
	once         sync.Once
)

func GetConfigClient() *config_client.IConfigClient {
	once.Do(initConfigClient)
	return configClient
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
	}
}

func CreateConfigClient(host string, port uint64, username string, password string) (c *config_client.IConfigClient, err error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(host, port, constant.WithContextPath("/nacos")),
	}
	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithUsername(username),
		constant.WithPassword(password),
		constant.WithNamespaceId(""),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(false),
		constant.WithLogDir(conf.GetConf().BaseDir+string(filepath.Separator)+"logs"+string(filepath.Separator)+"nacos"+string(filepath.Separator)+"log"),
		constant.WithCacheDir(conf.GetConf().BaseDir+string(filepath.Separator)+"logs"+string(filepath.Separator)+"nacos"+string(filepath.Separator)+"cache"),
		constant.WithLogLevel("debug"),
	)
	_ = logger.InitLogger(logger.BuildLoggerConfig(cc))
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
