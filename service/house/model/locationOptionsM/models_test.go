package locationOptionsM

import (
	"io/ioutil"
	"testing"
	"worldbar/DB/pgsql"
	"worldbar/config"
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