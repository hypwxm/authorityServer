package appGorm

import (
	"log"
	"os"
	"time"

	"github.com/hypwxm/authorityServer/config"

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
				SlowThreshold: time.Second, // 慢 SQL 阈值
				LogLevel:      logger.Info, // Log level
				Colorful:      true,        // 禁用彩色打印
			},
		)

		dbAddr := config.Config.Pgsql
		psql_db, err = gorm.Open(postgres.Open(dbAddr), &gorm.Config{
			Logger: newLogger,
		})

		if err != nil {
			psql_db = nil
			return nil
		}
		db, err := psql_db.DB()
		if err != nil {
			psql_db = nil
			return nil
		}
		// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
		db.SetMaxIdleConns(10)
		// SetMaxOpenConns 设置打开数据库连接的最大数量。
		db.SetMaxOpenConns(100)
		return psql_db
	}
	return psql_db
}
