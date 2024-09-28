package nacos

import (
	"github.com/IceFoxs/open-gateway/conf"
	con "github.com/IceFoxs/open-gateway/constant"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/common/logger"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	namingClient naming_client.INamingClient
	onceN        sync.Once
)

func GetNamingClient() naming_client.INamingClient {
	onceN.Do(initNamingClient)
	return namingClient
}

func initNamingClient() {
	address := conf.GetConf().Nacos.Address
	username := conf.GetConf().Nacos.Username
	password := conf.GetConf().Nacos.Password
	register := conf.GetConf().Registry.Register
	if register == con.REGISTRY_NACOS {
		n, err := CreateNamingClient(address, username, password)
		namingClient = n
		if err != nil {
			hlog.Errorf("initNacosNamingClient failed %s", err)
			return
		}
		hlog.Infof("initNacosNamingClient success")
	}
}
func CreateNamingClient(hosts []string, username string, password string) (iClient naming_client.INamingClient, err error) {
	var sc []constant.ServerConfig
	for _, address := range hosts {
		addresses := strings.Split(address, ":")
		host := addresses[0]
		port, _ := strconv.ParseUint(addresses[1], 0, 64)
		sc = append(sc, *constant.NewServerConfig(host, port, constant.WithContextPath("/nacos")))
	}
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
	iClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	return
}
