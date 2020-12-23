package envinit

import (
	"babygrowing/DB/pgsql"
	"babygrowing/config"
	"io/ioutil"

	"log"
	"testing"
)

func TestMyBabiesInit(t *testing.T) {
	if config.Env != "development" {
		log.Fatalln("环境错误")
	}
	db := pgsql.Open()

	sql, err := ioutil.ReadFile("sqls/g_my_babies.sql")
	if err != nil {
		log.Fatalln(err)
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

	sql, err := ioutil.ReadFile("sqls/g_daily.sql")
	if err != nil {
		log.Fatalln(err)
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

	sql, err := ioutil.ReadFile("sqls/g_media.sql")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatalln("环境错误")
	}
}
