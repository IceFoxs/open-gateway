package model

type GatewaySystemConfig struct {
	SystemId   string `json:"systemId" column:"system_id"`
	SystemName string `json:"systemName" column:"system_id"`
}

func (c *GatewaySystemConfig) TableName() string {
	return "gateway_system_config"
}
