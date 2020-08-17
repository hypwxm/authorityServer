package pgsql

import (
	"babygrowing/config"
	"babygrowing/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var psql_db *sqlx.DB

//port是数据库的端口号，默认是5432，如果改了，这里一定要自定义；
//user就是你数据库的登录帐号;
//dbname就是你在数据库里面建立的数据库的名字;
//sslmode就是安全验证模式;

//还可以是这种方式打开
//db, err := sql.Open("postgres", "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full")

func Open() *sqlx.DB {
	var err error
	if psql_db == nil {
		dbAddr := config.Config.Pgsql
		// psql_db, err = sqlx.Connect("postgres", "port=5432 user=postgres password=123456 dbname=brush sslmode=disable")
		psql_db, err = sqlx.Connect("postgres", dbAddr)
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
