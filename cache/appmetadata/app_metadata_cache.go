package appmetadata

import (
	"github.com/IceFoxs/open-gateway/cache/gatewaysystem"
	"github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/model"
	"github.com/IceFoxs/open-gateway/registry"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/json"
	cmap "github.com/orcaman/concurrent-map/v2"
	"sync"
)

var (
	appMetadataCache *AppMetadataCache
	once             sync.Once
)

type AppMetadataCache struct {
	appMetadata cmap.ConcurrentMap[string, model.AppMetadata]
	methods     cmap.ConcurrentMap[string, bool]
}

func GetAppMetadataCache() *AppMetadataCache {
	once.Do(initCache)
	return appMetadataCache
}

func initCache() {
	appMetadataCache = &AppMetadataCache{
		appMetadata: cmap.New[model.AppMetadata](),
		methods:     cmap.New[bool](),
	}
	hlog.SystemLogger().Infof("init appMetadataCache cache")
}

func (a *AppMetadataCache) GetAllMethods() []string {
	return a.methods.Keys()
}
func (a *AppMetadataCache) PutCache(appMetadata model.AppMetadata) {
	a.appMetadata.Set(appMetadata.AppName, appMetadata)
	if len(appMetadata.Methods) > 0 {
		for _, method := range appMetadata.Methods {
			a.methods.Set(method, true)
		}
	}
}

func (a *AppMetadataCache) RefreshCache(appName string, group string) {
	data, err := registry.GetRegisterClient().GetConfig(appName, group)
	if err != nil {
		hlog.Errorf("AppMetadata GetConfig %s failed,error is %s", appName, err.Error())
		return
	}
	hlog.Infof("AppMetadata GetConfig[%s] is %s", appName, data)
	var amm model.AppMetadata
	err = json.Unmarshal([]byte(data), &amm)
	if err != nil {
		hlog.Errorf("AppMetadata GetConfig %s failed,error is %s", appName, err.Error())
		return
	}
	hlog.Infof("AppMetadataCache [%s] is %s", appName, amm)
	appMetadataCache.PutCache(amm)
}

func (a *AppMetadataCache) RefreshAllCache() {
	gsc := gatewaysystem.GetGatewaySystemCache()
	appNames := gsc.GetAllAppName()
	for _, name := range appNames {
		a.RefreshCache(name, constant.APP_METADATA)
		hlog.Infof("AppMetadata RefreshAllCache APP_METADATA [%s]", name)
		a.RefreshCache(name, constant.HTTP_APP_METADATA)
		hlog.Infof("AppMetadata RefreshAllCache HTTP_APP_METADATA [%s]", name)

	}
}
