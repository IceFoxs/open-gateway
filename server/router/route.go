package router

import (
	"context"
	"encoding/base64"
	"github.com/IceFoxs/open-gateway/cache/appmetadata"
	"github.com/IceFoxs/open-gateway/cache/gatewayconfig"
	"github.com/IceFoxs/open-gateway/cache/gatewaymethod"
	"github.com/IceFoxs/open-gateway/common"
	"github.com/IceFoxs/open-gateway/common/regex"
	"github.com/IceFoxs/open-gateway/conf"
	"github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/db/mysql"
	"github.com/IceFoxs/open-gateway/model"
	"github.com/IceFoxs/open-gateway/rpc"
	"github.com/IceFoxs/open-gateway/server/handler"
	"github.com/IceFoxs/open-gateway/sync"
	"github.com/IceFoxs/open-gateway/util/aes"
	rsaUtil "github.com/IceFoxs/open-gateway/util/rsa"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/dubbogo/gost/log/logger"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func AddRouter(h *server.Hertz) {
	var webPath = os.Getenv(constant.WEB_PATH)
	if len(webPath) == 0 {
		if strings.HasPrefix(conf.GetConf().App.WebPath, string(filepath.Separator)) || strings.HasPrefix(conf.GetConf().App.WebPath, ":") {
			webPath = conf.GetConf().App.WebPath
		} else {
			staticPath := conf.GetConf().BaseDir
			if len(staticPath) == 0 {
				staticPath, _ = os.Getwd()
			}
			hlog.Infof("static path is %s", staticPath)
			webPath = staticPath + string(filepath.Separator) + conf.GetConf().App.WebPath
		}
	}
	hlog.Infof("webPath is %s", webPath)
	h.StaticFS("/", &app.FS{Root: webPath, IndexNames: []string{"index.html"}})
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, "ok")
	})
	h.POST("/decryptContent", func(ctx context.Context, c *app.RequestContext) {
		var req common.DecryptContentReq
		err := c.BindAndValidate(&req)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, err)
			return
		}

		cache := gatewayconfig.GetGatewayConfigCache()
		gc, _ := cache.GetCache(req.AppId)
		de, e := base64.StdEncoding.DecodeString(req.EncryptContent)
		if e != nil {
			hlog.Errorf("base64 decode error %s", e)
			c.Abort()
			c.JSON(consts.StatusOK, common.Error(955, "加解密失败"))
			return
		}
		body := aes.AesDecryptECB(de, []byte(gc.AesKey))
		c.Data(consts.StatusOK, "application/json", body)
	})
	h.POST("/selectByAppId", func(ctx context.Context, c *app.RequestContext) {
		var req model.GatewayConfigReq
		err := c.BindAndValidate(&req)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, err)
			return
		}
		if req.PageSize != 0 {
			var res = model.GatewayConfigResponse{}
			g, total, _ := mysql.QueryGatewayConfigByPage(req.AppId, req.PageIndex, req.PageSize)
			res.Models = g
			res.Total = total
			c.JSON(consts.StatusOK, res)
			return
		}
		g, _ := mysql.GetGatewayConfig(req.AppId)
		c.JSON(consts.StatusOK, g)
	})
	h.POST("/addChannelConfig", func(ctx context.Context, c *app.RequestContext) {
		var req model.GatewayConfig
		err := c.BindAndValidate(&req)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, err)
			return
		}
		err = mysql.CreateGatewayConfig([]*model.GatewayConfig{
			{
				AppId:               req.AppId,
				AppName:             req.AppName,
				CallbackUrl:         req.CallbackUrl,
				RsaPrivateKey:       req.RsaPrivateKey,
				RsaPublicKey:        req.RsaPublicKey,
				ClientRsaPrivateKey: req.ClientRsaPrivateKey,
				ClientRsaPublicKey:  req.ClientRsaPublicKey,
				AesKey:              req.AesKey,
				AesType:             req.AesType,
				SignType:            req.SignType,
				IsEnable:            req.IsEnable,
			},
		})
		if err != nil {
			c.JSON(consts.StatusOK, 0)
			return
		}
		c.JSON(consts.StatusOK, 1)
	})

	h.POST("/updateChannelConfig", func(ctx context.Context, c *app.RequestContext) {
		var req model.GatewayConfig
		err := c.BindAndValidate(&req)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, err)
			return
		}
		err = mysql.UpdateGatewayConfig(&model.GatewayConfig{
			AppId:               strings.Trim(strings.Trim(req.AppId, "\n"), " "),
			AppName:             req.AppName,
			CallbackUrl:         req.CallbackUrl,
			RsaPrivateKey:       req.RsaPrivateKey,
			RsaPublicKey:        req.RsaPublicKey,
			ClientRsaPrivateKey: req.ClientRsaPrivateKey,
			ClientRsaPublicKey:  req.ClientRsaPublicKey,
			AesKey:              req.AesKey,
			AesType:             req.AesType,
			SignType:            req.SignType,
			IsEnable:            req.IsEnable,
		})
		if err != nil {
			c.JSON(consts.StatusOK, 0)
			return
		}
		c.JSON(consts.StatusOK, 1)
	})
	h.POST("/addSystemConfig", func(ctx context.Context, c *app.RequestContext) {
		var req model.GatewaySystemConfig
		err := c.BindAndValidate(&req)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, err)
			return
		}
		err = mysql.CreateGatewaySystem([]*model.GatewaySystemConfig{
			{
				SystemId:   strings.Trim(strings.Trim(req.SystemId, "\n"), " "),
				SystemName: req.SystemName,
			},
		})
		if err != nil {
			c.JSON(consts.StatusOK, 0)
			return
		}
		appmetadata.GetAppMetadataCache().RefreshCacheByAppName([]string{req.SystemId})
		a, ok := appmetadata.GetAppMetadataCache().GetAppMetadata(req.SystemId)
		if ok {
			gatewaymethod.GetGatewayMethodCache().RefreshAllCache(a.Methods)
		}
		c.JSON(consts.StatusOK, 1)
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
	h.GET("/channelConfig/refresh", func(ctx context.Context, c *app.RequestContext) {
		sync.GetConfChangeClientHelper().Publish("GATEWAY_CHANNEL", "FPS_GROUP", strings.ReplaceAll(time.Now().Format("20060102150405.000"), ".", "")+"|"+uuid.NewString())
		c.JSON(consts.StatusOK, "ok")
	})
	h.GET("/gateway/refresh", func(ctx context.Context, c *app.RequestContext) {
		sync.GetConfChangeClientHelper().Publish("GATEWAY_SYSTEM", "FPS_GROUP", strings.ReplaceAll(time.Now().Format("20060102150405.000"), ".", "")+"|"+uuid.NewString())
		c.JSON(consts.StatusOK, "ok")
	})
	h.POST("/selectBySysId", func(ctx context.Context, c *app.RequestContext) {
		var req model.GatewaySystemReq
		err := c.BindAndValidate(&req)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, err)
			return
		}
		var g []*model.GatewaySystemConfig
		if req.PageSize != 0 {
			var res = model.GatewaySystemResponse{}
			var total int64
			g, total, _ = mysql.GetGatewaySystemConfigByPage(req.SysId, req.PageIndex, req.PageSize)
			res.Models = g
			res.Total = total
			c.JSON(consts.StatusOK, res)
			return
		}
		g, _ = mysql.GetGatewaySystemConfig(req.SysId)
		c.JSON(consts.StatusOK, g)
	})

	h.POST("/api/json", validFileName, validSign, func(ctx context.Context, c *app.RequestContext) {
		var r, _ = c.Get(common.REQ)
		req := r.(common.RequiredReq)
		var fr, _ = c.Get(common.FILENAME_REQ)
		fileReq := fr.(*regex.FilenameReq)
		hlog.Infof("fileReq:%s", fileReq)
		hlog.Infof("Req:%s", common.ToJSON(req))
		bizContent, ok := c.Get(common.REQ_BODY)
		if !ok {
			bizContent = req.BizContent
		}
		hlog.Infof("filename:%s,bizContent:%s", fileReq.FilenamePre, bizContent)
		rpc.Invoke(context.TODO(), c, req, fileReq, bizContent)
	})
}

func validFileName(ctx context.Context, c *app.RequestContext) {
	// BindAndValidate
	var req common.RequiredReq
	err := c.BindAndValidate(&req)
	if err != nil {
		logger.Errorf("validFileName error: %s", err.Error())
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
	cache := gatewayconfig.GetGatewayConfigCache()
	gc, _ := cache.GetCache(fileReq.AppId)
	if req.SignType != gc.SignType {
		hlog.Errorf("signType not match  %s", req.SignType)
		c.Abort()
		c.JSON(consts.StatusOK, common.Error(957, "签名类型不匹配"))
		return
	}
	if req.SignType == "NONE" {
		return
	}
	publicKey, _ := rsaUtil.Base64PublicKeyToRSA(gc.RsaPublicKey)
	var param = make(map[string]interface{})
	param["bizContent"] = req.BizContent
	param["version"] = req.Version
	param["filename"] = req.Filename
	param["signType"] = req.SignType
	param["encryptType"] = req.EncryptType
	param["timestamp"] = req.Timestamp
	param["sign"] = req.Sign
	sortedParams := rsaUtil.SortParam(param, true)
	err := rsaUtil.RSAVerifyByString(sortedParams, req.Sign, publicKey)
	if err != nil {
		hlog.Errorf("Signature verification failed %s", err)
		c.Abort()
		c.JSON(consts.StatusOK, common.Error(954, "验签失败:"+err.Error()))
		return
	}
	// 验证签名的有效性
	if gc.AesType != req.EncryptType {
		hlog.Errorf("encryptType not match")
		c.Abort()
		c.JSON(consts.StatusOK, common.Error(955, "encryptType not match"))
		return
	}
	if gc.AesType == constant.ENCRYPT_TYPE_AES {
		de, e := base64.StdEncoding.DecodeString(req.BizContent)
		if e != nil {
			hlog.Errorf("base64 decode error %s", e)
			c.Abort()
			c.JSON(consts.StatusOK, common.Error(955, "加解密失败"))
			return
		}
		body := aes.AesDecryptECB(de, []byte(gc.AesKey))
		c.Set(common.REQ_BODY, string(body))
		return
	}
	if gc.AesType == constant.ENCRYPT_TYPE_NONE {
		c.Set(common.REQ_BODY, req.BizContent)
	}

}
