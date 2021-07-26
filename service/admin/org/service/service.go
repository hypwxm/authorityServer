package service

import (
	"errors"
	"strings"

	"github.com/hypwxm/authorityServer/DB/appGorm"
	"github.com/hypwxm/authorityServer/service/admin/org/dao"
	"github.com/hypwxm/authorityServer/service/admin/org/dbModel"
	"github.com/hypwxm/authorityServer/util/interfaces"
)

type CreateModel struct {
	dbModel.GOrg
	UserId string `json:"userId"`
}

func Create(entity *CreateModel) (string, error) {
	db := appGorm.Open()
	return dao.Insert(db, &entity.GOrg)
}

func Modify(query interfaces.QueryInterface) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.GetID()) == "" {
		return errors.New("更新条件错误")
	}
	db := appGorm.Open()
	return dao.Update(db, query)
}

func List(query interfaces.QueryInterface) (interfaces.ModelMapSlice, int64, error) {
	db := appGorm.Open()
	list, total, err := dao.List(db, query)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func Del(query interfaces.QueryInterface) error {
	db := appGorm.Open()
	return dao.Delete(db, query)
}

func Get(query interfaces.QueryInterface) (interfaces.ModelInterface, error) {
	db := appGorm.Open()
	return dao.Get(db, query)
}
