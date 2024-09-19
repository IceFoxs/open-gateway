package rpc

import (
	"context"
	"errors"
	"github.com/IceFoxs/open-gateway/cache/gatewaymethod"
	"github.com/IceFoxs/open-gateway/client"
	"github.com/IceFoxs/open-gateway/common"
	"github.com/IceFoxs/open-gateway/constant"
	ge "github.com/IceFoxs/open-gateway/rpc/generic"
	"github.com/IceFoxs/open-gateway/util"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"time"
)

func Invoke(ctx context.Context, filename string, param interface{}) (interface{}, error) {
	hlog.Infof("Invoke filename:[%s]", filename)
	gmm, ok := gatewaymethod.GetGatewayMethodCache().GetCache(filename)
	if !ok {
		return nil, errors.New("未找到服务调用信息")
	}
	hlog.Infof("found rpc invoke metadata --------------- %s", gmm)
	if gmm.RpcType == constant.RPC_DUBOO {
		re := ge.NewRefConf1(gmm.InterfaceName, "nacos", "interface", "dubbo", "127.0.0.1:8848", "nacos", "nacos")
		time.Sleep(1 * time.Second)
		toMap, err := util.JsonStringToMap(param.(string))
		if err != nil {
			return nil, err
		}
		hlog.Infof("toMap %s", common.ToJSON(toMap))
		util.ConvertHessianMap(toMap)
		data, err := ge.Invoke(re, gmm.MethodName, gmm.ParameterTypeName, util.ConvertHessianMap(toMap))
		return data, err
	}
	if gmm.RpcType == constant.RPC_HTTP {
		data, err := client.GetHttpClient().Post(ctx, gmm.AppName, gmm.Path, param)
		return string(data.([]byte)), err
	}
	return nil, nil
}
