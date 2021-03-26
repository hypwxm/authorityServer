package model

import (
	"babygrow/DB/appGorm"
	"babygrow/config"
	"context"
	"testing"
)

func TestModels(t *testing.T) {
	if config.Env != "development" {
		t.Fatal("环境错误")
	}
	db := appGorm.Open()
	err := db.AutoMigrate(&GFamilyMembers{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdate(t *testing.T) {
	e := new(GFamilyMembers)
	e.Update(&UpdateByIDQuery{
		ID:       "1",
		Nickname: "Asdasd",
	})
	t.Fatal("sa")
}

func TestDel(t *testing.T) {
	e := new(GFamilyMembers)
	e.Delete(context.Background(), &DeleteQuery{
		IDs: []string{"1"},
	})
	t.Fatal("sa")
}
