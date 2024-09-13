package router

import (
	"context"
	ge "github.com/IceFoxs/open-gateway/rpc/generic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"time"
)

func AddRouter(h *server.Hertz, dir string) {
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		re := ge.NewRefConf1("com.hundsun.manager.model.proto.ConfRefreshRpcService", "nacos", "interface", "dubbo", "127.0.0.1:8848", "nacos", "nacos")
		time.Sleep(1 * time.Second)
		c.JSON(consts.StatusOK, ge.ConfRefresh(re))
		//c.JSON(consts.StatusOK, utils.H{"message": "pong"})
	})
	h.StaticFS("/", &app.FS{Root: dir + "/static", IndexNames: []string{"index.html"}})

	//h.Static("/", "/Users/icefox/Documents/go-project/github.com/IceFoxs/open-gateway/static")
	h.POST("/api/json", func(ctx context.Context, c *app.RequestContext) {
		body, _ := c.Body()
		c.Write(body)
		c.SetContentType("application/json; charset=utf-8")
	})
}
