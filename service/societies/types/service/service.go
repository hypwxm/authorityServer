package service

import (
	"babygrowing/service/societies/types/model"
)

func Create(entity *model.WbSocietiesType) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.WbSocietiesType).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.WbSocietiesType).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.WbSocietiesType).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.WbSocietiesType).GetByID(query)
}

func ToggleDisabled(query *model.DisabledQuery) error {
	return new(model.WbSocietiesType).ToggleDisabled(query)
}
