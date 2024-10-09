package cmd

import (
	"github.com/IceFoxs/open-gateway/cache/appmetadata"
	"github.com/IceFoxs/open-gateway/cache/gatewayconfig"
	"github.com/IceFoxs/open-gateway/cache/gatewaymethod"
	"github.com/IceFoxs/open-gateway/cache/gatewaysystem"
	"github.com/IceFoxs/open-gateway/conf"
	"github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/db"
	"github.com/IceFoxs/open-gateway/logger"
	"github.com/IceFoxs/open-gateway/registry"
	"github.com/IceFoxs/open-gateway/rpc/dubbo"
	"github.com/IceFoxs/open-gateway/rpc/http"
	"github.com/IceFoxs/open-gateway/server"
	"github.com/IceFoxs/open-gateway/sync/config/nacos"
	"github.com/IceFoxs/open-gateway/sync/config/zookeeper"
	"github.com/IceFoxs/open-gateway/sync/inits"
)

import (
	"github.com/IceFoxs/open-gateway/server/router"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Start() {
	Init()
	h, err := server.CreateServer()
	if err != nil {
		hlog.Errorf("create server failed: %s", err.Error())
	}
	router.AddRouter(h)
	h.Spin()
}

func Init() {
	l := logger.NewLogger()
	hlog.SetLogger(l)
	inits.Init()
	dubbo.InitDefaultDubboClient()
	dsn := conf.GetConf().MySQL.DSN
	db.Init(dsn, l.Logger())
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
	configType := conf.GetConf().SyncConfig.ConfigType
	if configType == constant.REGISTRY_NACOS {
		nacos.GetConfChangeClient()
	}
	if configType == constant.REGISTRY_ZOOKEEPER {
		zookeeper.GetConfChangeClient()
	}
}
