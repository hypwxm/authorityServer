package service

import (
	"babygrowing/service/societies/base/model"
)

func Create(entity *model.WbSocieties) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.WbSocieties).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.WbSocieties).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.WbSocieties).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.WbSocieties).GetByID(query)
}

func ToggleDisabled(query *model.DisabledQuery) error {
	return new(model.WbSocieties).ToggleDisabled(query)
}

func UpdateSort(query *model.UpdateSortQuery) error {
	return new(model.WbSocieties).UpdateSort(query)
}

func UpdateStatus(query *model.UpdateStatusQuery) error {
	return new(model.WbSocieties).UpdateStatus(query)
}
