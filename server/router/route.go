package router

import (
	"context"
	"github.com/IceFoxs/open-gateway/cache/gatewayconfig"
	"github.com/IceFoxs/open-gateway/common"
	"github.com/IceFoxs/open-gateway/common/regex"
	"github.com/IceFoxs/open-gateway/db/mysql"
	"github.com/IceFoxs/open-gateway/model"
	"github.com/IceFoxs/open-gateway/rpc"
	ge "github.com/IceFoxs/open-gateway/rpc/generic"
	"github.com/IceFoxs/open-gateway/server/handler"
	rsaUtil "github.com/IceFoxs/open-gateway/util/rsa"
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"time"
)

func AddRouter(h *server.Hertz, dir string) {
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, "ok")
	})

	h.POST("/selectByAppId", func(ctx context.Context, c *app.RequestContext) {
		var req common.GatewayConfigReq
		err := c.BindAndValidate(&req)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, err)
			return
		}
		g, _ := mysql.GetGatewayChannelConfig(req.AppId)
		c.JSON(consts.StatusOK, g)
	})
	h.POST("/selectAppMethods", func(ctx context.Context, c *app.RequestContext) {
		var req model.GatewayMethodRequest
		err := c.BindAndValidate(&req)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, err)
			return
		}
		g, _ := handler.QueryGatewayMethodInfo(req)
		c.JSON(consts.StatusOK, g)
	})
	h.POST("/selectBySysId", func(ctx context.Context, c *app.RequestContext) {
		var req common.GatewaySystemReq
		err := c.BindAndValidate(&req)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, err)
			return
		}
		g, _ := mysql.GetGatewaySystemConfig(req.SysId)
		c.JSON(consts.StatusOK, g)
	})
	h.GET("/generic", func(ctx context.Context, c *app.RequestContext) {
		re := ge.NewRefConf1("com.hundsun.manager.model.proto.ConfRefreshRpcService", "nacos", "interface", "dubbo", "127.0.0.1:8848", "nacos", "nacos")
		time.Sleep(1 * time.Second)
		var m = make(map[string]hessian.Object)
		m["confType"] = "BANK_TEST"
		m["confContent"] = "TEST|20240930"
		data, _ := ge.Invoke(re, "confRefresh", "com.hundsun.manager.model.req.ConfRefreshRequest", m)

		c.JSON(consts.StatusOK, data)
	})
	h.StaticFS("/", &app.FS{Root: dir + "/static", IndexNames: []string{"index.html"}})

	h.POST("/api/json", validFileName, validSign, func(ctx context.Context, c *app.RequestContext) {
		var r, _ = c.Get(common.REQ)
		req := r.(common.RequiredReq)
		var fr, _ = c.Get(common.FILENAME_REQ)
		fileReq := fr.(*regex.FilenameReq)
		hlog.Infof("fileReq %s", fileReq)
		rpc.Invoke(context.TODO(), c, fileReq.FilenamePre, req.BizContent)
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
	cache := gatewayconfig.GetGatewayConfigCache()
	gConfig, ok := cache.GetCache(filenameReq.AppId)
	if !ok || gConfig.IsEnable == 0 {
		c.Abort()
		c.JSON(consts.StatusOK, common.Error(301, "账号不存在"))
		return
	}
	hlog.Infof("validFileName output json:%s", common.ToJSON(req))
}

func validSign(ctx context.Context, c *app.RequestContext) {
	var r, _ = c.Get(common.REQ)
	req := r.(common.RequiredReq)
	var fr, _ = c.Get(common.FILENAME_REQ)
	fileReq := fr.(*regex.FilenameReq)
	if req.SignType == "NONE" {
		return
	}
	cache := gatewayconfig.GetGatewayConfigCache()
	gc, _ := cache.GetCache(fileReq.AppId)
	publicKey, _ := rsaUtil.Base64PublicKeyToRSA(gc.RsaPublicKey)
	var param = make(map[string]interface{})
	param["bizContent"] = req.BizContent
	param["version"] = req.Version
	param["filename"] = req.Filename
	param["signType"] = req.SignType
	param["encryptType"] = req.EncryptType
	param["timestamp"] = req.Timestamp
	param["sign"] = req.Sign
	sortedParams := rsaUtil.SortParam(param)
	err := rsaUtil.RSAVerifyByString(sortedParams, req.Sign, publicKey)
	if err != nil {
		hlog.Errorf("Signature verification failed %s", err)
		c.Abort()
		c.JSON(consts.StatusOK, common.Error(954, "验签失败:"+err.Error()))
	}
	// 验证签名的有效性
}
