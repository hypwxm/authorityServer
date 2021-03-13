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
	db.AutoMigrate(&GMessage{})
}

func TestDB(t *testing.T) {
	db := appGorm.Open()
	db.Create(&GMessage{})
}

func TestFind(t *testing.T) {
	db := appGorm.Open()
	m := new(GMessage)
	db.Find(m)
	t.Fatalf("%+v", m)
}

func TestDel(t *testing.T) {
	db := appGorm.Open()
	m := new(GMessage)
	db.First(m)
	db.Debug().Delete(m)
	t.Fatalf("%+v", m)
}
