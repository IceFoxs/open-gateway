package main

import (
	ge "github.com/IceFoxs/open-gateway/rpc/generic"
	"time"
)

func main() {
	re := ge.NewRefConf1("com.hundsun.manager.model.proto.ConfRefreshRpcService", "nacos", "interface", "dubbo", "127.0.0.1:8848", "nacos", "nacos")
	time.Sleep(1 * time.Second)
	ge.ConfRefresh(re)
}
