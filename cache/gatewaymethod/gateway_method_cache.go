package gatewaymethod

import (
	"github.com/IceFoxs/open-gateway/cache/appmetadata"
	"github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/model"
	"github.com/IceFoxs/open-gateway/registry"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/json"
	cmap "github.com/orcaman/concurrent-map/v2"
	"sync"
)

var (
	gatewayMethodCache *GatewayMethodCache
	once               sync.Once
)

type GatewayMethodCache struct {
	m cmap.ConcurrentMap[string, model.GatewayMethodMetadata]
}

func GetGatewayMethodCache() *GatewayMethodCache {
	once.Do(initCache)
	return gatewayMethodCache
}

func initCache() {
	gatewayMethodCache = &GatewayMethodCache{m: cmap.New[model.GatewayMethodMetadata]()}
	hlog.SystemLogger().Infof("init GatewaySystem cache")
}

func (g *GatewayMethodCache) PutCache(filename string, gm model.GatewayMethodMetadata) {
	g.m.Set(filename, gm)
}

func (g *GatewayMethodCache) RefreshCache(filename string) {
	data, err := registry.GetRegisterClient().GetConfig(filename, constant.GATEWAY_META_DATA)
	if err != nil {
		hlog.Errorf("GetConfig %s failed,error is %s", filename, err.Error())
		return
	}
	hlog.Infof("GetConfig[%s] is %s", filename, data)
	var gmm model.GatewayMethodMetadata
	err = json.Unmarshal([]byte(data), &gmm)
	if err != nil {
		hlog.Errorf("GetConfig %s failed,error is %s", filename, err.Error())
		return
	}
	hlog.Infof("GatewayMethodMetadata [%s] is %s", filename, gmm)
	g.PutCache(filename, gmm)
}

func (g *GatewayMethodCache) RefreshAllCache() {
	amc := appmetadata.GetAppMetadataCache()
	methods := amc.GetAllMethods()
	for _, method := range methods {
		g.RefreshCache(method)
	}
}

func (g *GatewayMethodCache) GetCache(gatewayMethodName string) (model.GatewayMethodMetadata, bool) {
	return g.m.Get(gatewayMethodName)
}

func (g *GatewayMethodCache) DeleteCache(appId string) {
	g.m.Remove(appId)
}

func (g *GatewayMethodCache) ClearCache() {
	g.m.Clear()
}

func (*GatewayMethodCache) ToJson() string {
	b, _ := gatewayMethodCache.m.MarshalJSON()
	return string(b)
}
