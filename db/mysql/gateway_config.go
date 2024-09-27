package mysql

import (
	"github.com/IceFoxs/open-gateway/model"
)

func QueryGatewayConfigByPage(keyword string, page, pageSize int64) ([]*model.GatewayConfig, int64, error) {
	db := DB.Model(model.GatewayConfig{})
	if len(keyword) != 0 {
		db = db.Where(DB.Or("app_name like ?", "%"+keyword+"%"))
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var res []*model.GatewayConfig
	if err := db.Limit(int(pageSize)).Offset(int(pageSize * (page - 1))).Find(&res).Error; err != nil {
		return nil, 0, err
	}
	return res, total, nil
}

func CreateGatewayConfig(gscs []*model.GatewayConfig) error {
	return DB.Create(gscs).Error
}

func UpdateGatewayConfig(gsc *model.GatewayConfig) error {
	return DB.Updates(gsc).Error
}
func GetGatewayConfig(keyword string) ([]*model.GatewayConfig, error) {
	db := DB.Model(model.GatewayConfig{})
	if len(keyword) != 0 {
		db = db.Where(DB.Or("app_name like ?", "%"+keyword+"%"))
	}
	var res []*model.GatewayConfig
	if err := db.Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
