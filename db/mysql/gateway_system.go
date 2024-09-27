package mysql

import "github.com/IceFoxs/open-gateway/model"

func CreateGatewaySystem(gscs []*model.GatewaySystemConfig) error {
	return DB.Create(gscs).Error
}
func GetGatewaySystemConfigByPage(keyword string, page, pageSize int64) ([]*model.GatewaySystemConfig, int64, error) {
	db := DB.Model(model.GatewaySystemConfig{})
	if len(keyword) != 0 {
		db = db.Where(DB.Or("system_name like ?", "%"+keyword+"%")).Or("system_id like ?", "%"+keyword+"%")
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var res []*model.GatewaySystemConfig
	if err := db.Limit(int(pageSize)).Offset(int(pageSize * (page - 1))).Find(&res).Error; err != nil {
		return nil, 0, err
	}
	return res, total, nil
}

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
