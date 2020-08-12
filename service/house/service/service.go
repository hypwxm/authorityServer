package service

import (
	"worldbar/service/house/model"
)

func Create(entity *model.WbHouse) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.WbHouse).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, error) {
	return new(model.WbHouse).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.WbHouse).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.WbHouse).GetByID(query)
}

func ToggleDisabled(query *model.DisabledQuery) error {
	return new(model.WbHouse).ToggleDisabled(query)
}
