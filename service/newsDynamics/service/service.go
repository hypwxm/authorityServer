package service

import (
	"worldbar/service/newsDynamics/model"
)

func Create(entity *model.WbNewsDynamics) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.WbNewsDynamics).Update(updateQuery)
}

func List(query *model.Query) ([]*model.WbNewsDynamics, int64, error) {
	return new(model.WbNewsDynamics).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.WbNewsDynamics).Delete(query)
}

func UpdateSort(query *model.UpdateSortQuery) error {
	return new(model.WbNewsDynamics).UpdateSort(query)
}
