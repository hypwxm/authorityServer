package service

import (
	"babygrow/service/admin/org/model"
	"log"
	"testing"
)

func TestModels(t *testing.T) {
	org := &model.GOrg{
		Name: "顶级组织",
	}
	_, err := org.Insert()
	if err != nil {
		log.Fatal(err)
	}
}
