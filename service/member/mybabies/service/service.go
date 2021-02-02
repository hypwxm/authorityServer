package service

import (
	"babygrowing/service/member/mybabies/model"
)

func Create(entity *model.GMyBabies) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.GMyBabies).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.GMyBabies).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.GMyBabies).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.GMyBabies).GetByID(query)
}

func ToggleDisabled(query *model.DisabledQuery) error {
	return new(model.GMyBabies).ToggleDisabled(query)
}
