package dbModel

import (
	"testing"

	"github.com/hypwxm/authorityServer/DB/appGorm"
)

func TestModels(t *testing.T) {
	db := appGorm.Open()
	// db.Migrator().DropTable(&Media{})
	db.AutoMigrate(&GRole{})

}
