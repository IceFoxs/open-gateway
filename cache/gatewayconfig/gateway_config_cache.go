package gatewayconfig

import (
	"github.com/IceFoxs/open-gateway/db/mysql"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	cmap "github.com/orcaman/concurrent-map/v2"
	"sync"
)

var (
	gatewayConfigCache *GatewayConfigCache
	once               sync.Once
)

type GatewayConfig struct {
	AppId         string `json:"appId"`
	AppName       string `json:"appName"`
	AesKey        string `json:"aesKey"`
	AesType       string `json:"aesType"`
	RsaPublicKey  string `json:"rsaPublicKey"`
	RsaPrivateKey string `json:"rsaPrivateKey"`
	SignType      string `json:"signType"`
	IsEnable      int    `json:"isEnable"`
}

type GatewayConfigCache struct {
	m cmap.ConcurrentMap[string, GatewayConfig]
}

func GetGatewayConfigCache() *GatewayConfigCache {
	once.Do(initCache)
	return gatewayConfigCache
}

func initCache() {
	gatewayConfigCache = &GatewayConfigCache{m: cmap.New[GatewayConfig]()}
	hlog.SystemLogger().Infof("init GatewayConfig cache")
}

func (*GatewayConfigCache) PutCache(gatewayConfig GatewayConfig) {
	gatewayConfigCache.m.Set(gatewayConfig.AppId, gatewayConfig)
}

func (g *GatewayConfigCache) RefreshCache() {
	configs, _ := mysql.GetGatewayChannelConfig("")
	for _, config := range configs {
		gConfig := GatewayConfig{
			AppId:         config.AppId,
			AppName:       config.AppName,
			AesKey:        config.AesKey,
			RsaPublicKey:  config.RsaPublicKey,
			RsaPrivateKey: config.RsaPrivateKey,
			AesType:       config.AesType,
			SignType:      config.SignType,
			IsEnable:      config.IsEnable,
		}
		g.PutCache(gConfig)
	}
}

func (*GatewayConfigCache) GetCache(appId string) (GatewayConfig, bool) {
	return gatewayConfigCache.m.Get(appId)
}

func (*GatewayConfigCache) DeleteCache(appId string) {
	gatewayConfigCache.m.Remove(appId)
}

func (*GatewayConfigCache) ClearCache() {
	gatewayConfigCache.m.Clear()
}

func (*GatewayConfigCache) ToJson() string {
	b, _ := gatewayConfigCache.m.MarshalJSON()
	return string(b)
}
