package model

import (
	"babygrow/DB/appGorm"
	"babygrow/config"
	"testing"
)

func TestModels(t *testing.T) {
	if config.Env != "development" {
		t.Fatal("环境错误")
	}
	db := appGorm.Open()
	db.Migrator().DropTable(&GDailyComment{})
	db.AutoMigrate(&GDailyComment{})
}
