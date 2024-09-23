package registry_test

import (
	"github.com/IceFoxs/open-gateway/registry"
	"github.com/IceFoxs/open-gateway/registry/nacos"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"testing"
)

func TestNacosRegisterClient(t *testing.T) {
	var rc registry.RegisterClient
	rc, _ = nacos.NewRegisterClient()
	// create config client
	err := rc.PublishConfig("131", "gatewayMetadata", "13213")
	if err != nil {
		return
	}
	data, _ := rc.GetConfig("131", "gatewayMetadata")
	hlog.Infof("data is %s", data)
}
