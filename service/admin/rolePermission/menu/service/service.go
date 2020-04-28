package service

import (
	"worldbar/service/admin/rolePermission/menu/model"
)

func Create(query *model.SaveQuery) (string, error) {
	return new(model.WbAdminRoleMenuPermission).Save(query)
}

func List(query *model.Query) ([]*model.ListModel, error) {
	return new(model.WbAdminRoleMenuPermission).List(query)
}
