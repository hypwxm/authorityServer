package service

import (
	"babygrowing/service/admin/menu/model"
)

func Create(entity *model.GMenu) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.GMenu).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.GMenu).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.GMenu).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.GMenu).GetByID(query)
}
