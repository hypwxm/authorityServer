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
		log.Fatalln(err)
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
		log.Fatalln(err)
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
		log.Fatalln(err)
	}
}

func TestOrgInit(t *testing.T) {
	if config.Env != "development" {
		log.Fatalln("环境错误")
	}
	db := pgsql.Open()

	sql, err := ioutil.ReadFile("sqls/g_org.sql")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatalln(err)
	}
}

func TestAdminUserInit(t *testing.T) {
	if config.Env != "development" {
		log.Fatalln("环境错误")
	}
	db := pgsql.Open()

	sql, err := ioutil.ReadFile("sqls/g_admin_user.sql")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatalln(err)
	}
}

func TestAdminRoleInit(t *testing.T) {
	if config.Env != "development" {
		log.Fatalln("环境错误")
	}
	db := pgsql.Open()

	sql, err := ioutil.ReadFile("sqls/g_admin_role.sql")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatalln(err)
	}
}


func TestAdminUserRoleInit(t *testing.T) {
	if config.Env != "development" {
		log.Fatalln("环境错误")
	}
	db := pgsql.Open()

	sql, err := ioutil.ReadFile("sqls/g_admin_user_role.sql")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatalln(err)
	}
}
