package db

import "github.com/IceFoxs/open-gateway/db/mysql"

func Init(dsn string) {
	mysql.Init(dsn)
}
