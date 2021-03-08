package service

import (
	"babygrow/event"
	"babygrow/service/member/family/model"
	"context"
)

func Create(ctx context.Context, entity *model.GFamily) (string, error) {
	return entity.Insert(ctx)
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.GFamily).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.GFamily).List(query)
}

func Del(query *model.DeleteQuery, ch chan int) error {
	return new(model.GFamily).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.GFamily).GetByID(query)
}

func ToggleDisabled(query *model.DisabledQuery) error {
	return new(model.GFamily).ToggleDisabled(query)
}

func init() {
	event.Ebus.Subscribe("memberSv:familyDelete", Del)
}
