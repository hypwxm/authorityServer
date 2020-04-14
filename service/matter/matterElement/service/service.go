package service

import (
	"worldbar/service/matter/matterElement/model"
)

func Create(entity *model.WbMatterElement) (string, error) {
	return entity.Insert()
}

func Modify(updateQuery *model.UpdateByIDQuery) error {
	return new(model.WbMatterElement).Update(updateQuery)
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return new(model.WbMatterElement).List(query)
}

func Del(query *model.DeleteQuery) error {
	return new(model.WbMatterElement).Delete(query)
}

func Get(query *model.GetQuery) (*model.GetModel, error) {
	return new(model.WbMatterElement).GetByID(query)
}
