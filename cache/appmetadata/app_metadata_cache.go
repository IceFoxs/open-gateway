package appmetadata

import (
	"github.com/IceFoxs/open-gateway/cache/gatewaymethod"
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
	methods     cmap.ConcurrentMap[string, string]
}

func GetAppMetadataCache() *AppMetadataCache {
	once.Do(initCache)
	return appMetadataCache
}

func initCache() {
	appMetadataCache = &AppMetadataCache{
		appMetadata: cmap.New[model.AppMetadata](),
		methods:     cmap.New[string](),
	}
	hlog.Infof("init appMetadataCache cache")
}

func (a *AppMetadataCache) GetAllMethods() []string {
	return a.methods.Keys()
}

func (a *AppMetadataCache) GetAppMetadata(appId string) (model.AppMetadata, bool) {
	return a.appMetadata.Get(appId)
}
func (a *AppMetadataCache) GetAppName(methodName string) string {
	appname, _ := a.methods.Get(methodName)
	return appname
}
func (a *AppMetadataCache) PutCache(appMetadata model.AppMetadata) {
	a.appMetadata.Set(appMetadata.AppName, appMetadata)
	if len(appMetadata.Methods) > 0 {
		for _, method := range appMetadata.Methods {
			a.methods.Set(method, appMetadata.AppName)
		}
	}
}
func (g *AppMetadataCache) AddListen(appname string) {
	for _, k := range g.appMetadata.Keys() {
		if appname == k {
			registry.GetRegisterClient().Subscribe(k, constant.APP_METADATA,
				appMetadataCache.Listen)
			registry.GetRegisterClient().Subscribe(k, constant.HTTP_APP_METADATA,
				appMetadataCache.Listen)
		}
	}
}

func (g *AppMetadataCache) Listen(group, dataId, data string) {
	hlog.Infof("Config Refresh  group:[%s],dataId:[%s],data:[%s]", group, dataId, data)
	g.RefreshCacheByAppName([]string{dataId})
	gatewaymethod.GetGatewayMethodCache().RefreshAllCache(g.GetAllMethods())
}

func (a *AppMetadataCache) RefreshCache(appName string, group string) {
	data, err := registry.GetRegisterClient().GetConfig(appName, group)
	if err != nil {
		hlog.Errorf("AppMetadata GetConfig %s failed,error is %s", appName, err.Error())
		return
	}
	hlog.Infof("AppMetadata GetConfig[%s][%s] is %s", appName, group, data)
	if len(data) == 0 {
		amm := model.AppMetadata{
			AppName: appName,
			Methods: []string{},
		}
		hlog.Errorf("AppMetadata GetConfig[%s][%s] is empty", appName, group)
		appMetadataCache.PutCache(amm)
		return
	}
	var amm model.AppMetadata
	err = json.Unmarshal([]byte(data), &amm)
	if err != nil {
		amm = model.AppMetadata{
			AppName: appName,
			Methods: []string{},
		}
		hlog.Errorf("AppMetadata GetConfig %s failed,error is %s", appName, err.Error())
		appMetadataCache.PutCache(amm)
	} else {
		hlog.Infof("AppMetadataCache [%s] is %s", appName, amm)
		appMetadataCache.PutCache(amm)
	}
}
func (a *AppMetadataCache) RefreshAllCache(appNames []string) {
	for _, name := range appNames {
		a.RefreshCache(name, constant.APP_METADATA)
		hlog.Infof("AppMetadata RefreshAllCache APP_METADATA [%s]", name)
		a.RefreshCache(name, constant.HTTP_APP_METADATA)
		hlog.Infof("AppMetadata RefreshAllCache HTTP_APP_METADATA [%s]", name)
	}
}

func (a *AppMetadataCache) RefreshCacheByAppName(appNames []string) {
	for _, name := range appNames {
		a.RefreshCache(name, constant.APP_METADATA)
		hlog.Infof("AppMetadata RefreshAllCache APP_METADATA [%s]", name)
		a.RefreshCache(name, constant.HTTP_APP_METADATA)
		hlog.Infof("AppMetadata RefreshAllCache HTTP_APP_METADATA [%s]", name)
	}
	for _, name := range appNames {
		hlog.Infof("AppMetadata addListen [%s]", name)
		a.AddListen(name)
	}
}
