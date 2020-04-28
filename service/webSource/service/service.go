package service

import (
	"worldbar/service/webSource/model"
)

func Create(entity *model.WbSettingsSource) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.WbSettingsSource).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.WbSettingsSource).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.WbSettingsSource).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.WbSettingsSource).GetByID(query)
}

func ToggleDisabled(query *model.DisabledQuery) error {
	return new(model.WbSettingsSource).ToggleDisabled(query)
}
