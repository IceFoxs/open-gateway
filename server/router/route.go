package router

import (
	"context"
	"github.com/IceFoxs/open-gateway/common"
	ge "github.com/IceFoxs/open-gateway/rpc/generic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"time"
)

func AddRouter(h *server.Hertz, dir string) {
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		re := ge.NewRefConf1("com.hundsun.manager.model.proto.ConfRefreshRpcService", "nacos", "interface", "dubbo", "127.0.0.1:8848", "nacos", "nacos")
		time.Sleep(1 * time.Second)
		data, _ := ge.ConfRefresh(re)
		c.JSON(consts.StatusOK, data)
		//c.JSON(consts.StatusOK, utils.H{"message": "pong"})
	})
	h.StaticFS("/", &app.FS{Root: dir + "/static", IndexNames: []string{"index.html"}})

	h.POST("/api/json", func(ctx context.Context, c *app.RequestContext) {
		// BindAndValidate
		var req common.RequiredReq
		err := c.BindAndValidate(&req)
		if err != nil {
			c.JSON(consts.StatusOK, utils.H{"message": err.Error()})
			return
		}
		c.JSON(consts.StatusOK, req)
	})
}
