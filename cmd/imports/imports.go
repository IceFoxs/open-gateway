package imports

import (
	"github.com/IceFoxs/open-gateway/cache/appmetadata"
	"github.com/IceFoxs/open-gateway/cache/gatewayconfig"
	"github.com/IceFoxs/open-gateway/cache/gatewaymethod"
	"github.com/IceFoxs/open-gateway/cache/gatewaysystem"
	"github.com/IceFoxs/open-gateway/conf"
	"github.com/IceFoxs/open-gateway/db"
	"github.com/IceFoxs/open-gateway/logger"
	"github.com/IceFoxs/open-gateway/registry"
	"github.com/IceFoxs/open-gateway/rpc/dubbo"
	"github.com/IceFoxs/open-gateway/rpc/http"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func init() {
	hlog.SetLogger(logger.NewLogger())
	dubbo.InitDefaultDubboClient()
	dsn := conf.GetConf().MySQL.DSN
	db.Init(dsn)
	registry.GetRegisterClient()
	http.GetHttpClient()
	c := gatewayconfig.GetGatewayConfigCache()
	c.RefreshCache()
	gsc := gatewaysystem.GetGatewaySystemCache()
	gsc.RefreshCache()
	appNames := gsc.GetAllAppName()
	amc := appmetadata.GetAppMetadataCache()
	amc.RefreshCacheByAppName(appNames)
	gmc := gatewaymethod.GetGatewayMethodCache()
	methods := amc.GetAllMethods()
	gmc.RefreshAllCache(methods)
}
