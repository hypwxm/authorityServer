package service

import (
	"babygrow/service/admin/rolePermission/menu/model"
)

func Create(query *model.SaveQuery) (string, error) {
	return new(model.GRoleMenu).Save(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.GRoleMenu).Delete(query)
}

func List(query *model.Query) ([]*model.ListModel, error) {
	return new(model.GRoleMenu).List(query)
}
