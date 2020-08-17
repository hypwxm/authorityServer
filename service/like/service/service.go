package service

import (
	"babygrowing/service/like/model"
)

func Create(entity *model.WbLike) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.WbLike).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.WbLike).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.WbLike).Delete(query)
}

func Get(query *model.GetQuery) (*model.WbLike, error) {
	return new(model.WbLike).GetByID(query)
}
