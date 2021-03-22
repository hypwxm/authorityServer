package model

import (
	"babygrow/DB/appGorm"
	"babygrow/DB/pgsql"
	"babygrow/config"
	"context"
	"log"
	"testing"
)

func TestModels(t *testing.T) {
	if config.Env != "development" {
		t.Fatal("环境错误")
	}
	db := appGorm.Open()
	db.Migrator().DropTable(&GFamily{})
	db.AutoMigrate(&GFamily{})
}

func TestContext(t *testing.T) {
	ctx := context.Background()
	db := pgsql.Open()
	db.Close()
	cw := context.WithValue(ctx, "a", pgsql.Open())
	cy := context.WithValue(cw, "a", pgsql.Open())

	log.Fatal(cw.Value("a"), cy.Value("a"))
}
