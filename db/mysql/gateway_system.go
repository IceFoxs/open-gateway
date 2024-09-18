package mysql

import "github.com/IceFoxs/open-gateway/model"

func GetGatewaySystemConfig(keyword string) ([]*model.GatewaySystemConfig, error) {
	db := DB.Model(model.GatewaySystemConfig{})
	if len(keyword) != 0 {
		db = db.Where(DB.Or("system_name like ?", "%"+keyword+"%"))
	}
	var res []*model.GatewaySystemConfig
	if err := db.Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
