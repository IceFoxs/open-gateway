package gatewaysystem

import (
	"github.com/IceFoxs/open-gateway/db/mysql"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	cmap "github.com/orcaman/concurrent-map/v2"
	"sync"
)

var (
	gatewaySystemCache *GatewaySystemCache
	once               sync.Once
)

type GatewaySystem struct {
	SystemId   string `json:"systemId"`
	SystemName string `json:"systemName"`
}

type GatewaySystemCache struct {
	m cmap.ConcurrentMap[string, GatewaySystem]
}

func GetGatewaySystemCache() *GatewaySystemCache {
	once.Do(initCache)
	return gatewaySystemCache
}

func initCache() {
	gatewaySystemCache = &GatewaySystemCache{m: cmap.New[GatewaySystem]()}
	hlog.SystemLogger().Infof("init GatewaySystem cache")
}
func (*GatewaySystemCache) GetAllAppName() []string {
	return gatewaySystemCache.m.Keys()
}
func (*GatewaySystemCache) PutCache(ga GatewaySystem) {
	gatewaySystemCache.m.Set(ga.SystemId, ga)
}

func (*GatewaySystemCache) RefreshCache() {
	gsc, _ := mysql.GetGatewaySystemConfig("")
	for _, config := range gsc {
		gatewaySystemCache.PutCache(GatewaySystem{SystemId: config.SystemId, SystemName: config.SystemName})
	}
}

func (*GatewaySystemCache) GetCache(appId string) (GatewaySystem, bool) {
	return gatewaySystemCache.m.Get(appId)
}

func (*GatewaySystemCache) DeleteCache(appId string) {
	gatewaySystemCache.m.Remove(appId)
}

func (*GatewaySystemCache) ClearCache() {
	gatewaySystemCache.m.Clear()
}

func (*GatewaySystemCache) ToJson() string {
	b, _ := gatewaySystemCache.m.MarshalJSON()
	return string(b)
}
