package dbModel

import (
	"testing"

	"github.com/hypwxm/authorityServer/DB/appGorm"
	"github.com/hypwxm/authorityServer/config"
)

func TestModels(t *testing.T) {
	if config.Env != "development" {
		t.Fatal("环境错误")
	}
	db := appGorm.Open()
	db.Migrator().DropTable(&Media{})
	db.AutoMigrate(&Media{})

}
