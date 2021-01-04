package service

import (
	"babygrowing/service/admin/role/model"
)

func Create(entity *model.GAdminRole) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.GAdminRole).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.GAdminRole).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.GAdminRole).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.GAdminRole).GetByID(query)
}
