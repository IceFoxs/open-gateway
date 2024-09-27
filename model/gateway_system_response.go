package model

type GatewaySystemReq struct {
	SysId     string `json:"sysId,required"`
	PageIndex int64  `json:"pageIndex"`
	PageSize  int64  `json:"pageSize"`
}

type GatewaySystemResponse struct {
	Total  int64                  `json:"total"`
	Models []*GatewaySystemConfig `json:"models"`
}
