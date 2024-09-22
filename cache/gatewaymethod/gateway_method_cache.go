package gatewaymethod

import (
	"github.com/IceFoxs/open-gateway/common"
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
	if len(data) == 0 {
		hlog.Errorf("Config[%s][%s] changed is empty", filename, constant.GATEWAY_META_DATA)
		return
	}
	var gmm model.GatewayMethodMetadata
	err = json.Unmarshal([]byte(data), &gmm)
	if err != nil {
		hlog.Errorf("GetConfig %s failed,error is %s", filename, err.Error())
		return
	}
	hlog.Infof("GatewayMethodMetadata [%s] is %s", filename, common.ToJSON(gmm))
	g.PutCache(filename, gmm)
}

func (g *GatewayMethodCache) AddListen(method string) {
	registry.GetRegisterClient().Subscribe(method, constant.GATEWAY_META_DATA, gatewayMethodCache.Listen)
}
func (g *GatewayMethodCache) Listen(group, dataId, data string) {
	hlog.Infof("Config Refresh  group:[%s],dataId:[%s],data:[%s]", group, dataId, data)
	g.RefreshCache(dataId)
}

func (g *GatewayMethodCache) RefreshAllCache(methods []string) {
	for _, method := range methods {
		g.AddListen(method)
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
