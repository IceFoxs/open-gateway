package dubbo_test

import (
	"encoding/json"
	ge "github.com/IceFoxs/open-gateway/rpc/generic"
	"github.com/bytedance/gopkg/util/logger"
	"log"
	"testing"
	"time"
)

func TestGenericClient(t *testing.T) {
	re := ge.NewRefConf1("com.hundsun.manager.model.proto.ConfRefreshRpcService", "nacos", "interface", "dubbo", "127.0.0.1:8848", "nacos", "nacos")
	time.Sleep(1 * time.Second)
	data, err := ge.ConfRefresh(re)
	if err != nil {
		logger.Errorf("confRefresh error: %+v", err.Error())
		return
	}
	by, err := json.Marshal(data)
	log.Println("output json:", string(by), err)
}
