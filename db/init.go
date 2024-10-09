package db

import (
	"github.com/IceFoxs/open-gateway/db/mysql"
	"go.uber.org/zap"
)

func Init(dsn string, logger *zap.Logger) {
	mysql.Init(dsn, logger)
}
