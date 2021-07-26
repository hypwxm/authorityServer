package service

import (
	"log"
	"testing"

	"github.com/hypwxm/authorityServer/service/admin/org/model"
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
