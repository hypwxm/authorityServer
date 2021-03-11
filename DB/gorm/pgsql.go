package gorm

import (
	"babygrow/config"
	"babygrow/logger"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var psql_db *gorm.DB

func Open() *gorm.DB {
	var err error
	if psql_db == nil {
		dbAddr := config.Config.Pgsql
		// psql_db, err = sqlx.Connect("postgres", "port=5432 user=postgres password=123456 dbname=brush sslmode=disable")
		psql_db, err = gorm.Open(postgres.Open(dbAddr), &gorm.Config{})

		// psql_db, err = sql.Open("postgres", "port=5432 user=postgresql password=123456 dbname=brush sslmode=disable")
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"event": "数据库链接错误",
			}).Panic(err)
			return nil
		}
		psql_db.SetMaxIdleConns(200)
		psql_db.SetMaxOpenConns(50)
		logger.Logger.Info("pgsql数据库连接成功")
		return psql_db
	}
	return psql_db
}
