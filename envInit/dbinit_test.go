package envInit

import (
	"babygrowing/DB/pgsql"
	"babygrowing/config"
	"babygrowing/service/daily/model"
	media "babygrowing/service/media/model"
	mybabies "babygrowing/service/mybabies/model"

	"log"
	"testing"
)

func TestMyBabiesInit(t *testing.T) {
	if config.Env != "development" {
		log.Fatalln("环境错误")
	}
	db := pgsql.Open()

	sql, err := mybabies.GetSqlFile()
	if err != nil {
		log.Fatalln("err")
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatalln("环境错误")
	}
}

func TestBabyGrowningInit(t *testing.T) {
	if config.Env != "development" {
		log.Fatalln("环境错误")
	}
	db := pgsql.Open()

	sql, err := model.GetSqlFile()
	if err != nil {
		log.Fatalln("err")
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatalln("环境错误")
	}
}

func TestMediaInit(t *testing.T) {
	if config.Env != "development" {
		log.Fatalln("环境错误")
	}
	db := pgsql.Open()

	sql, err := media.GetSqlFile()
	if err != nil {
		log.Fatalln("err")
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatalln("环境错误")
	}
}
