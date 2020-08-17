package houseModel

import (
	"io/ioutil"
	"testing"
	"babygrowing/DB/pgsql"
	"babygrowing/config"
)

func TestModels(t *testing.T) {
	if config.Env != "development" {
		t.Fatal("环境错误")
	}
	db := pgsql.Open()
	b, err := ioutil.ReadFile("scheme.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(b))
	if err != nil {
		t.Fatal(err)
	}
}