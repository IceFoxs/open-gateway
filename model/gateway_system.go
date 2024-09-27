package model

type GatewaySystemConfig struct {
	SystemId   string `json:"systemId,required" column:"system_id" gorm:"primarykey" vd:"@:len($)>0; msg:'systemId不能为空'"`
	SystemName string `json:"systemName" column:"system_name"`
}

func (c *GatewaySystemConfig) TableName() string {
	return "gateway_system_config"
}
