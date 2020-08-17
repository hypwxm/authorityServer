package service

import (
	"babygrowing/service/matter/matter/model"
)

func Create(entity *model.WbMatter) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.WbMatter).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.WbMatter).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.WbMatter).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.WbMatter).GetByID(query)
}

func UpdateSort(query *model.UpdateSortQuery) error {
	return new(model.WbMatter).UpdateSort(query)
}

func UpdateStatus(query *model.UpdateStatusQuery) error {
	return new(model.WbMatter).UpdateStatus(query)
}
