package model

type GatewayConfigReq struct {
	AppId     string `json:"appId,required"`
	PageIndex int64  `json:"pageIndex"`
	PageSize  int64  `json:"pageSize"`
}

type GatewayConfigResponse struct {
	Total  int64            `json:"total"`
	Models []*GatewayConfig `json:"models"`
}
