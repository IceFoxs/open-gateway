package dubbo_test

import (
	"fmt"
	"github.com/IceFoxs/open-gateway/common"
	"github.com/IceFoxs/open-gateway/util"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/dubbogo/gost/log/logger"
	"testing"
)

func TestGenericClient(t *testing.T) {
	s := "{\"confType\": \"BANK_TEST\",\"confContent\": \"TEST|20240930\",\"user\":{\"confType\": \"BANK_TEST\",\"confContent\": \"TEST|20240930\"}}"
	var data map[string]interface{}
	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		return
	}
	logger.Infof("map:%s", common.ToJSON(data))
	fmt.Printf("map:%s", common.ToJSON(util.ConvertHessianMap(data)))

	//re := ge.NewRefConf1("com.hundsun.manager.model.proto.ConfRefreshRpcService", "nacos", "interface", "dubbo", "127.0.0.1:8848", "nacos", "nacos")
	//time.Sleep(1 * time.Second)
	//data, err := ge.ConfRefresh(re)
	//if err != nil {
	//	logger.Errorf("confRefresh error: %+v", err.Error())
	//	return
	//}
	//by, err := json.Marshal(data)
	//log.Println("output json:", string(by), err)
}
