package service

import (
	"worldbar/service/newsDynamicsComment/model"
)

func Create(entity *model.WbNewsDynamicsComment) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.WbNewsDynamicsComment).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.WbNewsDynamicsComment).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.WbNewsDynamicsComment).Delete(query)
}

func Get(query *model.GetQuery) (*model.WbNewsDynamicsComment, error) {
	return new(model.WbNewsDynamicsComment).GetByID(query)
}

func ToggleDisabled(query *model.DisabledQuery) error {
	return new(model.WbNewsDynamicsComment).ToggleDisabled(query)
}
