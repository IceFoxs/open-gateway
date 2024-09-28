package rpc

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/IceFoxs/open-gateway/cache/gatewayconfig"
	"github.com/IceFoxs/open-gateway/cache/gatewaymethod"
	"github.com/IceFoxs/open-gateway/common"
	"github.com/IceFoxs/open-gateway/common/regex"
	"github.com/IceFoxs/open-gateway/conf"
	"github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/model"
	"github.com/IceFoxs/open-gateway/rpc/dubbo"
	"github.com/IceFoxs/open-gateway/rpc/http"
	"github.com/IceFoxs/open-gateway/util"
	"github.com/IceFoxs/open-gateway/util/aes"
	rsaUtil "github.com/IceFoxs/open-gateway/util/rsa"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"
)

func Invoke(ctx context.Context, c *app.RequestContext, req common.RequiredReq, fileReq *regex.FilenameReq, param interface{}) {
	hlog.Infof("Invoke filename:[%s]", fileReq.FilenamePre)
	gmm, ok := gatewaymethod.GetGatewayMethodCache().GetCache(fileReq.FilenamePre)
	if !ok {
		Error(ctx, c, req, fileReq, 900, "未找到服务调用信息")
		return
	}
	hlog.Infof("found rpc invoke metadata --------------- %s", gmm)
	if gmm.RpcType == constant.RPC_DUBOO {
		toMap, err := util.JsonStringToMap(param.(string))
		if err != nil {
			Error(ctx, c, req, fileReq, 900, err.Error())
			return
		}
		hlog.Infof("toMap %s", common.ToJSON(toMap))
		util.ConvertHessianMap(toMap)
		var wrapResp = conf.GetConf().Dubbo.WrapResp
		if len(wrapResp) == 0 {
			if len(c.GetHeader(constant.WRAP_RESP_HEADER)) == 0 {
				wrapResp = constant.FALSE
			} else {
				wrapResp = string(c.GetHeader(constant.WRAP_RESP_HEADER))
			}
		}
		//re := ge.NewRefConf1(gmm.InterfaceName, "nacos", constant.RPC_INTERFACE_TYPE, "dubbo", "127.0.0.1:8848", "nacos", "nacos")
		//data, err := ge.Invoke(re, gmm.MethodName, gmm.ParameterTypeName, util.ConvertHessianMap(toMap), wrapResp)

		client := dubbo.SingletonDubboClient()
		data, err := client.Invoke(gmm, util.ConvertHessianMap(toMap), wrapResp)
		if data == nil {
			Error(ctx, c, req, fileReq, 900, "服务调用异常")
			return
		}
		if constant.TRUE == wrapResp {
			res := data.(map[string]interface{})
			codeStr := fmt.Sprint(res["code"])
			msg := fmt.Sprint(res["message"])
			if codeStr == "200" {
				Success(ctx, c, req, fileReq, res["data"])
			} else {
				code, _ := strconv.Atoi(codeStr)
				Error(ctx, c, req, fileReq, code, msg)
			}
			return
		}
		if err != nil {
			Error(ctx, c, req, fileReq, 900, err.Error())
			return
		}
		Success(ctx, c, req, fileReq, data)
		return
	}
	if gmm.RpcType == constant.RPC_HTTP {
		data, err := http.GetHttpClient().Post(ctx, gmm.AppName, gmm.Path, param)
		if err != nil {
			Error(ctx, c, req, fileReq, 900, err.Error())
			return
		}
		var res model.Response
		err = json.Unmarshal(data.([]byte), &res)
		if err != nil {
			Error(ctx, c, req, fileReq, 900, err.Error())
			return
		}
		code := res.Code
		msg := res.Message
		if code != 200 {
			Error(ctx, c, req, fileReq, code, msg)
			return
		}
		Success(ctx, c, req, fileReq, res.Data)
		return
	}
	c.JSON(consts.StatusOK, common.Succ(0, "请求成功", "{}", ""))
}
func Error(ctx context.Context, c *app.RequestContext, req common.RequiredReq, fileReq *regex.FilenameReq, statusCode int, msg string) {
	cache := gatewayconfig.GetGatewayConfigCache()
	gConfig, ok := cache.GetCache(fileReq.AppId)
	if ok {
		var param = make(map[string]interface{})
		param["statusCode"] = statusCode
		param["errorMsg"] = msg
		sign := Sign(param, req, gConfig.RsaPrivateKey)
		c.JSON(consts.StatusOK, common.ErrorWithSign(statusCode, msg, sign))
		return
	}
}
func Success(ctx context.Context, c *app.RequestContext, req common.RequiredReq, fileReq *regex.FilenameReq, data interface{}) {
	cache := gatewayconfig.GetGatewayConfigCache()
	gConfig, ok := cache.GetCache(fileReq.AppId)
	msg := "请求成功"
	statusCode := 0
	if ok {
		if req.EncryptType == constant.ENCRYPT_TYPE_AES {
			aesKey := gConfig.AesKey
			content, err := json.Marshal(data)
			body := aes.AesEncryptECB(content, []byte(aesKey))
			if err != nil {
				var param = make(map[string]interface{})
				param["statusCode"] = statusCode
				param["errorMsg"] = msg
				statusCode = 500
				msg = err.Error()
				sign := Sign(param, req, gConfig.RsaPrivateKey)
				c.JSON(consts.StatusOK, common.ErrorWithSign(statusCode, msg, sign))
				return
			}
			de := base64.StdEncoding.EncodeToString(body)
			var param = make(map[string]interface{})
			param["statusCode"] = statusCode
			param["errorMsg"] = msg
			param["bizContent"] = de
			sign := Sign(param, req, gConfig.RsaPrivateKey)
			c.JSON(consts.StatusOK, common.SuccContent(0, msg, de, sign))
			return
		}
		if req.EncryptType == constant.ENCRYPT_TYPE_NONE {
			var param = make(map[string]interface{})
			param["statusCode"] = statusCode
			param["errorMsg"] = msg
			param["bizContent"] = data
			sign := Sign(param, req, gConfig.RsaPrivateKey)
			c.JSON(consts.StatusOK, common.Succ(statusCode, msg, data, sign))
			return
		}
	}
}

func Sign(param map[string]interface{}, req common.RequiredReq, privateKey string) string {
	if req.SignType == constant.SIGN_TYPE_NONE {
		return ""
	}
	signStr := rsaUtil.SortParam(param, false)
	pk3, _ := rsaUtil.Base64ToPrivateKeyByPkcs8(privateKey)
	data, _ := rsaUtil.RSASignByString(signStr, pk3)
	return data
}
