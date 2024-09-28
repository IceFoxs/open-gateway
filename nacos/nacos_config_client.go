package nacos

import (
	"github.com/IceFoxs/open-gateway/conf"
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
	address := conf.GetConf().Nacos.Address
	username := conf.GetConf().Nacos.Username
	password := conf.GetConf().Nacos.Password
	client, err := CreateConfigClient(address, username, password)
	configClient = client
	if err != nil {
		hlog.Errorf("initNacosConfigClient failed %s", err)
		return
	}
	hlog.Infof("initNacosConfigClient success")
}

func CreateConfigClient(hosts []string, username string, password string) (c *config_client.IConfigClient, err error) {
	var sc []constant.ServerConfig
	for _, address := range hosts {
		addresses := strings.Split(address, ":")
		host := addresses[0]
		port, _ := strconv.ParseUint(addresses[1], 0, 64)
		sc = append(sc, *constant.NewServerConfig(host, port, constant.WithContextPath("/nacos")))
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
