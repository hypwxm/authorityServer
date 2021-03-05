package model

import (
	"babygrow/DB/pgsql"
	"babygrow/config"
	"context"
	"io/ioutil"
	"log"
	"testing"
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

func TestContext(t *testing.T) {
	ctx := context.Background()
	db := pgsql.Open()
	db.Close()
	cw := context.WithValue(ctx, "a", pgsql.Open())
	cy := context.WithValue(cw, "a", pgsql.Open())

	log.Fatal(cw.Value("a"), cy.Value("a"))
}
