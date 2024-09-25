package mysql

import (
	"github.com/IceFoxs/open-gateway/model"
)

func QueryGatewayChannelConfig(keyword string, page, pageSize int64) ([]*model.GatewayChannelConfig, int64, error) {
	db := DB.Model(model.GatewayChannelConfig{})
	if len(keyword) != 0 {
		db = db.Where(DB.Or("app_name like ?", "%"+keyword+"%"))
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var res []*model.GatewayChannelConfig
	if err := db.Limit(int(pageSize)).Offset(int(pageSize * (page - 1))).Find(&res).Error; err != nil {
		return nil, 0, err
	}
	return res, total, nil
}

func CreateGatewayChannelConfig(gscs []*model.GatewayChannelConfig) error {
	return DB.Create(gscs).Error
}

func UpdateGatewayChannelConfig(gsc *model.GatewayChannelConfig) error {
	return DB.Updates(gsc).Error
}
func GetGatewayChannelConfig(keyword string) ([]*model.GatewayChannelConfig, error) {
	db := DB.Model(model.GatewayChannelConfig{})
	if len(keyword) != 0 {
		db = db.Where(DB.Or("app_name like ?", "%"+keyword+"%"))
	}
	var res []*model.GatewayChannelConfig
	if err := db.Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
