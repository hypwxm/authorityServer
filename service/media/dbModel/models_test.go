package dbModel

import (
	"authorityServer/DB/appGorm"
	"authorityServer/config"
	"testing"
)

func TestModels(t *testing.T) {
	if config.Env != "development" {
		t.Fatal("环境错误")
	}
	db := appGorm.Open()
	db.Migrator().DropTable(&Media{})
	db.AutoMigrate(&Media{})

}
