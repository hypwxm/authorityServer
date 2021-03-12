package appGorm

import (
	"babygrow/config"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var psql_db *gorm.DB

func Open() *gorm.DB {
	var err error
	if psql_db == nil {

		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second,   // 慢 SQL 阈值
				LogLevel:      logger.Silent, // Log level
				Colorful:      false,         // 禁用彩色打印
			},
		)

		dbAddr := config.Config.Pgsql
		// psql_db, err = sqlx.Connect("postgres", "port=5432 user=postgres password=123456 dbname=brush sslmode=disable")
		psql_db, err = gorm.Open(postgres.Open(dbAddr), &gorm.Config{
			Logger: newLogger,
		})

		// psql_db, err = sql.Open("postgres", "port=5432 user=postgresql password=123456 dbname=brush sslmode=disable")
		if err != nil {

			return nil
		}
		// psql_db.SetMaxIdleConns(200)
		// psql_db.SetMaxOpenConns(50)
		// logger.Logger.Info("pgsql数据库连接成功")
		return psql_db
	}
	return psql_db
}
