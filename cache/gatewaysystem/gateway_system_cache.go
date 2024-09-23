package gatewaysystem

import (
	"github.com/IceFoxs/open-gateway/cache/appmetadata"
	"github.com/IceFoxs/open-gateway/cache/gatewaymethod"
	"github.com/IceFoxs/open-gateway/db/mysql"
	sy "github.com/IceFoxs/open-gateway/sync"
	"github.com/dubbogo/gost/log/logger"
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
	sy.GetConfChangeClientHelper().Subscribe("GATEWAY_SYSTEM", "FPS_GROUP", gatewaySystemCache.Listen)
	logger.Infof("init GatewaySystem cache")
}
func (*GatewaySystemCache) GetAllAppName() []string {
	return gatewaySystemCache.m.Keys()
}
func (*GatewaySystemCache) PutCache(ga GatewaySystem) {
	gatewaySystemCache.m.Set(ga.SystemId, ga)
}

func (g *GatewaySystemCache) RefreshCache() {
	gsc, _ := mysql.GetGatewaySystemConfig("")
	for _, config := range gsc {
		g.PutCache(GatewaySystem{SystemId: config.SystemId, SystemName: config.SystemName})
	}
}
func (g *GatewaySystemCache) Listen(group, dataId, data string) {
	logger.Infof("Config Refresh  group:[%s],dataId:[%s],data:[%s]", group, dataId, data)
	g.RefreshCache()
	amc := appmetadata.GetAppMetadataCache()
	amc.RefreshCacheByAppName(g.GetAllAppName())
	gatewaymethod.GetGatewayMethodCache().RefreshAllCache(amc.GetAllMethods())
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
