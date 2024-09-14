package router

import (
	"context"
	"github.com/IceFoxs/open-gateway/common"
	"github.com/IceFoxs/open-gateway/common/regex"
	ge "github.com/IceFoxs/open-gateway/rpc/generic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"time"
)

func AddRouter(h *server.Hertz, dir string) {
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		re := ge.NewRefConf1("com.hundsun.manager.model.proto.ConfRefreshRpcService", "nacos", "interface", "dubbo", "127.0.0.1:8848", "nacos", "nacos")
		time.Sleep(1 * time.Second)
		data, _ := ge.ConfRefresh(re)
		c.JSON(consts.StatusOK, data)
	})
	h.StaticFS("/", &app.FS{Root: dir + "/static", IndexNames: []string{"index.html"}})

	h.POST("/api/json", validFileName, func(ctx context.Context, c *app.RequestContext) {
		var r, _ = c.Get(common.REQ)
		req := r.(common.RequiredReq)
		//re := ge.NewRefConf1("com.hundsun.manager.model.proto.ConfRefreshRpcService", "nacos", "interface", "dubbo", "127.0.0.1:8848", "nacos", "nacos")
		//time.Sleep(1 * time.Second)
		//data, _ := ge.ConfRefresh(re)
		c.JSON(consts.StatusOK, common.Succ(0, req, "NONE"))
	})
}

func validFileName(ctx context.Context, c *app.RequestContext) {
	// BindAndValidate
	var req common.RequiredReq
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.Errorf("validFileName error: %s", err.Error())
		c.Abort()
		c.JSON(consts.StatusOK, common.Error(500, err.Error()))
		return
	}
	filenameReq, err := regex.MatchFileName(req.Filename)
	if err != nil {
		hlog.Errorf("MatchFileName error: %s", err.Error())
		c.Abort()
		c.JSON(consts.StatusOK, common.Error(600, err.Error()))
		return
	}
	c.Set(common.REQ, req)
	c.Set(common.FILENAME_REQ, filenameReq)
	hlog.Infof("validFileName output json:%s", common.ToJSON(req))
}
