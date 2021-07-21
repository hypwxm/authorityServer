package dbModel

import (
	"authorityServer/DB/appGorm"
	"testing"
)

func TestModels(t *testing.T) {
	db := appGorm.Open()
	// db.Migrator().DropTable(&Media{})
	db.AutoMigrate(&GMenu{})

}
