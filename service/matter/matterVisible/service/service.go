package service

import (
	"worldbar/service/matter/matterVisible/model"
)

func Create(entity []model.WbMatterVisible) error {
	return new(model.WbMatterVisible).Insert(entity)
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.WbMatterVisible).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.WbMatterVisible).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.WbMatterVisible).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.WbMatterVisible).GetByID(query)
}

func UpdateSort(query *model.UpdateSortQuery) error {
	return new(model.WbMatterVisible).UpdateSort(query)
}

func UpdateStatus(query *model.UpdateStatusQuery) error {
	return new(model.WbMatterVisible).UpdateStatus(query)
}
