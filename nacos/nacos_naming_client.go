package nacos

import (
	"github.com/IceFoxs/open-gateway/conf"
	con "github.com/IceFoxs/open-gateway/constant"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
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
			hlog.Errorf("initNacosNamingClient failed %s", err)
			return
		}
		hlog.Infof("initNacosNamingClient success")
	}
}

func CreateNamingClient(host string, port uint64, username string, password string) (iClient naming_client.INamingClient, err error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(host, port, constant.WithContextPath("/nacos")),
	}
	cc := *constant.NewClientConfig(
		constant.WithUsername(username),
		constant.WithPassword(password),
		constant.WithNamespaceId(""),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir(conf.GetConf().BaseDir+"/logs/nacos/log"),
		constant.WithCacheDir(conf.GetConf().BaseDir+"/logs/nacos/cache"),
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
