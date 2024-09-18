package gatewaysystem

import (
	"encoding/json"
	"github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/registry"
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

func (*GatewaySystemCache) PutCache(ga GatewaySystem) {
	gatewaySystemCache.m.Set(ga.SystemId, ga)
}

func (g *GatewaySystemCache) RefreshCache(filename string) {
	data, err := registry.GetRegisterClient().GetConfig(filename, constant.GATEWAY_META_DATA)

	if err != nil {
		hlog.Errorf("GetConfig %s failed,error is %s", filename, err.Error())
		return
	}
	hlog.Infof("GetConfig[%s] is %s", filename, data)
	var gConfig GatewaySystem
	err = json.Unmarshal([]byte(data), &gConfig)
	if err != nil {
		hlog.Errorf("GetConfig %s failed,error is %s", filename, err.Error())
		return
	}
	hlog.Infof("gConfig[%s] is %s", filename, gConfig)
	g.PutCache(gConfig)
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
