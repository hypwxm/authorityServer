package service

import (
	"babygrowing/service/member/daily/model"
)

func Create(entity *model.GDaily) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.GDaily).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.GDaily).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.GDaily).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.GDaily).GetByID(query)
}
