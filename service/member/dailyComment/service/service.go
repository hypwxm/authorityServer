package service

import (
	"babygrow/service/member/dailyComment/model"
)

func Create(entity *model.GDailyComment) (string, error) {
	return entity.Insert()
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.GDailyComment).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.GDailyComment).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.GDailyComment).GetByID(query)
}
