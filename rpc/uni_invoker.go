package rpc

import (
	"context"
	"github.com/IceFoxs/open-gateway/cache/gatewaymethod"
	"github.com/IceFoxs/open-gateway/client"
	"github.com/IceFoxs/open-gateway/common"
	"github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/model"
	ge "github.com/IceFoxs/open-gateway/rpc/generic"
	"github.com/IceFoxs/open-gateway/util"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"time"
)

func Invoke(ctx context.Context, c *app.RequestContext, filename string, param interface{}) {
	hlog.Infof("Invoke filename:[%s]", filename)
	gmm, ok := gatewaymethod.GetGatewayMethodCache().GetCache(filename)
	if !ok {
		c.JSON(consts.StatusOK, common.Error(900, "未找到服务调用信息"))
	}
	hlog.Infof("found rpc invoke metadata --------------- %s", gmm)
	if gmm.RpcType == constant.RPC_DUBOO {
		re := ge.NewRefConf1(gmm.InterfaceName, "nacos", "interface", "dubbo", "127.0.0.1:8848", "nacos", "nacos")
		time.Sleep(1 * time.Second)
		toMap, err := util.JsonStringToMap(param.(string))
		if err != nil {
			c.JSON(consts.StatusOK, common.Error(900, err.Error()))
			return
		}
		hlog.Infof("toMap %s", common.ToJSON(toMap))
		util.ConvertHessianMap(toMap)
		data, err := ge.Invoke(re, gmm.MethodName, gmm.ParameterTypeName, util.ConvertHessianMap(toMap))
		c.JSON(consts.StatusOK, common.Succ(0, data, "NONE"))
		return
	}
	if gmm.RpcType == constant.RPC_HTTP {
		data, err := client.GetHttpClient().Post(ctx, gmm.AppName, gmm.Path, param)
		var res model.Response
		err = json.Unmarshal(data.([]byte), &res)
		if err != nil {
			c.JSON(consts.StatusOK, common.Error(900, err.Error()))
			return
		}
		code := res.Code
		msg := res.Message
		if code != 200 {
			c.JSON(consts.StatusOK, common.Error(code, msg))
			return
		}
		c.JSON(consts.StatusOK, common.Succ(0, res.Data, "NONE"))
		return
	}
	c.JSON(consts.StatusOK, common.Succ(0, "", "NONE"))
}
