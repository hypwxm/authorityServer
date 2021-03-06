package envinit

import (
	"io/ioutil"

	"github.com/hypwxm/authorityServer/DB/pgsql"
	"github.com/hypwxm/authorityServer/config"

	"log"
	"testing"
)

func TestMyBabiesInit(t *testing.T) {
	if config.Env != "development" {
		log.Fatalln("环境错误")
	}
	db := pgsql.Open()

	sql, err := ioutil.ReadFile("sqls/g_member_baby.sql")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatalln(err)
	}
}

func TestauthorityServerInit(t *testing.T) {
	if config.Env != "development" {
		log.Fatalln("环境错误")

	}
	db := pgsql.Open()

	sql, err := ioutil.ReadFile("sqls/g_member_baby_grow.sql")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatalln(err)
	}
}

func TestauthorityServerCommentInit(t *testing.T) {
	if config.Env != "development" {
		log.Fatalln("环境错误")

	}
	db := pgsql.Open()

	sql, err := ioutil.ReadFile("sqls/g_member_baby_grow_comment.sql")
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

func TestGMenuInit(t *testing.T) {
	if config.Env != "development" {
		log.Fatalln("环境错误")
	}
	db := pgsql.Open()

	sql, err := ioutil.ReadFile("sqls/g_menu.sql")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatalln(err)
	}
}

func TestGRoleMenuInit(t *testing.T) {
	if config.Env != "development" {
		log.Fatalln("环境错误")
	}
	db := pgsql.Open()

	sql, err := ioutil.ReadFile("sqls/g_role_menu.sql")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatalln(err)
	}
}

func TestGMemberInit(t *testing.T) {
	if config.Env != "development" {
		log.Fatalln("环境错误")
	}
	db := pgsql.Open()

	sql, err := ioutil.ReadFile("sqls/g_member.sql")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatalln(err)
	}
}

func TestGMemberBabyRelationInit(t *testing.T) {
	if config.Env != "development" {
		log.Fatalln("环境错误")
	}
	db := pgsql.Open()

	sql, err := ioutil.ReadFile("sqls/g_member_baby_relation.sql")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatalln(err)
	}
}

func TestGMemberFamilyInit(t *testing.T) {
	if config.Env != "development" {
		log.Fatalln("环境错误")
	}
	db := pgsql.Open()

	sql, err := ioutil.ReadFile("sqls/g_member_family.sql")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatalln(err)
	}
}

func TestGMemberFamilyMemberInit(t *testing.T) {
	if config.Env != "development" {
		log.Fatalln("环境错误")
	}
	db := pgsql.Open()

	sql, err := ioutil.ReadFile("sqls/g_member_family_member.sql")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatalln(err)
	}
}
