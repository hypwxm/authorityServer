package service

import (
	"babygrow/service/member/family/model"
)

func Create(entity *model.GFamily) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.GFamily).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.GFamily).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.GFamily).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.GFamily).GetByID(query)
}

func ToggleDisabled(query *model.DisabledQuery) error {
	return new(model.GFamily).ToggleDisabled(query)
}
