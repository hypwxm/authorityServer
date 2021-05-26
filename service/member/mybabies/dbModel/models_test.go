package dbModel

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
	db.Migrator().DropTable(&GMyBabies{})
	db.AutoMigrate(&GMyBabies{})

}

func TestModels2(t *testing.T) {
	if config.Env != "development" {
		t.Fatal("环境错误")
	}
	db := appGorm.Open()
	db.Migrator().DropTable(&GMemberBabyRelationApply{})
	db.AutoMigrate(&GMemberBabyRelationApply{})
}
