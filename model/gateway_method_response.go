package model

type GatewayMethodResponse struct {
	Total         int64                 `json:"total"`
	MethodDetails []GatewayMethodDetail `json:"methodDetails"`
}

type GatewayMethodRequest struct {
	PageIndex  int    `json:"pageIndex"`
	PageSize   int    `json:"pageSize"`
	MethodName string `json:"methodName"`
}

type GatewayMethodDetail struct {
	SystemName     string `json:"systemName"`
	MethodName     string `json:"methodName"`
	MethodMetaInfo string `json:"methodMetaInfo"`
}
