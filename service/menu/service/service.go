package service

import (
	"babygrowing/service/menu/model"
)

func Create(entity *model.WbSettingsMenu) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.WbSettingsMenu).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, error) {
	return new(model.WbSettingsMenu).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.WbSettingsMenu).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.WbSettingsMenu).GetByID(query)
}

func ToggleDisabled(query *model.DisabledQuery) error {
	return new(model.WbSettingsMenu).ToggleDisabled(query)
}
