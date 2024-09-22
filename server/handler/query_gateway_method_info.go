package handler

import (
	"github.com/IceFoxs/open-gateway/cache/appmetadata"
	"github.com/IceFoxs/open-gateway/cache/gatewaymethod"
	"github.com/IceFoxs/open-gateway/model"
	"github.com/cloudwego/hertz/pkg/common/json"
	"sort"
	"strings"
)

func QueryGatewayMethodInfo(req model.GatewayMethodRequest) (model.GatewayMethodResponse, error) {
	var res = model.GatewayMethodResponse{}
	methods := appmetadata.GetAppMetadataCache().GetAllMethods()
	if len(req.MethodName) > 0 {
		methods = filterMethods(methods, req.MethodName)
	}
	sort.Strings(methods)
	var methodDetails = make([]model.GatewayMethodDetail, 0)
	for _, method := range methods {
		gm, ok := gatewaymethod.GetGatewayMethodCache().GetCache(method)
		if ok {
			var m []byte = nil
			if &gm != nil {
				m, _ = json.Marshal(gm)
			}
			methodDetails = append(methodDetails, model.GatewayMethodDetail{
				MethodMetaInfo: string(m),
				MethodName:     gm.GatewayMethodName,
				SystemName:     appmetadata.GetAppMetadataCache().GetAppName(gm.GatewayMethodName),
			})
		} else {
			methodDetails = append(methodDetails, model.GatewayMethodDetail{
				MethodName: method,
				SystemName: appmetadata.GetAppMetadataCache().GetAppName(method),
			})
		}
	}
	res.Total = int64(len(methods))
	res.MethodDetails = paginate(methodDetails, req.PageIndex, req.PageSize)
	return res, nil
}
func filterMethods(methods []string, filter string) []string {
	var filteredMethods []string
	for _, p := range methods {
		if strings.Contains(p, strings.ToUpper(filter)) {
			filteredMethods = append(filteredMethods, p)
		}
	}
	return filteredMethods
}
func paginate(data []model.GatewayMethodDetail, pageNumber int, pageSize int) []model.GatewayMethodDetail {
	// 获取原始数据长度
	dataLen := len(data)

	// 计算分页开始和结束的索引
	startIndex := (pageNumber - 1) * pageSize
	endIndex := startIndex + pageSize

	// 处理边界情况
	if startIndex > dataLen {
		return []model.GatewayMethodDetail{}
	}
	if endIndex > dataLen {
		endIndex = dataLen
	}
	// 返回分页后的切片
	return data[startIndex:endIndex]
}
