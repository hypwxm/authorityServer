package service

import (
	"authorityServer/DB/appGorm"
	"authorityServer/service/admin/menu/dao"
	"authorityServer/service/admin/menu/dbModel"
	"authorityServer/util/interfaces"
	"errors"
	"strings"
)

type CreateModel struct {
	dbModel.GMenu
}

func Create(entity *CreateModel) (string, error) {
	db := appGorm.Open()
	return dao.Insert(db, &entity.GMenu)
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
