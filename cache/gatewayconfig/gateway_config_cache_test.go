package gatewayconfig_test

import (
	"github.com/IceFoxs/open-gateway/cache/gatewayconfig"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"testing"
)

func TestCache(t *testing.T) {
	g1 := gatewayconfig.GatewayConfig{
		AppId: "112",
	}
	cache1 := gatewayconfig.GetGatewayConfigCache()
	cache1.PutCache(g1)
	g2, _ := cache1.GetCache("112")
	hlog.SystemLogger().Infof("g2  - %v", g2)
	hlog.SystemLogger().Infof("cache1  - %v", cache1.ToJson())
	cache1.ClearCache()
	g2, _ = cache1.GetCache("112")
	hlog.SystemLogger().Infof("g2  - %v", g2)
	hlog.SystemLogger().Infof("cache1  - %v", cache1.ToJson())

}
